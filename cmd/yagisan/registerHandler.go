package main

import (
	"crypto/md5"
	"crypto/sha512"
	"fmt"
	"net/http"
	"time"

	// "github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shimiwaka/yagisan/connector"
	"github.com/shimiwaka/yagisan/schema"
)

func register(db *gorm.DB, email string, username string, password string, description string) (string, error) {
	box := schema.Box{Username: username, Email: email, Password: password, Description: description}
	db.Create(&box)

	seed := []byte(username + fmt.Sprint(time.Now().UnixNano()))
	token := fmt.Sprintf("%x", md5.Sum(seed))

	accessToken := schema.AccessToken{Box: box.ID, Token: token}
	db.Create(&accessToken)
	return fmt.Sprintf("{\"success\":true, \"token\":\"%s\"}", token), nil
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	e := r.ParseForm()
	if e != nil {
		fmt.Fprintf(w, errorMessage("parse error occured"))
		return
	}

	email := r.Form.Get("email")
	username := r.Form.Get("username")
	password := fmt.Sprintf("%x", sha512.Sum512([]byte(r.Form.Get("password"))))
	description := r.Form.Get("description")

	db := connector.ConnectDB()
	defer db.Close()

	result, err := register(db, email, username, password, description)
	if err != nil {
		fmt.Fprintf(w, errorMessage(fmt.Sprintf("%s", err)))
		return
	}
	fmt.Fprintf(w, result)
}
