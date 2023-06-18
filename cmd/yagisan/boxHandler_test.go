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

type RegisterTestCase struct {
	Email         string
	UserName      string
	Password      string
	Description   string
	ExpectStatus  int
	ExpectMessage string
}

type ShowBoxTestCase struct {
	AccessToken  string
	ExpectStatus int
	ExpectBody   string
}

func doRegisterTest(t *testing.T, db *gorm.DB, tc RegisterTestCase) {
	assert := assert.New(t)

	values := url.Values{}
	values.Set("email", tc.Email)
	values.Add("username", tc.UserName)
	values.Add("password", tc.Password)
	values.Add("description", tc.Description)

	r := httptest.NewRequest(http.MethodPost, "http://example.com/register", strings.NewReader(values.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	err := register(db, w, r)
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
			token := schema.AccessToken{}
			db.First(&token, "token = ?", r.Token)
			assert.Equal(r.Token, token.Token)

			box := schema.Box{}
			db.First(&box, "id = ?", token.Box)
			assert.Equal(box.Username, tc.UserName)
		}
	}
}

func TestRegister(t *testing.T) {
	db := connector.ConnectTestDB()
	defer db.Close()

	initializeDB(db)

	tcs := []RegisterTestCase{
		{
			Email:        "hoge@hoge.com",
			UserName:     "fuga",
			Password:     "hogefuga",
			Description:  "my question box",
			ExpectStatus: http.StatusOK,
		},
		{
			Email:         "hoge2@hoge.com",
			UserName:      "fuga",
			Password:      "hogefuga",
			Description:   "my question box",
			ExpectStatus:  http.StatusBadRequest,
			ExpectMessage: "Error 1062: Duplicate entry 'fuga' for key 'username'",
		},
		{
			Email:         "hoge@hoge.com",
			UserName:      "fuga2",
			Password:      "hogefuga",
			Description:   "my question box",
			ExpectStatus:  http.StatusBadRequest,
			ExpectMessage: "Error 1062: Duplicate entry 'hoge@hoge.com' for key 'email'",
		},
		{
			Email:         "",
			UserName:      "fuga",
			Password:      "hogefuga",
			Description:   "my question box",
			ExpectStatus:  http.StatusBadRequest,
			ExpectMessage: "lack of parameters",
		},
		{
			Email:         "hoge@hoge.com",
			UserName:      "",
			Password:      "hogefuga",
			Description:   "my question box",
			ExpectStatus:  http.StatusBadRequest,
			ExpectMessage: "lack of parameters",
		},
		{
			Email:         "hoge@hoge.com",
			UserName:      "fuga",
			Password:      "",
			Description:   "my question box",
			ExpectStatus:  http.StatusBadRequest,
			ExpectMessage: "lack of parameters",
		},
		{
			Email:         "hoge@hoge.com",
			UserName:      "fuga",
			Password:      "hogefug",
			Description:   "my question box",
			ExpectStatus:  http.StatusBadRequest,
			ExpectMessage: "password must be at least 8 characters",
		},
	}

	for _, tc := range tcs {
		doRegisterTest(t, db, tc)
	}
}

func doShowBoxTest(t *testing.T, db *gorm.DB, tc ShowBoxTestCase) {
	assert := assert.New(t)

	values := url.Values{}
	values.Set("accessToken", tc.AccessToken)

	r := httptest.NewRequest(http.MethodPost, "http://example.com/box/show", strings.NewReader(values.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	err := showBox(db, w, r)
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

		if resp.StatusCode == http.StatusOK {
			assert.Equal(tc.ExpectBody, body)
		}
	}
}

func TestShowBox(t *testing.T) {
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
		Token: "DUMMY",
	}
	db.Create(&accessToken1)

	box2 := schema.Box{
		Username:    "fuga",
		Password:    "xxxxxxxxxxx",
		Email:       "fuga@fuga.com",
		Description: "",
	}
	db.Create(&box2)

	accessToken2 := schema.AccessToken{
		Box:   box2.ID,
		Token: "DUMMY2",
	}
	db.Create(&accessToken2)

	for i := 0; i < 3; i++ {
		question := schema.Question{
			Box:  box1.ID,
			Body: fmt.Sprintf("I Love U(%d).", i),
		}
		db.Create(&question)
	}

	tcs := []ShowBoxTestCase{
		{
			AccessToken:  "DUMMY",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "{\"success\":true,\"questions\":[{\"box\":1,\"email\":\"\",\"ip\":\"\",\"user_agent\":\"\",\"body\":\"I Love U(2).\",\"token\":\"\",\"visible\":false},{\"box\":1,\"email\":\"\",\"ip\":\"\",\"user_agent\":\"\",\"body\":\"I Love U(1).\",\"token\":\"\",\"visible\":false},{\"box\":1,\"email\":\"\",\"ip\":\"\",\"user_agent\":\"\",\"body\":\"I Love U(0).\",\"token\":\"\",\"visible\":false}]}\n",
		},
		{
			AccessToken:  "non exist",
			ExpectStatus: http.StatusBadRequest,
		},
		{
			AccessToken:  "DUMMY2",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "{\"success\":true,\"questions\":[{\"box\":1,\"email\":\"\",\"ip\":\"\",\"user_agent\":\"\",\"body\":\"I Love U(2).\",\"token\":\"\",\"visible\":false},{\"box\":1,\"email\":\"\",\"ip\":\"\",\"user_agent\":\"\",\"body\":\"I Love U(1).\",\"token\":\"\",\"visible\":false},{\"box\":1,\"email\":\"\",\"ip\":\"\",\"user_agent\":\"\",\"body\":\"I Love U(0).\",\"token\":\"\",\"visible\":false}]}\n",
		},
	}

	for _, tc := range tcs {
		doShowBoxTest(t, db, tc)
	}
}
