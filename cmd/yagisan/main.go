package main

import (
	"fmt"
	"net/http"

	// "net/http/cgi"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	rootPath := os.Getenv("SCRIPT_NAME")

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"}}))

	r.Get(rootPath+"/", pingHandler)
	r.Post(rootPath+"/login", loginHandler)
	r.Post(rootPath+"/register", registerHandler)
	r.Post(rootPath+"/question", sendQuestionHandler)
	r.Post(rootPath+"/answer", sendAnswerHandler)
	r.Get(rootPath+"/confirm/{qToken}", confirmQuestionHandler)
	r.Post(rootPath+"/box/show", showBoxHandler)
	r.Post(rootPath+"/box/update", updateBoxHandler)
	r.Post(rootPath+"/question/{qToken}", getQuestionHandler)

	http.ListenAndServe(":9999", r)
	// cgi.Serve(r)
}

func errorMessage(s string) string {
	return fmt.Sprintf("{\"success\":false, \"message\":\"%s\"}", s)
}
