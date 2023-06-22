package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	// "time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shimiwaka/yagisan/connector"
	"github.com/shimiwaka/yagisan/schema"
)

func sendAnswer(db *gorm.DB, w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return errors.New("parse error occured")
	}

	questionId, _ := strconv.Atoi(r.Form.Get("question"))
	accessToken := r.Form.Get("accessToken")
	body := r.Form.Get("body")
	box := validateAccessToken(db, accessToken)

	if box.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("invalid access token")
	}

	question := schema.Question{}
	err = db.First(&question, "ID = ?", questionId).Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("invalid question id")
	}

	if question.Box != box.ID {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("invalid access token")
	}

	answer := schema.Answer{
		Question: uint(questionId),
		Body:     body,
	}

	db.Create(&answer)
	fmt.Fprintf(w, "{\"success\":true}")

	return nil
}

func sendAnswerHandler(w http.ResponseWriter, r *http.Request) {
	db := connector.ConnectDB()
	defer db.Close()

	err := sendAnswer(db, w, r)

	if err != nil {
		fmt.Fprintf(w, "{\"success\":false,\"message\":\"%s\"}", err)
	}
}
