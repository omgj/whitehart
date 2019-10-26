package main

import (
	"io"
	"fmt"
	"io/ioutil"
	"log"
	"context"
	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
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
	num := "+61"+number[1:]
	ctx := context.Background()
	people := fs.Collection("people")
	who := people.Doc(num)
	ds, err := who.Get(ctx)
	if err != nil {
		log.Println("couldn't find person")
		w.Write([]byte(`err`))
		return
	}
	dm := ds.Data()
	log.Print(dm)
	if code == dm["code"].(string) {
		if (int(time.Now().Unix())-int(dm["codevalidity"].(int64)))<30 {
			uuids := uuid.New()
			sessionsecret := os.Getenv("SESHSECRET")
			uui := uuids.String()+sessionsecret
			_, errr := fs.Collection("sessions").Doc(uui).Set(context.Background(), map[string]interface{}{
				"sessioncreated": int(time.Now().Unix()),
				"user": num,
			})
			if errr != nil {
				w.Write([]byte(`err`))
				return
			}
			coo := http.Cookie{Name: "whart", Value: uuids.String(), Path: "/"}
			http.SetCookie(w, &coo)
			w.Write([]byte(num))
			return
		}
	}
	w.Write([]byte(`err`))
}

func txtpwd(w http.ResponseWriter, r *http.Request) {
	num := "+61"+r.FormValue("numuser")[1:]
	log.Println("confirming: ", num)
	code := sendSms(num)
	_, e := fs.Collection("people").Doc(num).Set(context.Background(), map[string]interface{}{
		"codevalidity": int(time.Now().Unix()),
		"code": code,
	})
	if e != nil {
		w.Write([]byte(`err`))
		return
	}
	w.Write([]byte(`log`))
}

func sendSms(to string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomnumber := strconv.Itoa(r.Intn(1000000))
	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From",twilioNumber)
	msgData.Set("Body", randomnumber)
	msgDataReader := *strings.NewReader(msgData.Encode())
	client := &http.Client{}
	twilioSid := os.Getenv("TWILSID")
	twilioUrl := "https://api.twilio.com/2010-04-01/Accounts/"+twilioSid+"/Messages.json"
	req, _ := http.NewRequest("POST", twilioUrl, &msgDataReader)
	twilioAuth := os.Getenv("TWILAUTH")
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
	c, e := r.Cookie("whart")
	if e != nil {
		a, _ := ioutil.ReadFile("public.html")
		io.WriteString(w, string(a))
		return
	}
	sessionsecret := os.Getenv("SESHSECRET")
	q, er := fs.Collection("sessions").Doc(c.Value+sessionsecret).Get(context.Background())
	if er != nil {
		aaaa, _ := ioutil.ReadFile("public.html")
		io.WriteString(w, string(aaaa))
		log.Println("session not created")
	}
	mm := q.Data()
	log.Println(mm["user"].(string))
	io.WriteString(w, mm["user"].(string))
}
