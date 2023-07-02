package main

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"fmt"
	"net/http"

	"time"

	// "github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shimiwaka/yagisan/connector"
	"github.com/shimiwaka/yagisan/schema"
)

func login(db *gorm.DB, w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return errors.New("parse error occured")
	}

	username := r.Form.Get("username")
	rawPassword := r.Form.Get("password")
	password := fmt.Sprintf("%x", sha512.Sum512([]byte(rawPassword)))

	if username == "" || rawPassword == "" {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("lack of parameters")
	}

	fails := []schema.LoginFailLog{}
	db.Find(&fails, "username = ? and created_at >= ?", username, time.Now().Add(time.Minute * -30))

	if len(fails) > 2 {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("account is locked")
	}

	box := schema.Box{}

	err = db.First(&box, "username = ? and password = ?", username, password).Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		fail := schema.LoginFailLog{
			Username: username,
		}
		db.Create(&fail)
		return err
	}

	bytes := make([]byte, 64)
	rand.Read(bytes)

	token := fmt.Sprintf("%x", md5.Sum(bytes))

	accessToken := schema.AccessToken{Box: box.ID, Token: token}
	db.Create(&accessToken)
	fmt.Fprintf(w, "{\"success\":true, \"token\":\"%s\"}", token)

	return nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	db := connector.ConnectDB()
	defer db.Close()

	err := login(db, w, r)

	if err != nil {
		fmt.Fprintf(w, "{\"success\":false,\"message\":\"%s\"}", err)
	}
}
