package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var port string
var routes map[string]string

func init() {
	port = fmt.Sprintf(":%s", os.Getenv("HTTP_PORT"))
	err := json.Unmarshal([]byte(os.Getenv("ROUTES")), &routes)
	if err != nil {
		fmt.Printf("ROUTES: %v\n, Err: %v\n", os.Getenv("ROUTES"), err)
		os.Exit(1)
	}
	if len(routes) == 0 {
		fmt.Printf("ROUTES: %v\n, Err: %v\n", os.Getenv("ROUTES"), errors.New("no route found"))
		os.Exit(1)
	}
}

func main() {
	fmt.Printf("ReverseProxy running on %s\n", port)
	http.ListenAndServe(port, http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	for i := range routes {
		if strings.HasPrefix(r.URL.Path, i) {
			urlRoute, _ := url.Parse(routes[i])
			httputil.NewSingleHostReverseProxy(urlRoute).ServeHTTP(w, r)
		}
	}
}
