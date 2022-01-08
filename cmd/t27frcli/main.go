package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/cch71/T27FundraisingLambda/frgql"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-password/password"
)

////////////////////////////////////////////////////////////////////////////
//
func makeGqlReq(gqlFn *string) {
	query, err := os.ReadFile(*gqlFn)
	if err != nil {
		log.Panic("Failed opening file: ", *gqlFn, " Err: ", err)
	}

	if err := frgql.OpenDb(); err != nil {
		log.Panic("Failed to initialize db:", err)
	}
	defer frgql.CloseDb()

	rJSON, err := frgql.MakeGqlQuery(string(query))
	if err != nil {
		log.Panic("GraphQL Query Failed: ", err)
	}

	var unmarshalledJson interface{}
	if err := json.Unmarshal([]byte(rJSON), &unmarshalledJson); err != nil {
		log.Panic("Parsing results failed: ", err)
	}

	rJSON, err = json.MarshalIndent(unmarshalledJson, "", "\t")
	if err != nil {
		log.Panic("Indenting results failed: ", err)
	}

	log.Printf("%s", rJSON)
}

var usersGql = `
mutation {
    setUsers(users: [{
***USERS***
    }])
}
`

////////////////////////////////////////////////////////////////////////////
//
func troopMaster2Gql(fnIn *string, fnOut *string, csvFnOut *string) {

	orecs := [][]string{
		{"PatrolName", "FullName", "UserID", "Password"},
	}
	gqlUserEntries := []string{}

	csvFile, err := os.ReadFile(*fnIn)
	if err != nil {
		log.Panic("Failed opening file: ", *fnIn, " Err: ", err)
	}
	r := csv.NewReader(strings.NewReader(string(csvFile)))
	// skip first line
	if _, err := r.Read(); err != nil {
		log.Fatal(err)
	}

	recs, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, rec := range recs {
		pw, err := password.Generate(24, 5, 5, false, false)
		if err != nil {
			log.Fatal(err)
		}
		pw = strings.ReplaceAll(pw, ",", "^")
		pw = strings.ReplaceAll(pw, "l", "k")
		pw = strings.ReplaceAll(pw, "O", "r")
		pw = strings.ReplaceAll(pw, "o", "A")
		pw = strings.ReplaceAll(pw, "i", "b")
		pw = strings.ReplaceAll(pw, "I", "Q")
		pw = strings.ReplaceAll(pw, "0", "5")
		pw = strings.ReplaceAll(pw, "1", "9")
		pw = strings.ReplaceAll(pw, "\"", "*")
		pw = strings.ReplaceAll(pw, "\\", "S")

		// log.Printf("Generated Password: ", pw)

		userid := strings.ToLower(string(rec[0][0]) + rec[1])
		fullName := fmt.Sprint(rec[0], " ", rec[1])
		entry := fmt.Sprint("        id: \"", userid, "\"\n")
		entry = entry + fmt.Sprint("        password: \"", pw, "\"\n")
		entry = entry + fmt.Sprint("        group: \"", rec[3], "\"\n")
		entry = entry + fmt.Sprint("        name: \"", fullName, "\"")
		gqlUserEntries = append(gqlUserEntries, entry)
		orecs = append(orecs, []string{rec[3], fullName, userid + "@bsatroop27.us", pw})
		// fmt.Println(record)
	}
	gqlOut := strings.ReplaceAll(usersGql, "***USERS***", strings.Join(gqlUserEntries, "\n    },{\n"))
	// log.Println(gqlOut)
	os.WriteFile(*fnOut, []byte(gqlOut), 0666)

	f, err := os.Create(*csvFnOut)
	defer f.Close()

	if err != nil {

		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, rec := range orecs {
		if err := w.Write(rec); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}

}

////////////////////////////////////////////////////////////////////////////
//
type Auth0Creds struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
}

////////////////////////////////////////////////////////////////////////////
//
type Jwt struct {
	Token string `json:"access_token"`
	Type  string `json:"token_type"`
}

////////////////////////////////////////////////////////////////////////////
//
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
	return jwt
}

////////////////////////////////////////////////////////////////////////////
//
func getUsers(jwt *Jwt) {

	url := fmt.Sprint(os.Getenv("AUTH0_BASE_URL"), "api/v2/users")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panic("Failed creating request: ", url, " Err: ", err)
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Bearer "+jwt.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(res)
		log.Panic("Failed making request: ", url, " Err: ", err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	log.Println(string(body))
}

////////////////////////////////////////////////////////////////////////////
//
func main() {

	credentialsFile := path.Join(os.Getenv("HOME"), ".t27fr", "credentials")
	_ = godotenv.Load(credentialsFile)

	gqlFilenamePtr := flag.String("gql", "", "GraphGQ File")
	troopMasterInFilenamePtr := flag.String("troopmaster", "", "TroopMaster CSV user file name")
	gqlFileNameOut := flag.String("outgql", "", "Filenmae for gql output required with troopmaster flag")
	csvFileNameOut := flag.String("outcsv", "", "Filenmae for csv output required with troopmaster flag")
	doGetAdminToken := flag.Bool("useadmintoken", false, "Sets config token")
	flag.Parse()

	if *doGetAdminToken {
		jwt := getToken()
		os.Setenv("AUTH0_ADMIN_TOKEN", jwt.Token)
	}

	if len(*gqlFilenamePtr) > 0 {
		makeGqlReq(gqlFilenamePtr)
	} else if len(*troopMasterInFilenamePtr) > 0 && len(*gqlFileNameOut) > 0 && len(*csvFileNameOut) > 0 {
		troopMaster2Gql(troopMasterInFilenamePtr, gqlFileNameOut, csvFileNameOut)
	} else {
		log.Fatal("Invalid cli flag combination")
		flag.PrintDefaults()
	}

}
