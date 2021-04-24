package gohttproute

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type (
	Routes []route
	ctxKey int
)

const (
	RouteKey ctxKey = iota
)

var routes Routes = Routes{}

func init() {
	routes.AddRoute("view", "GET", "/view", index)
}

func (r *Routes) AddRoute(name, method, pattern string, handler http.HandlerFunc) {
	*r = append(*r, route{
		name:    name,
		method:  method,
		regex:   regexp.MustCompile("^" + pattern + "$"),
		handler: handler,
	})
}

func (r *Routes) GetRoute(req *http.Request) (route, error) {
	url := req.URL.Path
	for _, route := range routes {
		if route.regex.MatchString(url) {
			return route, nil
		}
	}
	return route{}, errors.New("no route found for request")
}

type route struct {
	method  string
	name    string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	route, err := routes.GetRoute(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	ctx := context.WithValue(r.Context(), RouteKey, route)
	route.handler(w, r.WithContext(ctx))

	// for _, route := range routes {
	// 	matches := route.regex.FindStringSubmatch(r.URL.Path)
	// 	if len(matches) > 0 {
	// 		if r.Method != route.method {
	// 			allow = append(allow, route.method)
	// 			continue
	// 		}
	// 	}
	// 	ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
	// 	route.handler(w, r.WithContext(ctx))
	// 	return
	// }

	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "index\n")
}
