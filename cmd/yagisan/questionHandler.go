package main

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"fmt"
	"net/http"

	// "time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shimiwaka/yagisan/connector"
	"github.com/shimiwaka/yagisan/schema"
)

func sendQuestion(db *gorm.DB, w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return errors.New("parse error occured")
	}

	rawEmail := r.Form.Get("email")
	email := fmt.Sprintf("%x", sha512.Sum512([]byte(rawEmail)))
	context := r.Form.Get("context")
	boxName := r.Form.Get("boxname")

	if rawEmail == "" {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("please input email")
	}

	box := schema.Box{}
	err = db.First(&box, "username = ?", boxName).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	bytes := make([]byte, 16)
	rand.Read(bytes)

	token := fmt.Sprintf("%x", md5.Sum(bytes))

	question := schema.Question{
		Box: box.ID,
		Email: email,
		IP: "",
		UserAgent: "",
		Body: context,
		Token: token,
		Visible: false,
	}
	db.Create(&question)
	fmt.Fprintf(w, "{\"success\":true, \"token\":\"%s\"}", token)
	
	return nil
}

func sendQuestionHandler(w http.ResponseWriter, r *http.Request) {
	db := connector.ConnectDB()
	defer db.Close()

	err := sendQuestion(db, w, r)

	if err != nil {
		fmt.Fprintf(w, "{\"success\":false,\"message\":\"%s\"}", err)
	}
}
