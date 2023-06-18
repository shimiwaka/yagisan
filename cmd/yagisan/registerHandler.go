package main

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"net/http"
	"errors"

	// "time"

	// "github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shimiwaka/yagisan/connector"
	"github.com/shimiwaka/yagisan/schema"
)

func register(db *gorm.DB, w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return errors.New("parse error occured")
	}

	email := r.Form.Get("email")
	username := r.Form.Get("username")
	rawPassword := r.Form.Get("password")
	password := fmt.Sprintf("%x", sha512.Sum512([]byte(rawPassword)))
	description := r.Form.Get("description")

	if email == "" || username == "" || rawPassword == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(email)
		fmt.Println(username)
		fmt.Println(rawPassword)
		return errors.New("lack of parameters")
	}

	box := schema.Box{Username: username, Email: email, Password: password, Description: description}
	db.Create(&box)

	bytes := make([]byte, 64)
	rand.Read(bytes)

	token := fmt.Sprintf("%x", md5.Sum(bytes))

	accessToken := schema.AccessToken{Box: box.ID, Token: token}
	db.Create(&accessToken)
	fmt.Fprintf(w, "{\"success\":true, \"token\":\"%s\"}", token)

	return nil
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	db := connector.ConnectDB()
	defer db.Close()

	err := register(db, w, r)

	if err != nil {
		fmt.Fprintf(w, "{\"success\":false,\"message\":\"%s\"}", err)
	}
}
