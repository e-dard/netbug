# rdb

Package `rdb` provides an `http.Handler` for accessing the profiling and debug tools available in the `/net/http/pprof` and `/runtime/pprof` packages.

The advantages of using `rdb` over the existing `/net/http/pprof` handlers are:

 1. You can register the handler under an arbitrary route-prefix. A use-case might be to have a secret endpoint for keeping this information hidden from prying eyes on production boxes;
 2. It pulls together all the handlers from `/net/http/pprof` and `/runtime/pprof` into a single index page, for when you can't quite remember the URL for the profile you want; and
 3. You can register the handlers onto `http.ServeMux`'s that aren't `http.DefaultServeMux`.

**Note**:
It still imports `/net/http/pprof`, which means the `/debug/pprof` routes in that package get registered on `http.DefaultServeMux`. If you're using this package to avoid those routes being registered, you should use it with your own `http.ServeMux`.

`rdb` is trying to cater for the situation where you want all profiling tools available remotely on your running services, but you don't want to expose the `/debug/pprof` routes that `net/http/pprof` forces you to expose.

## How do I use it?
All you have to do is give `rdb` the `http.ServeMux` you want to register the handlers on, and you're away.

```go
package main

import (
	"log"
	"net/http"

	"github.com/e-dard/remote-debug"
)

func main() {
	r := http.NewServeMux()
	rdb.Register("/some-prefix/", r)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
```

Visiting `http://localhost:8080/some-prefix/debug/pprof/` will then return:

![](http://cl.ly/image/3x1U3O1B2L3C/Screen%20Shot%202015-03-07%20at%2019.53.28.png)

### What can you do with it?

It just wraps the behaviour of the `/net/http/pprof` and `/runtime/pprof` packages.
Check out their documentation to see what's available.
As an example though, if you want to run a 30-second CPU profile on your running service it's really simple:

```
$ go tool pprof https://example.com/some-hidden-prefix/debug/pprof/profile
```

## Background
The [net/http/pprof](http://golang.org/pkg/net/http/pprof/) package is great.
It let's you access profiling and debug information about your running services, via `HTTP`, and even plugs straight into `go tool pprof`.
You can find out more about using the `net/http/pprof` package at the bottom of [this blog post](http://blog.golang.org/profiling-go-programs).

However, there are a couple of problems with the `net/http/pprof` package.

 1. It assumes you're cool about the relevant handlers being registered under the `/debug/pprof` route.
 2. It assumes you're cool about handlers being registered on `http.DefaultServeMux`.

You can sort of fix (1) and (2) by digging around the `net/http/pprof` package and registering all all the exported handlers under different paths on your own `http.ServeMux`, but you still have the problem of the index page—which is useful to visit if you don't profile much—using hard-coded paths. It doesn't quite work well.
Also, the index page doesn't provide you with easy links to the debug information that the `net/http/pprof` has handlers for.

So, `rdb` is just a simple package to fix (1) and (2).

