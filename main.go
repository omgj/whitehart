package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"context"
	"cloud.google.com/go/firestore"
	"net/http"
	"os"
)

var fs *firestore.Client

type User struct {
	Mobile string `firestore:"mobile"`
}

func init() {
	var err error
	ctx := context.Background()
	fs, err = firestore.NewClient(ctx, "whartbar")
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", public)
	http.HandleFunc("/boss", boss)
	http.HandleFunc("/login", login)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("listening on 8080")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func login(w http.ResponseWriter, r *http.Request) {
	a := r.FormValue("user")
	log.Println(a)
	if a == "" {
		a = "Empty"
	}
	peep := fs.Doc("people/when")
	_, err := peep.Set(context.Background(), User{
		Mobile: a,
	})
	if err != nil {
		panic(err)
	}

}

func boss(w http.ResponseWriter, r *http.Request) {
	a, _ := ioutil.ReadFile("boss.html")
	io.WriteString(w, string(a))
}

func public(w http.ResponseWriter, r *http.Request) {
	a, _ := ioutil.ReadFile("public.html")
	io.WriteString(w, string(a))
}
