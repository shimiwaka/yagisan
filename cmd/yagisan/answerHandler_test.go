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
	"strconv"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shimiwaka/yagisan/connector"
	"github.com/shimiwaka/yagisan/schema"
	"github.com/stretchr/testify/assert"
)

type SendAnswerTestCase struct {
	Question      uint
	Body          string
	AccessToken	  string
	ExpectStatus  int
	ExpectMessage string
}

func doSendAnswerTest(t *testing.T, db *gorm.DB, tc SendAnswerTestCase) {
	assert := assert.New(t)

	values := url.Values{}
	values.Set("question", strconv.Itoa(int(tc.Question)))
	values.Add("body", tc.Body)
	values.Add("accessToken", tc.AccessToken)

	r := httptest.NewRequest(http.MethodPost, "http://example.com/answer", strings.NewReader(values.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	err := sendAnswer(db, w, r)
	if err != nil {
		fmt.Fprintf(w, "{\"success\":false,\"message\":\"%s\"}", err)
	}

	resp := w.Result()
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	body := string(raw)

	assert.Equal(tc.ExpectStatus, resp.StatusCode)
	
	if body != "" {
		r := schema.SendAnswerResponse{}
		_ = json.Unmarshal(raw, &r)

		if tc.ExpectMessage != "" {
			assert.Equal(tc.ExpectMessage, r.Message)
		}
	}
}

func TestSendAnswer(t *testing.T) {
	db := connector.ConnectTestDB()
	defer db.Close()

	initializeDB(db)

	box1 := schema.Box{
		Username:    "hoge",
		Password:    "xxxxxxxxxxx",
		Email:       "hoge@hoge.com",
		Description: "",
	}
	db.Create(&box1)

	accessToken1 := schema.AccessToken{
		Box:   box1.ID,
		Token: "DUMMYXXXX",
	}
	db.Create(&accessToken1)

	question1 := schema.Question{
		Box:     box1.ID,
		Body:    "I love U.",
		Token:   "XXXX",
		Visible: false,
	}
	db.Create(&question1)

	tcs := []SendAnswerTestCase{
		{
			Question:     question1.ID,
			AccessToken:  "DUMMYXXXX",
			Body:         "I love U too!",
			ExpectStatus: http.StatusOK,
		},
	}

	for _, tc := range tcs {
		doSendAnswerTest(t, db, tc)
	}
}
