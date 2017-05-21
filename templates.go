package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/aymerick/raymond"
	humanize "github.com/dustin/go-humanize"
	"github.com/gorilla/sessions"
)

func init() {
	raymond.RegisterHelper("hrt", func(t time.Time) string {
		if t.IsZero() {
			return ""
		}

		return humanize.Time(t)
	})

	raymond.RegisterHelper("ellipsis", func(max int, text string) string {
		if len(text) > max {
			return text[0:max] + "..."
		}

		return text
	})

	raymond.RegisterHelper("ne", func(a interface{}, b interface{}, options *raymond.Options) interface{} {
		if raymond.Str(a) != raymond.Str(b) {
			return options.Fn()
		}

		return ""
	})
}

func ListPartials(path string) map[string]string {
	globs := []string{}

	rice.MustFindBox("static").Walk("", func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(path, "partials/") {
			globs = append(globs, path)
		}

		return nil
	})

	var j int
	var name string

	out := make(map[string]string)
	for i := 0; i < len(globs); i++ {
		name = filepath.Base(globs[i])
		j = strings.Index(name, ".")
		if j > -1 {
			name = name[0:j]
		}

		out[name] = rice.MustFindBox("static").MustString(globs[i])
	}

	return out
}

type Loader struct {
	Partials string

	ctx func(w http.ResponseWriter, r *http.Request) map[string]interface{}
}

func NewLoader(partials string, defaultCtx func(w http.ResponseWriter, r *http.Request) map[string]interface{}) *Loader {
	return &Loader{Partials: partials, ctx: defaultCtx}
}

func (l *Loader) Get(path string) *raymond.Template {
	tmpl := raymond.MustParse(rice.MustFindBox("static").MustString(path))
	tmpl.RegisterPartials(ListPartials(l.Partials))

	return tmpl
}

func (l *Loader) Load(w http.ResponseWriter, r *http.Request, path string, ctx interface{}) {
	if l.ctx == nil {
		w.Write([]byte(l.Get(path).MustExec(ctx)))
		return
	}

	out := l.ctx(w, r)
	if out == nil {
		panic("default context returned nil")
	}

	out["ctx"] = ctx

	w.Write([]byte(l.Get(path).MustExec(out)))
}

func defaultCtx(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	session := getSession(r)
	// messages := session.Flashes("messages")

	// We have to save the session, otherwise the flashes aren't properly
	// cleared either.
	// err := session.Save(r, w)
	// if err != nil {
	// 	panic(err)
	// }

	return map[string]interface{}{
		"full_url": r.URL.String(),
		"url":      r.URL,
		"sess":     session.Values,
		// "messages": messages,
	}
}

// func flashMessage(w http.ResponseWriter, r *http.Request, status, message string) {
// 	session := getSession(r)
// 	session.AddFlash(Message{status, message}, "messages")
// 	err := session.Save(r, w)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func getSession(r *http.Request) *sessions.Session {
	session, _ := sess.Get(r, "sessionid")

	return session
}

func getSessionString(r *http.Request, key string) string {
	val, ok := getSession(r).Values[key]
	if ok {
		return val.(string)
	}

	return ""
}
