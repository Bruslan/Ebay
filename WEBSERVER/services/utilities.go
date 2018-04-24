package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"../data"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

const capcthaURL = "https://www.google.com/recaptcha/api/siteverify"

var isSmall = regexp.MustCompile(`^[a-z]+$`).MatchString
var isCap = regexp.MustCompile(`^[A-Z]+$`).MatchString
var isNumb = regexp.MustCompile(`^[0-9]+$`).MatchString

// captcha reponse struct
type ApiCaptchaResponse struct {
	success     bool
	challengeTs time.Time
	hostname    string
	errorCodes  []int
}

// configuration struct:
type Configuration struct {
	Address      string
	AddressSSL   string
	Redirect     string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

var logger *log.Logger
var Config Configuration

// triggers parameter loading(config.json), creates logfile ***
func init() {
	loadConfig()
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file in services package", err)
	}
	logger = log.New(file, "Info ", log.Ldate|log.Ltime|log.Lshortfile)
}

// loads config.json parameters ***
func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	Config = Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
}

// verification of capctha:
func verifyCaptcha(remoteip, cap_resp string) (err error) {
	resp, err := http.PostForm(capcthaURL,
		url.Values{"secret": {"6LcBKkoUAAAAAF5UcvuWKV-7TqDXp9s1i_PAM3wn"},
			"remoteip": {remoteip}, "reponse": {cap_resp}})
	if err != nil {
		danger("HTTP post form captcha error:", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		danger("Read captcha body error:", err)
	}
	var res ApiCaptchaResponse
	err = json.Unmarshal(b, &res)
	if err != nil {
		danger("Parse json error: ", err)
	}
	if res.success {
		return nil
	}
	return err
}

// Checks if the user is logged in and has a session, if not err is not nil ***
func sessionCheck(writer http.ResponseWriter, request *http.Request) (sess data.Session, err error) {
	// check if cookie exists
	cookie, err := request.Cookie("_ianzncookie")
	if err == nil {
		// check if session valid
		device := request.Header["User-Agent"][0]
		sess = data.Session{Uuid: cookie.Value, Device: device}
		if ok := sess.SessValid(); !ok {
			err = errors.New("Session Invalid")
		}
	}
	return
}

// parse HTML templates
func parseTemplateFiles(filenames ...string) (t *template.Template) {
	var files []string
	t = template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("html/%s.html", file))
	}
	t = template.Must(t.ParseFiles(files...))
	return
}

// passses html to agent
func generateHTML(writer http.ResponseWriter, data interface{}, startfile string, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("html/%s.html", file))
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(writer, startfile, data)
}

// Convenience function to redirect to the error message page
func error_message(writer http.ResponseWriter, request *http.Request, msg string) {
	url := []string{"/err?msg=", msg}
	http.Redirect(writer, request, strings.Join(url, ""), http.StatusFound)
}

// for logging
func info(args ...interface{}) {
	logger.SetPrefix("INFO ")
	logger.Println(args...)
}

func danger(args ...interface{}) {
	logger.SetPrefix("ERROR ")
	logger.Println(args...)
}

func warning(args ...interface{}) {
	logger.SetPrefix("WARNING ")
	logger.Println(args...)
}

////////////////////////////////////////////////////////////////////
// LOGIN CHECK FUNCTIONS

func time_diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

// validate password
func check_pw(pw string) bool {

	// checks pw length
	if utf8.RuneCountInString(pw) < 8 || utf8.RuneCountInString(pw) > 20 {
		return false
	}

	// check if letter, cap letter and number in pw
	small := false
	capital := false
	numb := false
	for i := range pw {
		if isSmall(pw[i : i+1]) {
			small = true
		}
		if isCap(pw[i : i+1]) {
			capital = true
		}
		if isNumb(pw[i : i+1]) {
			numb = true
		}
	}
	if small && capital && numb {
		return true
	}
	return false
}

// validatye age
func check_age(bday time.Time) bool {

	if year, _, _, _, _, _ := time_diff(bday, time.Now()); year < 16 {
		return false
	}
	return true
}
