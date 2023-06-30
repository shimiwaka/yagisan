package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
    "unicode/utf8"

	"gopkg.in/yaml.v2"

	// "net/http/cgi"

	"github.com/go-chi/chi"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shimiwaka/yagisan/connector"
	"github.com/shimiwaka/yagisan/schema"
)

func questionHandler(w http.ResponseWriter, r *http.Request) {
	settings := schema.Settings{}
	b, _ := os.ReadFile("config.yaml")
	yaml.Unmarshal(b, &settings)

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

	t, err := template.ParseFiles("template/question.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "template error: %v", err)
		return
	}

	bodys := strings.Split(question.Body, "\n")
	answerBodys := strings.Split(answer.Body, "\n")

	shortAnswerBody := answer.Body
	if utf8.RuneCountInString(shortAnswerBody) > 100 {
		shortAnswerBody = string([]rune(answer.Body)[:100])
	}

	if err := t.Execute(w, struct {
		Token      string
		AnswerBody []string
		ShortAnswerBody string
		Body       []string
		MainUrl    string
		OgpImage   string
	}{
		Token:      qToken,
		AnswerBody: answerBodys,
		ShortAnswerBody : shortAnswerBody,
		Body:       bodys,
		MainUrl:    settings.MainUrl + "/mypage",
		OgpImage:   settings.SubUrl + "/image/" + qToken,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to execute template: %v", err)
		return
	}
}
