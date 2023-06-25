package main

import (
	"fmt"
	"net/http"
	"bytes"
	"encoding/binary"
	// "net/http/cgi"

	"github.com/go-chi/chi"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shimiwaka/yagisan/connector"
	"github.com/shimiwaka/yagisan/schema"
	"github.com/shimiwaka/str2img"
)

func imageHandler(w http.ResponseWriter, r *http.Request) {
	qToken := chi.URLParam(r, "qToken")
	if qToken == "" {
		// for test
		err := r.ParseForm()
		if err != nil {
			w.Header().Set("Content-Type","text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "parse error occured: %v", err)
			return
		}

		qToken = r.Form.Get("qToken")
	}

	db := connector.ConnectDB()
	defer db.Close()

	question := schema.Question{}
	err := db.First(&question, "token = ? and visible = true", qToken).Error
	if err != nil {
		w.Header().Set("Content-Type","text/plain")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "token is invalid : %v", err)
		return
	}

	w.Header().Set("Content-Type","image/png")

	generator := &str2img.Generator{
		ImageHeight: 630,
		ImageWidth:  1200,
		FontSize:    40.0,
		FontFile:    "Koruri-Regular.ttf",
		ImageBytes:  &bytes.Buffer{},
	}

	err = generator.Generate(question.Body)
	if err != nil {
		w.Header().Set("Content-Type","text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "image convert error : %v", err)
		return
	}
	err = binary.Write(w, binary.BigEndian, generator.ImageBytes.Bytes())
	if err != nil {
		panic(err)
	}
}