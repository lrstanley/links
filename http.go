// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	rice "github.com/GeertJohan/go.rice"
	gctx "github.com/gorilla/context"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/timshannon/bolthold"
)

var tmpl *Loader
var sess *sessions.CookieStore

func httpServer() {
	tmpl = NewLoader("partials/*", defaultCtx)
	gob.Register(FlashMessage{})
	updateGlobalStats(nil)

	sess = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))
	sess.MaxAge(86400)

	r := chi.NewRouter()

	if conf.Proxy {
		r.Use(middleware.RealIP)
	}
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.DefaultLogger)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.Recoverer)
	r.FileServer("/static", rice.MustFindBox("static").HTTPBox())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Load(w, r, "tmpl/index.html", nil)
	})

	r.Get("/:uid", expand)
	r.Post("/:uid", expand)
	r.Post("/", addForm)
	r.Post("/add", addAPI)

	if conf.TLS.Enable {
		debug.Printf("initializing https server on %s", conf.HTTP)
		debug.Fatal(http.ListenAndServeTLS(conf.HTTP, conf.TLS.Cert, conf.TLS.Key, gctx.ClearHandler(r)))
	}

	debug.Printf("initializing http server on %s", conf.HTTP)
	debug.Fatal(http.ListenAndServe(conf.HTTP, gctx.ClearHandler(r)))
}

func mustJSON(input interface{}) []byte {
	out, err := json.Marshal(input)
	if err != nil {
		panic(fmt.Sprintf("unable to marshal json: %s", input))
	}

	return out
}

type HTTPResp struct {
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
			tmpl.Load(w, r, "tmpl/auth.html", nil)
			return
		}

		decrypt := r.PostFormValue("decrypt")

		if hash(decrypt) != link.EncryptionHash {
			tmpl.Flash(w, r, "danger", "invalid decryption string provided")
			w.WriteHeader(http.StatusForbidden)
			tmpl.Load(w, r, "tmpl/auth.html", nil)
			return
		}

		// Assume at this point they have provided authentication for the
		// link, and we can redirect them.
	}

	link.AddHit()

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

	if err := link.Create(); err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		tmpl.Flash(w, r, "danger", err.Error())
		tmpl.Load(w, r, "tmpl/index.html", nil)
		return
	}

	tmpl.Load(w, r, "tmpl/index.html", map[string]interface{}{"link": link})
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

	if err := link.Create(); err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write(mustJSON(HTTPResp{Success: false, Message: err.Error()}))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(mustJSON(HTTPResp{Success: true, URL: link.Short()}))
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
