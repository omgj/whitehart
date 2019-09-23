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
	http.HandleFunc("/register", register)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("listening on 8080")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
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
		"ipaddr": r.RemoteAddr,
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
		w.Write([]byte(`uno`))
		return
	}
	ddd := dd.Data()
	if ddd["pword"] == apword {
		aa := http.Cookie{
			Name: "whitehart",
			Value: auser,
			Path: "/",
			MaxAge: 1000,
		}
		http.SetCookie(w, &aa)
		w.Write([]byte(`ok`))
		return
	}
	w.Write([]byte(`pno`))
}

func boss(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	var cuser string
	for _, c := range cookies {
		if c.Name == "whitehart" {
			cuser = c.Value
		}
	}
	if cuser != "" {
		a, _ := ioutil.ReadFile("boss.html")
		io.WriteString(w, string(a))
		return
	}
	a, _ := ioutil.ReadFile("login.html")
	io.WriteString(w, string(a))
}

func public(w http.ResponseWriter, r *http.Request) {
	a, _ := ioutil.ReadFile("public.html")
	io.WriteString(w, string(a))
}
