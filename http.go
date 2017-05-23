package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
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

	sess = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))
	sess.MaxAge(86400 * 2)

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
	r.Post("/add", addLink)

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

func render(w http.ResponseWriter, log string, resp *HTTPResp) {
	if log == "" {
		w.Write(mustJSON(resp))
		return
	}

	// TODO: render html, add messages as flashes and links as ctx.
}

type HTTPResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,ommitempty"`
	URL     string `json:"url,omitempty"`
}

func addLink(w http.ResponseWriter, r *http.Request) {
	raw := r.PostFormValue("url")
	passwd := r.PostFormValue("encrypt")
	// r.PostFormValue("location")

	// Check for old password supplying method.
	if passwd == "" {
		passwd = r.PostFormValue("password")
	}

	if len(raw) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(mustJSON(HTTPResp{Success: false, Message: "please supply a url to shorten"}))
		return
	}

	uri, err := url.Parse(raw)
	if err != nil || uri.Host == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write(mustJSON(HTTPResp{Success: false, Message: "unable to parse url: " + raw}))
		return
	}

	if !isValidScheme(uri.Scheme) {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write(mustJSON(HTTPResp{Success: false, Message: "invalid url scheme. allowed schemes: " + strings.Join(validSchemes, ", ")}))
		return
	}

	link := Link{URL: uri.String(), Created: time.Now(), Hits: 0}

	link.Author, _, _ = net.SplitHostPort(r.RemoteAddr)
	if link.Author == "" {
		link.Author = r.RemoteAddr
	}

	if passwd != "" {
		link.EncryptionHash = hash(passwd)
	}

	db := newDB(false)
	defer db.Close()

	// Check for dups.
	var result []Link
	err = db.Find(&result, bolthold.Where("URL").Eq(link.URL).And("EncryptionHash").Eq(link.EncryptionHash).Limit(1))
	if err != nil {
		panic(err)
	}

	// Assume there is a dup, just return it to the user.
	if len(result) > 0 {
		w.WriteHeader(http.StatusOK)
		w.Write(mustJSON(HTTPResp{Success: true, URL: result[0].Short()}))
		return
	}

	// Store it.
	for {
		link.UID = uuid(4)
		err = db.Insert(link.UID, link)
		if err != nil {
			if err == bolthold.ErrKeyExists {
				// Keep looping through until we're able to store one which
				// doesn't collide with a pre-existing key.
				continue
			}

			panic(err)
		}

		break
	}

	w.WriteHeader(http.StatusOK)
	w.Write(mustJSON(HTTPResp{Success: true, URL: link.Short()}))

	debug.Printf("%#v", link)
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
