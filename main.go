package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", public)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func public(w http.ResponseWriter, r *http.Request) {
	a, _ := ioutil.ReadFile("public.html")
	io.WriteString(w, string(a))
}