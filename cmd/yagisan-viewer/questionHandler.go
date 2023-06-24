package main

import (
	"fmt"
	"net/http"
	"html/template"

	// "net/http/cgi"

	"github.com/go-chi/chi"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shimiwaka/yagisan/connector"
	"github.com/shimiwaka/yagisan/schema"
)

func questionHandler(w http.ResponseWriter, r *http.Request) {
	qToken := chi.URLParam(r, "qToken")
	if qToken == "" {
		// for test
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "parse error occured: %v", err)
			return
		}

		qToken = r.Form.Get("qToken")
	}

	db := connector.ConnectDB()
	defer db.Close()

	question := schema.Question{}
	err := db.First(&question, "token = ? and visible = true", qToken).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "token is invalid : %v", err)
		return
	}

	answer := schema.Answer{}
	db.Order("id desc").First(&answer, "question = ?", question.ID)

	// resp := schema.GetQuestionReponse{
	// 	Success: true,
	// 	Email: question.Email,
	// 	IP: question.IP,
	// 	UserAgent: question.UserAgent,
	// 	Body: question.Body,
	// 	QuestionID: question.ID,
	// 	AnswerBody: answer.Body,
	// 	CreatedAt: question.CreatedAt,
	// }

	t, err := template.ParseFiles("template/question.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "template error: %v", err)
		return
	}

	if err := t.Execute(w, struct {
		Token   string;
		AnswerBody string;
		Body string;
	}{
		Token:  qToken,
		Body: question.Body,
		AnswerBody: answer.Body,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to execute template: %v", err)
		return
	}
}