package main

import "fmt"
import "net/http"
import "encoding/json"
import "encoding/hex"
import "io/ioutil"
import "crypto/md5"
import "os"
import "strings"
import "net/url"
import "net/http/cookiejar"

var config Configuration
var httpClient http.Client

type Configuration struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UniqueToken string `json:"unique_token"`
	Nim string `json:"nim"`
	Subjects []string `json:"subjects"`
}

func perform_login() bool {
	fmt.Println("Trying to login")

	resp, err := httpClient.Get("https://login.itb.ac.id/cas/login?service=https%3A%2F%2Fakademik.itb.ac.id%2Flogin%2FINA")
	if err != nil {
    	fmt.Println(err.Error())
    	return false
    }

    bodyBytes, _ := ioutil.ReadAll(resp.Body)
    body := string(bodyBytes)

    pos := strings.Index(body, "<form");
    pos = strings.Index(body[pos:], "action") + 8 + pos
    cnt := 0
    for body[pos+cnt] != '"' {
    	cnt++
 	}
 	formAction := "https://login.itb.ac.id" + body[pos:(pos+cnt)]

 	pos = strings.Index(body, "<input type=\"hidden\" name=\"lt\"") + 38
 	cnt = 0
 	for ;body[pos+cnt] != '"';cnt++ {}
 	formLt := body[pos:(pos+cnt)]

 	pos = strings.Index(body, "<input type=\"hidden\" name=\"execution\"") + 45
 	cnt = 0
 	for ;body[pos+cnt] != '"';cnt++ {}
 	formExec := body[pos:(pos+cnt)]

 	formEventId := "submit"
 	formSubmit := "LOGIN"

 	var urlValues url.Values = url.Values{}
 	urlValues.Add("username", config.Username)
 	urlValues.Add("password", config.Password)
 	urlValues.Add("lt", formLt)
 	urlValues.Add("execution", formExec)
 	urlValues.Add("_eventId", formEventId)
 	urlValues.Add("submit", formSubmit)

 	resp, err = httpClient.PostForm(formAction, urlValues)
 	if err != nil {
    	fmt.Println(err.Error())
    	return false
    }
    
    return resp.StatusCode == 200
}

func main() {
	fmt.Println("Initializing")

	raw, err := ioutil.ReadFile("./config.json")
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    err = json.Unmarshal(raw, &config)
    if err != nil {
    	fmt.Println(err.Error())
    	os.Exit(1)
    }

    fmt.Println("Running script using:")
    fmt.Println("  username:", config.Username)
    passwordArr := md5.Sum([]byte(config.Password))
    fmt.Println("  md5_password:", hex.EncodeToString(passwordArr[:]))
    fmt.Println("  unique token:", config.UniqueToken)
    fmt.Println("  nim:", config.Nim)
    fmt.Println("  subjects:", config.Subjects)

    cookieJarOpt := cookiejar.Options{}
    cookieJar, err := cookiejar.New(&cookieJarOpt)
    httpClient = http.Client {Jar: cookieJar}

    if perform_login() {
    	fmt.Println("Success login")
    } else {
    	fmt.Println("Failed login")
    }
}