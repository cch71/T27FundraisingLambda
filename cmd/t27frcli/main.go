package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Nerzal/gocloak/v12"
	"github.com/cch71/T27FundraisingLambda/frgql"
	"github.com/joho/godotenv"
)

var (
	_ = godotenv.Load(path.Join(os.Getenv("HOME"), ".t27fr", "credentials"))
)

////////////////////////////////////////////////////////////////////////////
//
func loginKcAdmin(ctx context.Context) (*gocloak.GoCloak, token) {
	client := gocloak.NewClient(os.Getenv("KEYCLOAK_URL"))
	kcId := os.Getenv("KEYCLOAK_ID")
	kcSecret := os.Getenv("KEYCLOAK_SECRET")
	realm := os.Getenv("KEYCLOAK_REALM")

	log.Printf("ID: %s, Secret: %s, Ream: %s", kcId, kcSecret, realm)

	token, err := client.LoginAdmin(ctx, kcId, kcSecret, realm)
	if err != nil {
		log.Fatalln("Login failed:", err.Error())
	}
	return client, token.AccessToken
}

////////////////////////////////////////////////////////////////////////////
//
func createKcUser(client *gocloak.GoCloak, token *string, user *UserInfo, password *string) error {

	if err := frgql.verifyAdminTokenFromCtx(ctx); err != nil {
		return false, err
	}

	kcUser := gocloak.User{
		FirstName:     gocloak.StringP(user.FirstName),
		LastName:      gocloak.StringP(user.LastName),
		Email:         gocloak.StringP(user.Id + "@bsatroop27.us"),
		EmailVerified: gocloak.BoolP(true),
		Enabled:       gocloak.BoolP(true),
		Username:      gocloak.StringP(user.Id),
		Groups:        &[]string{"FrSellers"},
	}

	if user.Group == "Admin" {
		kcUser.Groups = &[]string{"FrAdmins"}
	}

	userId, err := client.CreateUser(ctx, *token, realm, user)
	if err != nil {
		log.Fatalln("Oh no!, failed to create user: ", err.Error())
	}
	log.Printf("Created UserID: %s", userId)
	err = client.SetPassword(ctx, token.AccessToken, userId, realm, "YUv4m*6NZqynEJ@aPMMWQkEk", false)
	if err != nil {
		log.Fatalln("Oh no!, failed to set user password: ", err.Error())
	}
	log.Printf("Password for UserID: %s set", userId)

	return nil
}

////////////////////////////////////////////////////////////////////////////
//
func createKcUsers(ctx context.Context, users *[]frgql.UserInfo) {

	client, token := loginKcAdmin(ctx)
}

////////////////////////////////////////////////////////////////////////////
//
func syncKcUsers(ctx context.Context, db *string) {
	// syncUsersWithoutAuthCreds and addToLocalDb
	// generatePasswords and savePwToLocalDb
	// loginKcAdmin
	// for each user
	//  createKcUser
	//  setUsersAuthCredsTrue in localDb
	//  setUsersAuthCredsTrue in remoteDb Transaction
	// commit Transaction
}

