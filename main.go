package main

import (
	"io"
	"fmt"
	"io/ioutil"
	"log"
	"context"
	"cloud.google.com/go/firestore"
	"net/http"
	"net/url"
	"strings"
	"os"
	"math/rand"
	"encoding/json"
	"strconv"
	"time"
)

var fs *firestore.Client

const (
	twilioNumber = "+61480019099"
	twilioSid   = "AC836708d5d4b4c2ba9a14aaa0c0f692c2"
	twilioAuth  = "f0f64339eea0feb0aa54a15f70acc66b"
	twilioUrl = "https://api.twilio.com/2010-04-01/Accounts/"+twilioSid+"/Messages.json"
)

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
	http.HandleFunc("/txtpwd", txtpwd)
	http.HandleFunc("/codeconf", codeconf)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("listening on 8080")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func codeconf(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	number := r.FormValue("number")
	p, err := fs.Collection("people").Doc("61+"+number[1:]).Get(context.Background())
	if err != nil {
		w.Write([]byte(`wrong`))
		return
	}
	person := p.Data()
	if code == person["code"].(string) && person["codevalidity"].(int)-int(time.Now().Unix()) < 60 {
		w.Write([]byte(`ok`))
		return
	}
}

func txtpwd(w http.ResponseWriter, r *http.Request) {
	num := "+61"+r.FormValue("numuser")[1:]
	code := sendSms(num)
	_, e := fs.Collection("people").Doc(num).Set(context.Background(), map[string]interface{}{
		"codevalidity": int(time.Now().Unix()),
		"code": code,
	})
	if e != nil {
		w.Write([]byte(`err`))
		return
	}
	}
	w.Write([]byte(`log`))
}

func sendSms(to string) string {
	randomnumber := strconv.Itoa(rand.Intn(10000))
	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From",twilioNumber)
	msgData.Set("Body", randomnumber)
	msgDataReader := *strings.NewReader(msgData.Encode())
	client := &http.Client{}
	req, _ := http.NewRequest("POST", twilioUrl, &msgDataReader)
	req.SetBasicAuth(twilioSid, twilioAuth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(req)
	if (resp.StatusCode >= 200 && resp.StatusCode < 300) {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if (err == nil) {
			log.Println(data["sid"])
			return randomnumber
		}
		} else {
			log.Println(resp.Status);
			return ""
		}
	return ""
}

func public(w http.ResponseWriter, r *http.Request) {
	a, _ := ioutil.ReadFile("public.html")
	io.WriteString(w, string(a))
}
