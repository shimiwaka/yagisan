package main

import (
	"net/http"

	"net/http/cgi"
	"os"

	"github.com/go-chi/chi"
)

func main() {
	rootPath := os.Getenv("SCRIPT_NAME")

	r := chi.NewRouter()
	r.Get(rootPath+"/question/{qToken}", questionHandler)
	r.Get(rootPath+"/image/{qToken}", imageHandler)

	fileServer := http.FileServer(http.Dir("./static/"))
	r.Get(rootPath+"/static/*",
		func(w http.ResponseWriter, r *http.Request) {
			http.StripPrefix(rootPath+"/static/", fileServer).ServeHTTP(w, r)
		},
	)

	// http.ListenAndServe(":9998", r)
	cgi.Serve(r)
}
