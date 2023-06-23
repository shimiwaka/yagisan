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

type LoginTestCase struct {
	Username         string
	Password	string
	ExpectStatus  int
	ExpectMessage string
}

func doLoginTest(t *testing.T, db *gorm.DB, tc LoginTestCase) {
	assert := assert.New(t)

	values := url.Values{}
	values.Set("username", tc.Username)
	values.Add("password", tc.Password)

	r := httptest.NewRequest(http.MethodPost, "http://example.com/login", strings.NewReader(values.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	err := login(db, w, r)
	if err != nil {
		fmt.Fprintf(w, "{\"success\":false,\"message\":\"%s\"}", err)
	}

	resp := w.Result()
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	body := string(raw)

	assert.Equal(tc.ExpectStatus, resp.StatusCode)

	if body != "" {
		r := schema.LoginResponse{}
		_ = json.Unmarshal(raw, &r)

		if tc.ExpectMessage != "" {
			assert.Equal(tc.ExpectMessage, r.Message)
		}

		if resp.StatusCode == http.StatusOK {
			token := schema.AccessToken{}
			db.First(&token, "token = ?", r.Token)
			assert.Equal(r.Token, token.Token)
		}
	}

}

func TestLogin(t *testing.T) {
	db := connector.ConnectTestDB()
	defer db.Close()

	initializeDB(db)

	box1 := schema.Box{
		Username:    "username",
		Password:    "b109f3bbbc244eb82441917ed06d618b9008dd09b3befd1b5e07394c706a8bb980b1d7785e5976ec049b46df5f1326af5a2ea6d103fd07c95385ffab0cacbc86",
		Email:       "fuga@hoge.com",
		Description: "",
	}
	db.Create(&box1)

	tcs := []LoginTestCase{
		{
			Username:        "username",
			Password:		"password",
			ExpectStatus: http.StatusOK,
		},
		{
			Username:        "",
			Password:		"password",
			ExpectStatus: http.StatusBadRequest,
			ExpectMessage: "lack of parameters",
		},
		{
			Username:        "username",
			Password:		"",
			ExpectStatus: http.StatusBadRequest,
			ExpectMessage: "lack of parameters",
		},
		{
			Username:        "username",
			Password:		"pissword",
			ExpectStatus: http.StatusBadRequest,
			ExpectMessage: "record not found",
		},
	}

	for _, tc := range tcs {
		doLoginTest(t, db, tc)
	}
}
