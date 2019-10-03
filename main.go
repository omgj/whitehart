package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	// "context"
	// "cloud.google.com/go/firestore"
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"math/rand"
	"encoding/json"
	"os"
)

// var fs *firestore.Client

const (
	twilioNumber = "+61480019099"
	twilioSid   = "AC836708d5d4b4c2ba9a14aaa0c0f692c2"
	twilioAuth  = "f0f64339eea0feb0aa54a15f70acc66b"
	twilioUrl = "https://api.twilio.com/2010-04-01/Accounts/"+twilioSid+"/Messages.json"
)

// type User struct {
// 	Mobile string `firestore:"mobile"`
// }

// func init() {
// 	var err error
// 	ctx := context.Background()
// 	fs, err = firestore.NewClient(ctx, "whartbar")
// 	if err != nil {
// 		panic(err)
// 	}
// }

func main() {
	http.HandleFunc("/", public)
	http.HandleFunc("/boss", boss)
	http.HandleFunc("/txtpwd", txtpwd)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("listening on 8080")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func txtpwd(w http.ResponseWriter, r *http.Request) {
	num := r.FormValue("numuser")
	// log.Println(sendSms("+61"+num[1:]))
	person, err := fs.Collection("people").Doc(numuser).Get(context.Background())
	if err != nil {
		w.Write([]byte(`none`))
		return
	}
	log.Print(person.Data())
}

// func mbile(w http.ResponseWriter, r *http.Request) {
// 	mobile := r.FormValue("mbile")
// 	tform := "+61"+mobile[1:]
// 	sendSms()
// }

// func register(w http.ResponseWriter, r *http.Request) {
// 	auser := r.FormValue("user")
// 	apword := r.FormValue("pword")
// 	_, err := fs.Collection("people").Doc(auser).Get(context.Background())
// 	if err == nil {
// 		w.Write([]byte(`<span class="robo txt">user already exists</span>`))
// 		return
// 	}
// 	wr, err := fs.Collection("people").Doc(auser).Set(context.Background(), map[string]interface{}{
// 		"pword": apword,
// 		"ipaddr": r.RemoteAddr,
// 	})
// 	if err != nil {
// 		w.Write([]byte(`<span class="robo txt">unable to register</span>`))
// 		return
// 	}
// 	log.Println(wr.UpdateTime)
// }

func cookiereader(r *http.Request) string {
	cc := r.Cookies()
	for _, c := range cc {
		if c.Name == "whitehart" {
			return c.Value
		}
	}
	return ""
	// doc, err := fs.Collection("sessions").Doc(id).Get(context.Background())
	// if err != nil {
	// 	log.Println("document doesn't exist")
	// 	return ""
	// }
	
	// if doc["When"].(int) + 1000 < time.Now().Unix() {
	// 	log.Println("too old")
	// 	return ""
	// }
	
	// return doc["who"].(string)
}

func login(w http.ResponseWriter, r *http.Request) {

	if cookiereader(r) != "" {
		return
	}

	user := r.FormValue("user")


	fs.Collection("people").


	// password := r.FormValue("password")
	// raddr := strings.Split(r.RemoteAddr, ":")[0]

	// if strings.Contains(user, "04") {
	// 	spam := false
	// 	auser = auser[1:]
	// 	mobile := "+61"+auser
	// 	d, err := fs.Collection("codes").Doc(mobile).Get(context.Background())
	// 	if err != nil {
	// 		code := sendSms(mobile)
	// 		if code == "" {
	// 			w.Write([]byte(`nocode`))
	// 			return
	// 		}
	// 		_, err := fs.Collection("codes").Doc(mobile).Set(context.Background(), map[string]interface{}{
	// 			"Code": code,
	// 			"Remote": raddr,
	// 			"When": time.Now().Unix(),
	// 		});
	// 		if err != nil {
	// 			w.Write([]byte(`no`))
	// 			return
	// 		}
	// 		w.Write([]byte(`yes`))
	// 	}
	// 	dd := d.Data()
		
	// }
	log.Println(user)
	// apword := r.FormValue("pword")
	// dd, err := fs.Collection("people").Doc(auser).Get(context.Background())
	// if err != nil {
	// 	w.Write([]byte(`uno`))
	// 	return
	// }
	// ddd := dd.Data()
	// if ddd["pword"] == apword {
	// 	fs.Collection("sessions").Doc(auser)
	// 	fs.Collection("sessions").Doc(auser).Set(context.Background(), Login{
	// 		RemoteAddr: r.strings.Split(RemoteAddr, ":")[0],
	// 		LastLogin: time.Now().Unix(),

	// 	})
	// 	aa := http.Cookie{
	// 		Name: "whitehart",
	// 		Value: auser,
	// 		Path: "/",
	// 		MaxAge: 1000,
	// 	}
	// 	http.SetCookie(w, &aa)
	// 	w.Write([]byte(`ok`))
	// 	return
	// }
	// fs.Collection("sessions").Doc(auser)
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

func boss(w http.ResponseWriter, r *http.Request) {
	id := cookiereader(r)
	if id == "" {
		a, _ := ioutil.ReadFile("login.html")
		io.WriteString(w, string(a))
		return
	}
	a, _ := ioutil.ReadFile("boss.html")
	io.WriteString(w, string(a))
}

func public(w http.ResponseWriter, r *http.Request) {
	a, _ := ioutil.ReadFile("public.html")
	io.WriteString(w, string(a))
}
