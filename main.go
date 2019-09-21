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
	http.HandleFunc("/menu", menu)
	http.HandleFunc("/register", register)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("listening on 8080")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func menu(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`Helllllo`))
}

func register(w http.ResponseWriter, r *http.Request) {
	auser := r.FormValue("user")
	apword := r.FormValue("pword")
	_, err := fs.Collection("people").Doc(auser).Get(context.Background())
	if err == nil {
		w.Write([]byte(`<span class="robo txt">user already exists</span>`))
		return
	}
	wr, err := fs.Collection("people").Doc(auser).Set(context.Background(), map[string]interface{}{
		"pword": apword,
	})
	if err != nil {
		w.Write([]byte(`<span class="robo txt">unable to register</span>`))
		return
	}
	log.Println(wr.UpdateTime)
}

func login(w http.ResponseWriter, r *http.Request) {
	auser := r.FormValue("user")
	apword := r.FormValue("pword")
	dd, err := fs.Collection("people").Doc(auser).Get(context.Background())
	if err != nil {
		w.Write([]byte(`<span class="robo txt">no such user</span>`))
		return
	}
	ddd := dd.Data()
	if ddd["pword"] == apword {
		w.Write([]byte(`<span class="robo txt">correct pwd</span>`))
		return
	}
	w.Write([]byte(`<span class="robo txt">incorrect password</span>`))
}

func boss(w http.ResponseWriter, r *http.Request) {
	a, _ := ioutil.ReadFile("boss.html")
	io.WriteString(w, string(a))
}

func public(w http.ResponseWriter, r *http.Request) {
	a, _ := ioutil.ReadFile("public.html")
	io.WriteString(w, string(a))
}
