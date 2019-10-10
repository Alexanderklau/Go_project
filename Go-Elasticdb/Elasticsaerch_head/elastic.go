package main

import (
	"net/http"

	rice "github.com/GeertJohan/go.rice"
)

func main() {
	http.Handle("/", http.StripPrefix("/", http.FileServer(rice.MustFindBox("elastichead").HTTPBox())))
	http.ListenAndServe(":8000", nil)
}
