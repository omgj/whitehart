package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	stripe "github.com/stripe/stripe-go"
	customer "github.com/stripe/stripe-go/customer"
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
	http.HandleFunc("/whoami", whoami)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/txtpwd", txtpwd)
	http.HandleFunc("/cardtoken", cardtoken)
	http.HandleFunc("/codeconf", codeconf)
	http.HangleFunc("/addtocart", addtocart)	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("listening on 8080")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func cardtoken(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	stripe.Key = "sk_test_XCvlmn17AXsURaJN66uYs1Mk"
	cp := &stripe.CustomerParams{
		Phone: stripe.String("0416580041"),
	}
	cp.SetSource(token)
	c, _ := customer.New(cp)
	io.WriteString(w, c.ID)
}

func logout(w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{
		Name:    "whart",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),

		HttpOnly: true,
	}
	http.SetCookie(w, c)
	w.Write([]byte(`ok`))
}

func whoami(w http.ResponseWriter, r *http.Request) {
	c, e := r.Cookie("whart")
	if e != nil {
		w.Write([]byte(`none`))
		return
	}
	sessionsecret := os.Getenv("SESHSECRET")
	q, er := fs.Collection("sessions").Doc(c.Value + sessionsecret).Get(context.Background())
	if er != nil {
		w.Write([]byte(`none`))
		return
	}
	mm := q.Data()
	user := mm["user"].(string)
	muser := "0" + user[3:]
	w.Write([]byte(muser))
}

func codeconf(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	number := r.FormValue("number")
	num := "+61" + number[1:]
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
		if (int(time.Now().Unix()) - int(dm["codevalidity"].(int64))) < 30 {
			uuids := uuid.New()
			sessionsecret := os.Getenv("SESHSECRET")
			uui := uuids.String() + sessionsecret
			_, errr := fs.Collection("sessions").Doc(uui).Set(context.Background(), map[string]interface{}{
				"sessioncreated": int(time.Now().Unix()),
				"user":           num,
			})
			if errr != nil {
				w.Write([]byte(`err`))
				return
			}
			coo := http.Cookie{Name: "whart", Value: uuids.String(), Path: "/"}
			http.SetCookie(w, &coo)
			w.Write([]byte("0" + num[3:]))
			return
		}
	}
	w.Write([]byte(`err`))
}

func txtpwd(w http.ResponseWriter, r *http.Request) {
	num := "+61" + r.FormValue("numuser")[1:]
	log.Println("confirming: ", num)
	_, er := fs.Collection("people").Doc(num).Get(context.Background())
	code := sendSms(num)
	if er != nil {
		_, e := fs.Collection("people").Doc(num).Set(context.Background(), map[string]interface{}{
			"codevalidity": int(time.Now().Unix()),
			"code":         code,
			"role":         "customer",
		})
		if e != nil {
			w.Write([]byte(`err`))
			return
		}
		w.Write([]byte(`log`))
		return
	}
	_, err := fs.Collection("people").Doc(num).Update(context.Background(), []firestore.Update{{Path: "code", Value: code}})
	if err != nil {
		w.Write([]byte(`err`))
		return
	}
	_, errr := fs.Collection("people").Doc(num).Update(context.Background(), []firestore.Update{{Path: "codevalidity", Value: int(time.Now().Unix())}})
	if errr != nil {
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
	msgData.Set("From", twilioNumber)
	msgData.Set("Body", randomnumber)
	msgDataReader := *strings.NewReader(msgData.Encode())
	client := &http.Client{}
	twilioSid := os.Getenv("TWILSID")
	twilioUrl := "https://api.twilio.com/2010-04-01/Accounts/" + twilioSid + "/Messages.json"
	req, _ := http.NewRequest("POST", twilioUrl, &msgDataReader)
	twilioAuth := os.Getenv("TWILAUTH")
	req.SetBasicAuth(twilioSid, twilioAuth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			log.Println(data["sid"])
			return randomnumber
		}
	} else {
		log.Println(resp.Status)
		return ""
	}
	return ""
}

func public(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr)
	a, _ := ioutil.ReadFile("public.html")
	io.WriteString(w, string(a))
}
