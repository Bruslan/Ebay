package services

import (
	"bytes"
	"fmt"
	"../data"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Bruse Api:
func Bruse(writer http.ResponseWriter, request *http.Request) {

	var jsonStr = []byte(`{"title":"Bruse"}`)
	

		buf := new(bytes.Buffer)
		buf.ReadFrom(request.Body)
		s := buf.String()
	fmt.Println(s)
	
	u, _ := url.ParseRequestURI("http://0.0.0.0:80/test/" + s)
	
	client := &http.Client{}
	r, _ := http.NewRequest("GET", u.String(), bytes.NewBuffer(jsonStr))
	r.Header.Add("Content-Type", "application/json")
	fmt.Println(r.FormValue("search_text"))
	resp, _ := client.Do(r)
	fmt.Println(resp.Status)

	// resp, err := http.Get("http://129.187.229.141:8080/fb")
	// if err != nil {
	// 	fmt.Println("Error calling Brus")
	// }
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	writer.Write(body)
	fmt.Println(string(body))
	//fmt.Fprintf(writer, body)
}

// redirect to HTTPS:
func Redirect(writer http.ResponseWriter, request *http.Request) {
	host := strings.Split(request.Host, ":")[0]
	http.Redirect(writer, request, "https://"+host+":443", http.StatusMovedPermanently)
}

// Terms and Conditions, Privacy, Third Party, GET /terms
func About(writer http.ResponseWriter, request *http.Request) {
	generateHTML(writer, nil, "layout", "about", "terms", "privacy", "thirdparty")
}

// GET /err?msg=
func Err(writer http.ResponseWriter, request *http.Request) {

	vals := request.URL.Query()
	//fmt.Println("Printing out request.URL.Query() values", vals)
	_, err := sessionCheck(writer, request)
	if err != nil {
		generateHTML(writer, vals.Get("msg"), "layout", "publayout", "public.navbar", "error")
	} else {
		generateHTML(writer, vals.Get("msg"), "layout", "privlayout", "private.navbar", "error")
	}
}

// GET /
func Index(writer http.ResponseWriter, request *http.Request) {

	// check if user has valid session
	_, err := sessionCheck(writer, request)
	if err != nil {
		generateHTML(writer, nil, "layout", "publayout", "public.navbar", "index","search")
	} else {
		generateHTML(writer, nil, "layout", "privlayout", "profile", "gossip", "board", "settings", "public.thread")
	}
}

// GET /signup
func Signup(writer http.ResponseWriter, request *http.Request) {

	// generateHTML(writer, nil, "login.layout", "public.navbar", "signup")
	generateHTML(writer, "6LcBKkoUAAAAADH9MpWIO7H1rDgoqcIeUDCdYZWh", "layout", "signup.layout", "signup")
}

// POST /signup
func SignupAccount(writer http.ResponseWriter, request *http.Request) {

	err := request.ParseForm()
	if err != nil {
		warning(err, "Cannot parse form /signup", err)
	}
	// passwords check
	if !check_pw(request.PostFormValue("passw1")) {
		generateHTML(writer, "Your password must be 8-20 characters long, contain letters and numbers, and must not contain spaces, special characters, or emoji.", "layout", "signup.layout", "signup.err")
		return
	}
	if request.PostFormValue("passw1") != request.PostFormValue("passw2") {
		generateHTML(writer, "Passwords Do Not Match", "layout", "signup.layout", "signup.err")
		return
	}
	// check username for '@''
	if strings.Contains(request.PostFormValue("name"), "@") {
		generateHTML(writer, "Username cannot contain '@'", "layout", "signup.layout", "signup.err")
		return
	}
	// verify age > 16 years
	m, err := time.Parse("January", request.PostFormValue("age_month"))
	if err != nil {
		log.Println("Could not parse user birthday month", err)
	}
	day, err := strconv.Atoi(request.PostFormValue("age_day"))
	if err != nil {
		log.Println("Could not convert user birthday day", err)
	}
	year, err := strconv.Atoi(request.PostFormValue("age_year"))
	if err != nil {
		log.Println("Could not convert user birthday yaer", err)
	}
	birthday := time.Date(year, m.Month(), day, 0, 0, 0, 0, time.UTC)
	if !check_age(birthday) {
		generateHTML(writer, "Sorry, you must be at least 16 years old", "layout", "signup.layout", "signup.err")
		return
	}

	// verify captcha:
	// cap_resp := request.PostFormValue("g-recaptcha-response")
	// remoteip := strings.Split(request.RemoteAddr, ":")[0]
	// fmt.Println(cap_resp)
	// err = verifyCaptcha(remoteip, cap_resp)
	// if err != nil {
	// 	fmt.Println("capctha error: ", err)
	// 	generateHTML(writer, "Captcha Error", "layout", "signup.layout", "signup.err")
	//	return
	// }

	// create user in database and check if already exists:
	user := data.User{
		UserName:  request.PostFormValue("name"),
		FirstName: request.PostFormValue("first_name"),
		LastName:  request.PostFormValue("last_name"),
		Email:     request.PostFormValue("email"),
		Country:   request.PostFormValue("country"),
		Password:  request.PostFormValue("passw1"),
		Birthday:  birthday,
	}
	if err, stmt := user.Create(); err != nil || stmt != "" {
		// could not create new user
		generateHTML(writer, stmt, "layout", "signup.layout", "signup.err")
	} else {
		// new user created
		generateHTML(writer, "Have Fun", "layout", "signup.layout", "signup.success")
	}
}

// POST /authenticate
func Authenticate(writer http.ResponseWriter, request *http.Request) {

	err := request.ParseForm()
	if err != nil {
		warning(err, "Cannot parse form /authenticate", err)
	}

	user, err := data.GetUserEmail(request.PostFormValue("unameemail"))
	if err != nil {
		info(err, "Cannot find username or email")
	}
	if user.Password == data.Encrypt(request.PostFormValue("passw1")) {

		// create session
		device := request.Header["User-Agent"][0]
		sess := data.Session{UserId: user.Id, Device: device}
		// check if inactive session of user id and device already in session table
		err := sess.InactiveExists()
		if err != nil {
			info(err, "Could not delete inactive existing session")
		}
		// create new session table entry
		sess, err = user.CreateSession(device)
		if err != nil {
			info(err, "Cannot create session")
		}
		cookie := http.Cookie{
			Name:     "_ianzncookie",
			Value:    sess.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(writer, &cookie)
	}
	http.Redirect(writer, request, "/", 302)
}

// Get /delete_account
func DelAccount(writer http.ResponseWriter, request *http.Request) {

	// check cookie uuid
	cookie, err := request.Cookie("_ianzncookie")
	if err != http.ErrNoCookie {
		info("Failed to get cookie", err)

		// get user Id in session
		sess := data.Session{Uuid: cookie.Value}
		user, err := sess.User()
		if err != nil {
			warning("Could not find user to session uuid", err)
		}

		// delete all user sessions and user
		if err = user.DeleteSessions(); err != nil {
			warning("Could not delete all session from user", err)
		}
		if err = user.Delete(); err != nil {
			warning("Could not delete User", err)
		}
	}
	http.Redirect(writer, request, "/", 302)
}

// GET /logout
func Logout(writer http.ResponseWriter, request *http.Request) {

	// check cookie uuid
	cookie, err := request.Cookie("_ianzncookie")
	if err != http.ErrNoCookie {
		info("Failed to get cookie", err)
		sess := data.Session{Uuid: cookie.Value}
		if err = sess.SetInactive(); err != nil {
			warning("Could not set Session to inactive in lougout function", err)
		}
		/*if err = sess.Delete(); err != nil {
			warning("Could not delete Session with Logout function", err)
		}*/
	}
	http.Redirect(writer, request, "/", 302)
}
