// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"encoding/gob"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/lrstanley/pt"
	"github.com/timshannon/bolthold"
)

var (
	tmpl *pt.Loader
)

func httpServer(ctx context.Context, closer chan struct{}) {
	tmpl = pt.New("", pt.Config{
		CacheParsed:     !conf.Debug,
		Loader:          rice.MustFindBox("static").Bytes,
		ErrorLogger:     os.Stderr,
		DefaultCtx:      tmplDefaultCtx,
		NotFoundHandler: http.NotFound,
	})

	gob.Register(flashMessage{})
	updateGlobalStats(nil)

	r := chi.NewRouter()

	if conf.Proxy {
		r.Use(middleware.RealIP)
	}
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.DefaultLogger)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.Recoverer)
	r.Use(middleware.GetHead)

	// Mount the static directory (in-memory and disk) to the /static route.
	pt.FileServer(r, "/static", rice.MustFindBox("static").HTTPBox())

	if conf.Debug {
		r.Mount("/debug", middleware.Profiler())
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Render(w, r, "tmpl/index.html", nil)
	})
	r.Get("/-/abuse", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Render(w, r, "tmpl/abuse.html", pt.M{"safebrowsing": conf.SafeBrowsing.APIKey != ""})
	})

	r.Get("/{uid}", expand)
	r.Post("/{uid}", expand)
	r.Post("/", addForm)
	r.Post("/add", addAPI)

	srv := &http.Server{
		Addr:         conf.HTTP,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		debug.Printf("initializing http server on %s", conf.HTTP)

		var err error
		if conf.TLS.Enable {
			err = srv.ListenAndServeTLS(conf.TLS.Cert, conf.TLS.Key)
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			debug.Printf("http error: %v", err)
		}

		close(closer)
	}()

	// Wait for parent to ask to shutdown.
	<-ctx.Done()

	debug.Printf("requesting http server to shutdown")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil && err != http.ErrServerClosed {
		debug.Fatalf("unable to shutdown http server: %v", err)
	}
}

func tmplDefaultCtx(w http.ResponseWriter, r *http.Request) (ctx map[string]interface{}) {
	if ctx == nil {
		ctx = make(map[string]interface{})
	}

	cachedGlobalStats.mu.RLock()
	// Note that this copies a mutex, but it should never be re-locked, as
	// it's only being used in a template.
	stats := cachedGlobalStats
	cachedGlobalStats.mu.RUnlock()

	ctx = pt.M{
		"full_url":          r.URL.String(),
		"url":               r.URL,
		"commit":            commit,
		"version":           version,
		"stats":             &stats,
		"http_pre_include":  conf.HTTPPreInclude,
		"http_post_include": conf.HTTPPostInclude,
	}

	return ctx
}

type flashMessage struct {
	Type string
	Text string
}

type httpResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	URL     string `json:"url,omitempty"`
}

func expand(w http.ResponseWriter, r *http.Request) {
	db := newDB(true)

	var link Link
	err := db.Get(chi.URLParam(r, "uid"), &link)
	db.Close()

	if err != nil {
		if err == bolthold.ErrNotFound {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		panic(err)
	}

	if link.EncryptionHash != "" {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusForbidden)
			tmpl.Render(w, r, "tmpl/auth.html", nil)
			return
		}

		decrypt := r.PostFormValue("decrypt")

		if hash(decrypt) != link.EncryptionHash {
			w.WriteHeader(http.StatusForbidden)
			tmpl.Render(w, r, "tmpl/auth.html", pt.M{"message": flashMessage{"danger", "invalid decryption string provided"}})
			return
		}

		// Assume at this point they have provided authentication for the
		// link, and we can redirect them.
	}

	link.AddHit()

	if conf.SafeBrowsing.APIKey != "" {
		threats, err := safeBrowser.LookupURLs([]string{link.URL})
		if err != nil {
			debug.Printf("safebrowsing error: %v", err)

			if conf.SafeBrowsing.RedirectFallback {
				http.Redirect(w, r, link.URL, http.StatusFound)
				return
			}

			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}

		if len(threats) > 0 {
			tmpl.Render(w, r, "tmpl/safebrowsing.html", pt.M{"threats": threats[0], "link": link.URL})
			return
		}
	}

	http.Redirect(w, r, link.URL, http.StatusFound)
}

func addForm(w http.ResponseWriter, r *http.Request) {
	link := &Link{
		URL:            r.PostFormValue("url"),
		EncryptionHash: hash(r.PostFormValue("encrypt")),
	}

	link.Author, _, _ = net.SplitHostPort(r.RemoteAddr)
	if link.Author == "" {
		link.Author = r.RemoteAddr
	}

	if err := link.Create(nil); err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		tmpl.Render(w, r, "tmpl/index.html", pt.M{"message": flashMessage{"danger", err.Error()}})
		return
	}

	tmpl.Render(w, r, "tmpl/index.html", pt.M{"link": link})
}

func addAPI(w http.ResponseWriter, r *http.Request) {
	link := Link{
		URL:            r.PostFormValue("url"),
		EncryptionHash: hash(r.PostFormValue("encrypt")),
	}

	// Check for old password supplying method.
	if link.EncryptionHash == "" {
		link.EncryptionHash = hash(r.PostFormValue("password"))
	}

	link.Author, _, _ = net.SplitHostPort(r.RemoteAddr)
	if link.Author == "" {
		link.Author = r.RemoteAddr
	}

	if err := link.Create(nil); err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		pt.JSON(w, r, httpResp{Success: false, Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	pt.JSON(w, r, httpResp{Success: true, URL: link.Short()})
}

var validSchemes = []string{
	"http",
	"https",
	"ftp",
}

func isValidScheme(scheme string) bool {
	scheme = strings.ToLower(scheme)

	for i := 0; i < len(validSchemes); i++ {
		if validSchemes[i] == scheme {
			return true
		}
	}

	return false
}
