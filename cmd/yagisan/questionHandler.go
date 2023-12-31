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
	"net/smtp"
	"os"
	"strings"

	// "time"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shimiwaka/yagisan/connector"
	"github.com/shimiwaka/yagisan/schema"
	"gopkg.in/yaml.v2"
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

	if context == "" {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("context is blank")
	}

	if len(context) > 10000 {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("character count is over")
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

	visible := true

	if box.SecureMode {
		visible = false

		if rawEmail == "" {
			w.WriteHeader(http.StatusBadRequest)
			return errors.New("please input email")
		}	
	} else {
		email = ""
	}

	rawIP := r.RemoteAddr
	IP := fmt.Sprintf("%x", sha512.Sum512([]byte(rawIP)))

	question := schema.Question{
		Box:       box.ID,
		Email:     email,
		IP:        IP,
		UserAgent: r.UserAgent(),
		Body:      context,
		Token:     token,
		Visible:   visible,
	}
	db.Create(&question)

	testMode := os.Getenv("TEST_MODE")

	if testMode != "1" {
		settings := schema.Settings{}
		b, _ := os.ReadFile("config.yaml")
		yaml.Unmarshal(b, &settings)

		from := settings.MailAddress
		subject := "やぎさんゆうびん"

		recipients := []string{}
		body := ""
		auth := smtp.CRAMMD5Auth(settings.MailAddress, settings.MailPassword)
		msg := []byte{}

		if box.SecureMode {
			recipients = []string{rawEmail}
			body = fmt.Sprintf("以下のリンクにアクセスすると質問が送信されます：\n\n%s/confirm/%s", settings.ServiceHost, token)

			msg = []byte(strings.ReplaceAll(fmt.Sprintf("To: %s\nSubject: %s\n\n%s", strings.Join(recipients, ","), subject, body), "\n", "\r\n"))

			err = smtp.SendMail(fmt.Sprintf("%s:%d", settings.MailHost, 587), auth, from, recipients, msg)
			if err != nil {
				return err
			}
		} else {
			recipients = []string{box.Email}
			body = fmt.Sprintf("新しい質問が届きました！\n\n%s/mypage", settings.ServiceHost)

			msg = []byte(strings.ReplaceAll(fmt.Sprintf("To: %s\nSubject: %s\n\n%s", strings.Join(recipients, ","), subject, body), "\n", "\r\n"))

			err = smtp.SendMail(fmt.Sprintf("%s:%d", settings.MailHost, 587), auth, from, recipients, msg)
			if err != nil {
				return err
			}
		}
	}

	fmt.Fprintf(w, "{\"success\":true}")

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

func confirmQuestion(db *gorm.DB, w http.ResponseWriter, r *http.Request) error {
	qToken := chi.URLParam(r, "qToken")
	if qToken == "" {
		// for test
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return errors.New("parse error occured")
		}

		qToken = r.Form.Get("qToken")
	}

	question := schema.Question{}
	err := db.First(&question, "token = ?", qToken).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	db.Model(&question).Update("visible", true)
	testMode := os.Getenv("TEST_MODE")

	if testMode != "1" {
		box := schema.Box{}
		err = db.First(&box, "id = ?", question.Box).Error

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}

		settings := schema.Settings{}
		b, _ := os.ReadFile("config.yaml")
		yaml.Unmarshal(b, &settings)

		from := settings.MailAddress
		subject := "やぎさんゆうびん"

		auth := smtp.CRAMMD5Auth(settings.MailAddress, settings.MailPassword)

		recipients := []string{box.Email}
		body := fmt.Sprintf("新しい質問が届きました！\n\n%s/mypage", settings.ServiceHost)

		msg := []byte(strings.ReplaceAll(fmt.Sprintf("To: %s\nSubject: %s\n\n%s", strings.Join(recipients, ","), subject, body), "\n", "\r\n"))

		err = smtp.SendMail(fmt.Sprintf("%s:%d", settings.MailHost, 587), auth, from, recipients, msg)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(w, "{\"success\":true}")

	return nil
}

func confirmQuestionHandler(w http.ResponseWriter, r *http.Request) {
	db := connector.ConnectDB()
	defer db.Close()

	err := confirmQuestion(db, w, r)
	if err != nil {
		fmt.Fprintf(w, "{\"success\":false,\"message\":\"%s\"}", err)
	}
}

func getQuestion(db *gorm.DB, w http.ResponseWriter, r *http.Request) error {
	qToken := chi.URLParam(r, "qToken")
	if qToken == "" {
		// for test
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return errors.New("parse error occured")
		}

		qToken = r.Form.Get("qToken")
	}

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

	question := schema.Question{}
	err = db.First(&question, "token = ? and visible = true", qToken).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	if question.Box != box.ID {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("invalid access token")
	}

	answer := schema.Answer{}
	db.Order("id desc").First(&answer, "question = ?", question.ID)

	resp := schema.GetQuestionReponse{
		Success:    true,
		Email:      question.Email,
		IP:         question.IP,
		UserAgent:  question.UserAgent,
		Body:       question.Body,
		QuestionID: question.ID,
		AnswerBody: answer.Body,
		CreatedAt:  question.CreatedAt,
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

func getQuestionHandler(w http.ResponseWriter, r *http.Request) {
	db := connector.ConnectDB()
	defer db.Close()

	err := getQuestion(db, w, r)
	if err != nil {
		fmt.Fprintf(w, "{\"success\":false,\"message\":\"%s\"}", err)
	}
}
