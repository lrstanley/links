package main

import (
	"net/http"
	"time"

	rice "github.com/GeertJohan/go.rice"
	gctx "github.com/gorilla/context"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
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
		w.Write([]byte("ohai"))
	})

	if conf.TLS.Enable {
		debug.Printf("initializing https server on %s", conf.HTTP)
		debug.Fatal(http.ListenAndServeTLS(conf.HTTP, conf.TLS.Cert, conf.TLS.Key, gctx.ClearHandler(r)))
	}

	debug.Printf("initializing http server on %s", conf.HTTP)
	debug.Fatal(http.ListenAndServe(conf.HTTP, gctx.ClearHandler(r)))
}
