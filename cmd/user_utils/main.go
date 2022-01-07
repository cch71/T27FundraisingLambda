package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-password/password"
)

type Auth0Creds struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
}
type Jwt struct {
	Token string `json:"access_token"`
	Type  string `json:"token_type"`
}

func getToken() Jwt {
	url := fmt.Sprint(os.Getenv("AUTH0_BASE_URL"), "oauth/token")

	creds := Auth0Creds{
		ClientId:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		Audience:     os.Getenv("AUTH0_AUDIENCE"),
		GrantType:    "client_credentials",
	}
	b, err := json.Marshal(creds)
	if err != nil {
		log.Panic("Failed to create creds: ", err)
	}
	payload := strings.NewReader(string(b))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Panic("Failed creating token request: ", err)
	}
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	log.Println(res)
	if err != nil {
		log.Panic("Failed making token request: ", err)
	}

	defer res.Body.Close()
	//body, _ := ioutil.ReadAll(res.Body)
	//log.Println(string(body))

	jwt := Jwt{}
	if err := json.NewDecoder(res.Body).Decode(&jwt); err != nil {
		log.Panic("Failed parsing jwt resp: ", err)
	}
	jwt.Token = fmt.Sprint("Bearer ", jwt.Token)
	return jwt
}

func getUsers(jwt *Jwt) {

	url := fmt.Sprint(os.Getenv("AUTH0_BASE_URL"), "api/v2/users")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panic("Failed creating request: ", url, " Err: ", err)
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", jwt.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(res)
		log.Panic("Failed making request: ", url, " Err: ", err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	log.Println(string(body))
}

var createUserJson = `
{
  "email": "***EMAIL***",
  "email_verified": true,
  "app_metadata": {"full_name": "***NAME***"},
  "nickname": "***NICKNAME***",
  "connection": "Username-Password-Authentication",
  "password": "***PASSWORD***",
  "verify_email": false
}
`

func createUser(jwt *Jwt) {
	pw, err := password.Generate(32, 5, 5, false, false)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Generated Password: ", pw)

	paramStr := strings.ReplaceAll(createUserJson, "***PASSWORD***", pw)
	log.Println(paramStr)

	rawData := json.RawMessage(paramStr)
	param, err := rawData.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprint(os.Getenv("AUTH0_BASE_URL"), "api/v2/users")
	payload := strings.NewReader(string(param))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Panic("Failed creating request: ", url, " Err: ", err)
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", jwt.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(res)
		log.Panic("Failed making request: ", url, " Err: ", err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	log.Println(string(body))
}

func main() {

	credentialsFile := path.Join(os.Getenv("HOME"), ".t27fr", "credentials")
	_ = godotenv.Load(credentialsFile)

	jwt := getToken()
	// log.Println("Token: ", jwt.Token)
	getUsers(&jwt)
	createUser(&jwt)
}
