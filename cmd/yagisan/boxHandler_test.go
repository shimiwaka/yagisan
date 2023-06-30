package main

import (
	"crypto/sha512"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
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
	Page         int
	ExpectStatus int
	ExpectBody   string
}

type UpdateTestCase struct {
	NewEmail       string
	NewUserName    string
	NewPassword    string
	NewDescription string
	AccessToken    string
	Password       string
	ExpectStatus   int
	ExpectMessage  string
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
		{
			Email:         "hoge@hoge.com",
			UserName:      "fu",
			Password:      "hogefuga",
			Description:   "my question box",
			ExpectStatus:  http.StatusBadRequest,
			ExpectMessage: "username must be at least 3 characters",
		},
		{
			Email:         "hoge@hoge.com",
			UserName:      "ほげお",
			Password:      "hogefuga",
			Description:   "my question box",
			ExpectStatus:  http.StatusBadRequest,
			ExpectMessage: "username must be only alphabet, number and _.",
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
	values.Set("page", fmt.Sprintf("%d", tc.Page))

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
			re := regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{2}:[0-9]{2}")
			body = re.ReplaceAllString(body, "-")
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
			Box:     box1.ID,
			Body:    fmt.Sprintf("I Love U(%d).", i),
			Visible: true,
		}
		db.Create(&question)
	}

	tcs := []ShowBoxTestCase{
		{
			AccessToken:  "DUMMY",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "{\"success\":true,\"username\":\"hoge\",\"questions\":[{\"ID\":3,\"CreatedAt\":\"-\",\"UpdatedAt\":\"-\",\"DeletedAt\":null,\"body\":\"I Love U(2).\",\"token\":\"\"},{\"ID\":2,\"CreatedAt\":\"-\",\"UpdatedAt\":\"-\",\"DeletedAt\":null,\"body\":\"I Love U(1).\",\"token\":\"\"},{\"ID\":1,\"CreatedAt\":\"-\",\"UpdatedAt\":\"-\",\"DeletedAt\":null,\"body\":\"I Love U(0).\",\"token\":\"\"}]}\n",
		},
		{
			AccessToken:  "non exist",
			ExpectStatus: http.StatusBadRequest,
		},
		{
			AccessToken:  "DUMMY2",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "{\"success\":true,\"username\":\"fuga\",\"questions\":[]}\n",
		},
		{
			AccessToken:  "DUMMY",
			ExpectStatus: http.StatusOK,
			Page:         1,
			ExpectBody:   "{\"success\":true,\"username\":\"hoge\",\"questions\":[]}\n",
		},
	}

	for _, tc := range tcs {
		doShowBoxTest(t, db, tc)
	}
}