// ////////////////////////////////////////////////////////////////////////////
// //
// func userInfo2Gql(users *[]frgql.UserInfo, fnOut *string, csvFnOut *string) {
//      var usersGql = `
//      mutation {
//          addUsers(users: [{
//      ***USERS***
//          }])
//      }
//      `
// 	orecs := [][]string{
// 		{"PatrolName", "FullName", "UserID", "Password"},
// 	}
// 	gqlUserEntries := []string{}
//
// 	for _, user := range *users {
// 		pw, err := password.Generate(24, 5, 5, false, false)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		pw = strings.ReplaceAll(pw, ",", "^")
// 		pw = strings.ReplaceAll(pw, "l", "k")
// 		pw = strings.ReplaceAll(pw, "O", "r")
// 		pw = strings.ReplaceAll(pw, "o", "A")
// 		pw = strings.ReplaceAll(pw, "i", "b")
// 		pw = strings.ReplaceAll(pw, "I", "Q")
// 		pw = strings.ReplaceAll(pw, "0", "5")
// 		pw = strings.ReplaceAll(pw, "1", "9")
// 		pw = strings.ReplaceAll(pw, "\"", "*")
// 		pw = strings.ReplaceAll(pw, "\\", "S")
//
// 		// log.Printf("Generated Password: ", pw)
//
// 		userid := user.Id
// 		fullName := user.Name
// 		entry := fmt.Sprint("        id: \"", userid, "\"\n")
// 		entry = entry + fmt.Sprint("        password: \"", pw, "\"\n")
// 		entry = entry + fmt.Sprint("        group: \"", user.Group, "\"\n")
// 		entry = entry + fmt.Sprint("        firstName: \"", user.FirstName, "\"")
// 		entry = entry + fmt.Sprint("        lastName: \"", user.LastName, "\"")
// 		gqlUserEntries = append(gqlUserEntries, entry)
// 		orecs = append(orecs, []string{user.Group, fullName, userid + "@bsatroop27.us", pw})
// 		// fmt.Println(record)
// 	}
// 	gqlOut := strings.ReplaceAll(usersGql, "***USERS***", strings.Join(gqlUserEntries, "\n    },{\n"))
// 	// log.Println(gqlOut)
// 	os.WriteFile(*fnOut, []byte(gqlOut), 0666)
//
// 	f, err := os.Create(*csvFnOut)
// 	if err != nil {
// 		log.Fatalln("failed to open file", err)
// 	}
// 	defer f.Close()
//
// 	w := csv.NewWriter(f)
// 	defer w.Flush()
//
// 	for _, rec := range orecs {
// 		if err := w.Write(rec); err != nil {
// 			log.Fatalln("error writing record to file", err)
// 		}
// 	}
//
// }
//
// ////////////////////////////////////////////////////////////////////////////
// //
// func xlsx2Gql(fnIn *string, fnOut *string, csvFnOut *string) {
// 	f, err := excelize.OpenFile(*fnIn)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	defer f.Close()
// 	rows, err := f.GetRows("Final Roster")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	users := []frgql.UserInfo{}
// 	for idx, row := range rows {
// 		log.Println(row)
// 		if strings.HasPrefix(row[0], "New Scouts joined") || 0 == idx {
// 			continue
// 		}
// 		userid := strings.ToLower(string(row[0][0]) + row[1])
// 		fullName := fmt.Sprint(row[0], " ", row[1])
// 		users = append(users, frgql.UserInfo{Name: fullName, Id: userid, Group: row[3]})
// 	}
// 	userInfo2Gql(&users, fnOut, csvFnOut)
// }
//
// ////////////////////////////////////////////////////////////////////////////
// //
// func troopMaster2Gql(fnIn *string, fnOut *string, csvFnOut *string) {
// 	// Troop Master is Last Name, First Name, Patrol
//
// 	csvFile, err := os.ReadFile(*fnIn)
// 	if err != nil {
// 		log.Panic("Failed opening file: ", *fnIn, " Err: ", err)
// 	}
// 	r := csv.NewReader(strings.NewReader(string(csvFile)))
// 	// skip first line
// 	if _, err := r.Read(); err != nil {
// 		log.Fatal(err)
// 	}
//
// 	recs, err := r.ReadAll()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	users := []frgql.UserInfo{}
//
// 	for _, rec := range recs {
// 		userid := strings.ToLower(string(rec[0][0]) + rec[1])
// 		fullName := fmt.Sprint(rec[0], " ", rec[1])
// 		users = append(users, frgql.UserInfo{Name: fullName, Id: userid, Group: rec[3]})
// 	}
// 	userInfo2Gql(&users, fnOut, csvFnOut)
//
// }

////////////////////////////////////////////////////////////////////////////
//
func makeGqlReq(ctx context.Context, gqlFn *string) {

	_, token := loginKcAdmin(ctx)
	ctx = context.WithValue(ctx, "T27FrAuthorization", token)

	query, err := os.ReadFile(*gqlFn)
	if err != nil {
		log.Panic("Failed opening file: ", *gqlFn, " Err: ", err)
	}

	if err := frgql.OpenDb(); err != nil {
		log.Panic("Failed to initialize db:", err)
	}
	defer frgql.CloseDb()

	rJSON, err := frgql.MakeGqlQuery(ctx, string(query))
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

	log.Printf("JSON Resp:\n%s", rJSON)
}

////////////////////////////////////////////////////////////////////////////
// usage:
//  go run main.go gql --in <gql filename>
func main() {
	ctx := context.Background()

	gqlCmd := flag.NewFlagSet("gql", flag.ExitOnError)
	gqlCmdFilenameInPtr := flag.String("in", "", "GraphGQ File")

	syncKcUsersCmd := flag.NewFlagSet("synckcusers", flag.ExitOnError)
	syncKcUsersLocalDb := syncKcUsersCmd.String("db", "", "Local DB of Users")
	if len(os.Args) < 2 {
		fmt.Println("expected 'gql' or 'synckcusers' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "synckcusers":
		syncKcUsersCmd.Parse(os.Args[2:])
		syncKcUsers(ctx, syncKcUsersLocalDb)
	case "gql":
		gqlCmd.Parse(os.Args[2:])
		if 0 >= len(*gqlCmdFilenameInPtr) {
			log.Panic("in param required for gql request")
		}
		makeGqlReq(ctx, gqlCmdFilenameInPtr)
	default:
		flag.PrintDefaults()
		log.Fatal("Invalid cli flag combination")
		os.Exit(1)
	}

	log.Println("Done")
}
