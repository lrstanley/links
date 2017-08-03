// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"net/http"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/flosch/pongo2"
	"github.com/pressly/chi"
	"github.com/sharpner/pobin"
)

type FlashMessage struct {
	Type string
	Text string
}

const defaultSessionID = "sessionid"

var assetfs *pongo2.TemplateSet

func setupTmpl() {
	assetfs = pongo2.NewSet("", pobin.NewMemoryTemplateLoader(rice.MustFindBox("static").Bytes))

}

func tmpl(w http.ResponseWriter, r *http.Request, path string, ctx map[string]interface{}) {
	session, _ := sess.Get(r, defaultSessionID)
	messages := session.Flashes("messages")

	// We have to save the session, otherwise the flashes aren't properly
	// cleared either.
	err := session.Save(r, w)
	if err != nil {
		panic(err)
	}

	tpl := pongo2.Must(assetfs.FromFile(path))

	if ctx == nil {
		ctx = make(map[string]interface{})
	}

	// cachedGlobalStats.mu.RLock()
	// // Note that this copies a mutex, but it should never be re-locked, as
	// // it's only being used in a template.
	// stats := cachedGlobalStats
	// cachedGlobalStats.mu.RUnlock()

	ctx["full_url"] = r.URL.String()
	ctx["url"] = r.URL
	ctx["sess"] = session.Values
	ctx["messages"] = messages
	ctx["commit"] = commit
	ctx["version"] = version
	// ctx["stats"] = &stats

	out, err := tpl.ExecuteBytes(ctx)
	if err != nil {
		panic(err)
	}

	w.Write(out)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("url params not allowed in file server")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func flashMessage(w http.ResponseWriter, r *http.Request, status, message string) {
	session, _ := sess.Get(r, defaultSessionID)
	session.AddFlash(FlashMessage{status, message}, "messages")
	err := session.Save(r, w)
	if err != nil {
		panic(err)
	}
}