func doUpdateTest(t *testing.T, db *gorm.DB, tc UpdateTestCase) {
	assert := assert.New(t)

	values := url.Values{}
	values.Set("accessToken", tc.AccessToken)
	values.Set("password", tc.Password)
	values.Set("newEmail", tc.NewEmail)
	values.Set("newUsername", tc.NewUserName)
	values.Set("newDescription", tc.NewDescription)
	values.Set("newPassword", tc.NewPassword)

	r := httptest.NewRequest(http.MethodPost, "http://example.com/box/update", strings.NewReader(values.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	err := updateBox(db, w, r)
	if err != nil {
		fmt.Fprintf(w, "{\"success\":false,\"message\":\"%s\"}", err)
	}

	resp := w.Result()
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	body := string(raw)

	assert.Equal(tc.ExpectStatus, resp.StatusCode)

	if body != "" {
		r := schema.UpdateResponse{}
		_ = json.Unmarshal(raw, &r)

		if resp.StatusCode == http.StatusOK {
			accessToken := schema.AccessToken{}
			db.First(&accessToken, "token = ?", tc.AccessToken)

			box := schema.Box{}
			db.First(&box, "id = ?", accessToken.Box)

			if tc.NewEmail != "" {
				assert.Equal(tc.NewEmail, box.Email)
			}

			if tc.NewUserName != "" {
				assert.Equal(tc.NewUserName, box.Username)
			}

			if tc.NewPassword != "" {
				assert.Equal(fmt.Sprintf("%x", sha512.Sum512([]byte(tc.NewPassword))), box.Password)
			}

			assert.Equal(tc.NewDescription, box.Description)
		} else {
			assert.Equal(tc.ExpectMessage, r.Message)
		}
	}
}

func TestUpdate(t *testing.T) {
	db := connector.ConnectTestDB()
	defer db.Close()

	initializeDB(db)

	box1 := schema.Box{
		Username:    "hoge",
		Password:    fmt.Sprintf("%x", sha512.Sum512([]byte("xxxxxxxx"))),
		Email:       "hoge@hoge.com",
		Description: "",
	}
	db.Create(&box1)

	box2 := schema.Box{
		Username:    "hoge2",
		Password:    fmt.Sprintf("%x", sha512.Sum512([]byte("xxxxxxxx"))),
		Email:       "fuga@hoge.com",
		Description: "",
	}
	db.Create(&box2)

	accessToken := schema.AccessToken{
		Box:   box1.ID,
		Token: "DUMMY",
	}
	db.Create(&accessToken)

	tcs := []UpdateTestCase{
		{
			AccessToken:  "DUMMY",
			NewEmail:     "new@email.com",
			Password:     "xxxxxxxx",
			ExpectStatus: http.StatusOK,
		},
		{
			AccessToken:  "DUMMY",
			NewUserName:  "fuga",
			Password:     "xxxxxxxx",
			ExpectStatus: http.StatusOK,
		},
		{
			AccessToken:    "DUMMY",
			NewDescription: "I Love U.",
			Password:       "xxxxxxxx",
			ExpectStatus:   http.StatusOK,
		},
		{
			AccessToken:  "DUMMY",
			NewPassword:  "yyyyyyyy",
			Password:     "xxxxxxxx",
			ExpectStatus: http.StatusOK,
		},
		{
			AccessToken:    "DUMMY",
			NewEmail:       "new2@email.com",
			NewUserName:    "piyo",
			NewDescription: "I Love U too.",
			NewPassword:    "zzzzzzzz",
			Password:       "yyyyyyyy",
			ExpectStatus:   http.StatusOK,
		},
		{
			AccessToken:   "DUMMY!!!!!!!",
			Password:      "zzzzzzzz",
			ExpectStatus:  http.StatusBadRequest,
			ExpectMessage: "invalid access token",
		},
		{
			AccessToken:   "DUMMY",
			Password:      "aaaaaaa",
			ExpectStatus:  http.StatusBadRequest,
			ExpectMessage: "password is wrong",
		},
		{
			AccessToken:   "DUMMY",
			NewEmail:      "fuga@hoge.com",
			Password:      "zzzzzzzz",
			ExpectStatus:  http.StatusInternalServerError,
			ExpectMessage: "DB error occured",
		},
	}

	for _, tc := range tcs {
		doUpdateTest(t, db, tc)
	}
}

func TestProfile(t *testing.T) {
	db := connector.ConnectTestDB()
	defer db.Close()

	initializeDB(db)

	box := schema.Box{
		Username:    "hoge",
		Email:       "hoge@hoge.com",
		Description: "Please give me questions.",
		SecureMode:  true,
	}
	db.Create(&box)

	accessToken := schema.AccessToken {
		Box: box.ID,
		Token: "DUMMY",
	}
	db.Create(&accessToken)

	assert := assert.New(t)

	values := url.Values{}
	values.Set("username", "hoge")

	r := httptest.NewRequest(http.MethodPost, "http://example.com/box/update", strings.NewReader(values.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	err := profile(db, w, r)
	if err != nil {
		fmt.Fprintf(w, "{\"success\":false,\"message\":\"%s\"}", err)
	}

	resp := w.Result()
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	body := string(raw)

	assert.Equal(resp.StatusCode, http.StatusOK)
	assert.Equal("{\"success\":true,\"username\":\"hoge\",\"email\":\"\",\"description\":\"Please give me questions.\",\"SecureMode\":true,\"message\":\"\"}\n", body)

	values = url.Values{}
	values.Set("username", "hoge")
	values.Set("accessToken", "DUMMY")

	r = httptest.NewRequest(http.MethodPost, "http://example.com/box/update", strings.NewReader(values.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()
	err = profile(db, w, r)
	if err != nil {
		fmt.Fprintf(w, "{\"success\":false,\"message\":\"%s\"}", err)
	}

	resp = w.Result()
	raw, _ = io.ReadAll(resp.Body)
	body = string(raw)

	assert.Equal(resp.StatusCode, http.StatusOK)
	assert.Equal("{\"success\":true,\"username\":\"hoge\",\"email\":\"hoge@hoge.com\",\"description\":\"Please give me questions.\",\"SecureMode\":true,\"message\":\"\"}\n", body)

	values = url.Values{}
	values.Set("username", "fuga")

	r = httptest.NewRequest(http.MethodPost, "http://example.com/box/update", strings.NewReader(values.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()
	err = profile(db, w, r)
	if err != nil {
		fmt.Fprintf(w, "{\"success\":false,\"message\":\"%s\"}", err)
	}

	resp = w.Result()

	assert.Equal(resp.StatusCode, http.StatusBadRequest)
}
