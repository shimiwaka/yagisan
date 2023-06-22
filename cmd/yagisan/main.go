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
	r.Post(rootPath+"/register", registerHandler)
	r.Post(rootPath+"/question", sendQuestionHandler)
	r.Post(rootPath+"/answer", sendAnswerHandler)
	r.Get(rootPath+"/confirm/{qToken}", confirmQuestionHandler)
	r.Get(rootPath+"/box/show", showBoxHandler)

	// r.Post(rootPath + "/forget", forgetHandler)

	// r.Route(rootPath + "/board", func(r chi.Router) {
	// 	r.Get("/{boardToken}", showBoardHandler)
	// 	r.Post("/{boardToken}/newcolumn", addColumnHandler)
	// 	r.Get("/{boardToken}/deletecolumn/{idx}", deleteColumnHandler)
	// 	r.Get("/{boardToken}/check/{date}/{column}", checkHandler)
	// 	r.Get("/{boardToken}/uncheck/{date}/{column}", uncheckHandler)
	// 	r.Get("/{boardToken}/newpayment/{date}", newPaymentHandler)
	// 	r.Get("/{boardToken}/cancelpayment/{date}", cancelPaymentHandler)
	//   })

	http.ListenAndServe(":9999", r)
	// cgi.Serve(r)
}

func errorMessage(s string) string {
	return fmt.Sprintf("{\"success\":false, \"message\":\"%s\"}", s)
}
