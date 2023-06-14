package main

import (
	"fmt"
	"net/http"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Hello World!")
}
