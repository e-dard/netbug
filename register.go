// Package rdb provide an `http.Handler` that can be registered on an
// `http.ServeMux` of your choice to access remote profiling.
//
// rdb provides some advantages over the `/net/http/pprof` and
// `/runtime/pprof` packages:
//	1. You can register the handler under an arbitrary route-prefix.
//	   A use-case might be to have a secret endpoint for keeping this
//	   information hidden from prying eyes on production boxes;
//	2. It pulls together all the handlers from `/net/http/pprof` and
//	   `runtime/pprof` into a single index page, for when you can't
//	   quite remember the URL for the profile you want; and
//	3. You can register the handlers onto `http.ServeMux`'s that
//	   aren't `http.DefaultServeMux`.
//
//
// The simplest integration of `rdb` looks like:
//
//	package main
//
//	import (
//		"log"
//		"net/http"
//
//		"github.com/e-dard/remote-debug"
//	)
//
//	func main() {
//		r := http.NewServeMux()
//		rdb.Register("/some-prefix/", r)
//		if err := http.ListenAndServe(":8080", r); err != nil {
//			log.Fatal(err)
//		}
//	}
//
package rdb

import (
	"net/http"
	nhpprof "net/http/pprof"
	"runtime/pprof"
	"strings"
	"text/template"
)

// Register registers all the available handlers in `/net/http/pprof` on
// the provided `http.ServeMux`, ensuring that the routes are prefixed
// with the provided prefix argument.
//
// The provided prefix needs to have a trailing slash. The full list of
// routes registered can be seen by visiting the index page.
func Register(prefix string, mux *http.ServeMux) {
	tmpInfo := struct {
		Profiles []*pprof.Profile
		Info     []string
		Prefix   string
	}{

		Profiles: pprof.Profiles(),
		Info:     []string{"cmdline", "symbol"},
		Prefix:   prefix,
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, prefix+"debug/pprof/")
		switch name {
		case "":
			// Index page.
			indexTmpl.Execute(w, tmpInfo)
		case "cmdline":
			nhpprof.Cmdline(w, r)
		case "profile":
			nhpprof.Profile(w, r)
		case "symbol":
			nhpprof.Symbol(w, r)
		default:
			// Provides access to all profiles under runtime/pprof
			nhpprof.Handler(name).ServeHTTP(w, r)
		}
	}
	mux.HandleFunc(prefix, h)
}

var indexTmpl = template.Must(template.New("index").Parse(`<html>
  <head>
    <title>{{.Prefix}}debug/pprof/</title>
  </head>
  {{.Prefix}}debug/pprof/<br>
  <br>
  <body>
    profiles:<br>
    <table>
    {{range .Profiles}}
      <tr><td align=right>{{.Count}}<td><a href="{{$.Prefix}}debug/pprof/{{.Name}}?debug=1">{{.Name}}</a>
    {{end}}
    <tr><td align=right><td><a href="{{$.Prefix}}debug/pprof/profile">CPU</a>
    </table>
    <br>
    debug information:<br>
    <table>
    {{range .Info}}
      <tr><td align=right><td><a href="{{$.Prefix}}debug/pprof/{{.}}">{{.}}</a>
    {{end}}
    <tr><td align=right><td><a href="{{.Prefix}}debug/pprof/goroutine?debug=2">full goroutine stack dump</a><br>
    <table>
  </body>
</html>`))
