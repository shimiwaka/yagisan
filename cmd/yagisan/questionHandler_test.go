package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/jinzhu/gorm"

	"encoding/json"
	"io"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shimiwaka/yagisan/connector"
	"github.com/shimiwaka/yagisan/schema"
	"github.com/stretchr/testify/assert"
)

type SendQuestionTestCase struct {
	Email         string
	Context		  string
	BoxName		  string
	ExpectStatus  int
	ExpectMessage string
}

func doSendQuestionTest(t *testing.T, db *gorm.DB, tc SendQuestionTestCase) {
	assert := assert.New(t)

	values := url.Values{}
	values.Set("email", tc.Email)
	values.Add("context", tc.Context)
	values.Add("boxname", tc.BoxName)

	r := httptest.NewRequest(http.MethodPost, "http://example.com/send", strings.NewReader(values.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	err := sendQuestion(db, w, r)
	if err != nil {
		fmt.Fprintf(w, "{\"success\":false,\"message\":\"%s\"}", err)
	}

	resp := w.Result()
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	body := string(raw)

	assert.Equal(tc.ExpectStatus, resp.StatusCode)

	if body != "" {
		r := schema.RegisterResponse{}
		_ = json.Unmarshal(raw, &r)

		if tc.ExpectMessage != "" {
			assert.Equal(tc.ExpectMessage, r.Message)
		}

		if resp.StatusCode == http.StatusOK {
			question := schema.Question{}
			db.First(&question, "token = ?", r.Token)
			assert.Equal(r.Token, question.Token)
			assert.Equal(tc.Context, question.Body)
			assert.Equal(false, question.Visible)
		}
	}
}

func TestSendQuestion(t *testing.T) {
	db := connector.ConnectTestDB()
	defer db.Close()

	initializeDB(db)

	box := schema.Box{
		Username: "hoge",
		Password: "xxxxxxxxxxx",
		Email: "hoge@hoge.com",
		Description: "",
	}
	db.Create(&box)

	tcs := []SendQuestionTestCase{
		{
			Email:        "hoge@hoge.com",
			Context:      "I love U.",
			BoxName:		  "hoge",
			ExpectStatus: http.StatusOK,
		},
		{
			Email:        "hoge@hoge.com",
			Context:      "I love U.",
			BoxName:		  "unexist",
			ExpectStatus: http.StatusBadRequest,
			ExpectMessage: "record not found",
		},
		{
			Email:        "",
			Context:      "I love U.",
			BoxName:		  "unexist",
			ExpectStatus: http.StatusBadRequest,
			ExpectMessage: "please input email",
		},
	}

	for _, tc := range tcs {
		doSendQuestionTest(t, db, tc)
	}
}