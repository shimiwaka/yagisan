package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"github.com/jinzhu/gorm"
	// "bytes"
	"io"
    "strings"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shimiwaka/yagisan/connector"
	// "github.com/shimiwaka/yagisan/schema"
)

func doRegisterTest(t *testing.T, db *gorm.DB) {
    // jsonParam := `{"email":"hoge@hoge.com", "username":"hogehogehoge", "password":"hogehoge", "description":"hogehoge"}`

    values := url.Values{}
    values.Set("email", "hoge@hoge.com")
    values.Add("username", "fuga")
    values.Add("password", "hogefuga")
    values.Add("description", "my question box")

	r := httptest.NewRequest(http.MethodPost, "http://example.com/register", strings.NewReader(values.Encode()))
    r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	err := register(db, w, r)

	if err != nil {
		fmt.Printf("%s", err)
	}

	resp := w.Result()
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	body := string(raw)
	fmt.Println(body)

	fmt.Printf("%#v", resp)
}

func TestRegister(t *testing.T) {
	db := connector.ConnectTestDB()
	defer db.Close()

	doRegisterTest(t, db)
}