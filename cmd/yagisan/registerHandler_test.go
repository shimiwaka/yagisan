package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/jinzhu/gorm"
	"bytes"
	"io"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shimiwaka/yagisan/connector"
	// "github.com/shimiwaka/yagisan/schema"
)

func doRegisterTest(t *testing.T, db *gorm.DB) {
    jsonParam := `{"email":"hoge@hoge.com", "username":"hogefuga", "password":"hogehoge", "description":"hogehoge"}`

	r := httptest.NewRequest(http.MethodPost, "http://example.com/register", bytes.NewBuffer([]byte(jsonParam)))
	w := httptest.NewRecorder()
	err := register(db, w, r)

	if err != nil {
		fmt.Println("Error occured")
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