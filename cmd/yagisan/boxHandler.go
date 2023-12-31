package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

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
		return errors.New("lack of parameters")
	}

	if !regexp.MustCompile("^[0-9a-zA-Z_]+$").MatchString(username) {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("username must be only alphabet, number and _.")
	}

	if len(rawPassword) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("password must be at least 8 characters")
	}

	if len(username) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("username must be at least 3 characters")
	}

	box := schema.Box{Username: username, Email: email, Password: password, Description: description}
	err = db.Create(&box).Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

func registerHandler(w http.ResponseWriter, r *http.Request) {
	db := connector.ConnectDB()
	defer db.Close()

	err := register(db, w, r)

	if err != nil {
		fmt.Fprintf(w, "{\"success\":false,\"message\":\"%s\"}", err)
	}
}

func showBox(db *gorm.DB, w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return errors.New("parse error occured")
	}

	accessToken := r.Form.Get("accessToken")
	box := validateAccessToken(db, accessToken)
	page, _ := strconv.Atoi(r.Form.Get("page"))

	if box.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("invalid access token")
	}

	questions := []schema.Question{}
	err = db.Limit(10).Offset(page*10).Order("id desc").Find(&questions, "box = ? and visible = true", box.ID).Error

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return errors.New("error in finding questions")
	}

	resp := schema.ShowBoxResponse{
		Username:  box.Username,
		Success:   true,
		Questions: questions,
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(&resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return errors.New("failed to encode json")
	}

	fmt.Fprint(w, buf.String())
	return nil
}

func showBoxHandler(w http.ResponseWriter, r *http.Request) {
	db := connector.ConnectDB()
	defer db.Close()

	err := showBox(db, w, r)

	if err != nil {
		fmt.Fprintf(w, "{\"success\":false,\"message\":\"%s\"}", err)
	}
}

func updateBox(db *gorm.DB, w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return errors.New("parse error occured")
	}

	accessToken := r.Form.Get("accessToken")
	box := validateAccessToken(db, accessToken)

	if box.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("invalid access token")
	}

	rawPassword := r.Form.Get("password")
	password := fmt.Sprintf("%x", sha512.Sum512([]byte(rawPassword)))

	if box.Password != password {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("password is wrong")
	}

	newEmail := r.Form.Get("newEmail")
	newUsername := r.Form.Get("newUsername")
	rawNewPassword := r.Form.Get("newPassword")
	newPassword := fmt.Sprintf("%x", sha512.Sum512([]byte(rawNewPassword)))
	newDescription := r.Form.Get("newDescription")
	newSecureMode := r.Form.Get("newSecureMode")

	if !regexp.MustCompile("^[0-9a-zA-Z_]+$").MatchString(newUsername) && newUsername != "" {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("username must be only alphabet, number and _.")
	}

	if len(rawNewPassword) < 8 && rawNewPassword != "" {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("password must be at least 8 characters")
	}

	if len(newUsername) < 3 && newUsername != "" {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("username must be at least 3 characters")
	}

	if box.Email != newEmail && newEmail != "" {
		err = db.Model(&box).Update("email", newEmail).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return errors.New("DB error occured")
		}
	}

	if box.Username != newUsername && newUsername != "" {
		err = db.Model(&box).Update("username", newUsername).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return errors.New("DB error occured")
		}
	}

	if box.Password != newPassword && rawNewPassword != "" {
		err = db.Model(&box).Update("password", newPassword).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return errors.New("DB error occured")
		}
	}

	if box.Description != newDescription {
		err = db.Model(&box).Update("description", newDescription).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return errors.New("DB error occured")
		}
	}

	if newSecureMode == "true" {
		err = db.Model(&box).Update("secure_mode", true).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return errors.New("DB error occured")
		}
	}

	if newSecureMode == "false" {
		err = db.Model(&box).Update("secure_mode", false).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return errors.New("DB error occured")
		}
	}

	fmt.Fprintf(w, "{\"success\":true}")

	return nil
}

func updateBoxHandler(w http.ResponseWriter, r *http.Request) {
	db := connector.ConnectDB()
	defer db.Close()

	err := updateBox(db, w, r)

	if err != nil {
		fmt.Fprintf(w, "{\"success\":false,\"message\":\"%s\"}", err)
	}
}

func profile(db *gorm.DB, w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return errors.New("parse error occured")
	}
	username := r.Form.Get("username")
	accessToken := r.Form.Get("accessToken")

	box := schema.Box{}
	if accessToken != "" {
		box = validateAccessToken(db, accessToken)
		if box.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return errors.New("invalid access token")
		}
		username = box.Username
	} else {
		err = db.First(&box, "username = ?", username).Error
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("invalid username")
	}

	resp := schema.ProfileResponse{
		Success:     true,
		Username:    username,
		Description: box.Description,
		SecureMode:  box.SecureMode,
	}

	if accessToken != "" {
		resp.Email = box.Email
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(&resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return errors.New("failed to encode json")
	}

	fmt.Fprint(w, buf.String())
	return nil
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	db := connector.ConnectDB()
	defer db.Close()

	err := profile(db, w, r)

	if err != nil {
		fmt.Fprintf(w, "{\"success\":false,\"message\":\"%s\"}", err)
	}
}
