package main

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/cch71/T27FundraisingLambda/frgql"
	"github.com/google/uuid"
	"github.com/sethvargo/go-password/password"
)

// //////////////////////////////////////////////////////////////////////////
func LoginKcAdmin(ctx context.Context) (*gocloak.GoCloak, string) {
	client := gocloak.NewClient(os.Getenv("KEYCLOAK_URL"))
	kcId := os.Getenv("KEYCLOAK_ID")
	kcSecret := os.Getenv("KEYCLOAK_SECRET")
	realm := os.Getenv("KEYCLOAK_REALM")

	// log.Printf("ID: %s, Secret: %s, Ream: %s", kcId, kcSecret, realm)

	token, err := client.LoginAdmin(ctx, kcId, kcSecret, realm)
	if err != nil {
		log.Fatalln("Login failed:", err.Error())
	}
	return client, token.AccessToken
}

// ///////////////////////////////////////////////////////////////////////
func createKcUser(ctx context.Context, client *gocloak.GoCloak, token *string, user UserInfo) {
	log.Printf("Creating Kc Users:\n%v", user)

	kcUser := gocloak.User{
		FirstName:     gocloak.StringP(user.FirstName),
		LastName:      gocloak.StringP(user.LastName),
		Email:         gocloak.StringP(user.Id + "@bsatroop27.us"),
		EmailVerified: gocloak.BoolP(true),
		Enabled:       gocloak.BoolP(true),
		Username:      gocloak.StringP(user.Id),
		Groups:        &[]string{"FrSellers"},
	}

	if user.Group == "Admins" {
		kcUser.Groups = &[]string{"FrAdmins"}
	}

	realm := os.Getenv("KEYCLOAK_REALM")
	userId, err := client.CreateUser(ctx, *token, realm, kcUser)
	if err != nil {
		log.Fatalln("Oh no!, failed to create user: ", err.Error())
	}
	log.Printf("Created UserID: %s: %s", user.Id, userId)
	err = client.SetPassword(ctx, *token, userId, realm, user.Password, false)
	if err != nil {
		log.Fatalln("Oh no!, failed to set user password: ", err.Error())
	}
	log.Printf("Password for UserID: %s: %s set", user.Id, userId)
}

// ///////////////////////////////////////////////////////////////////////
func createKcUsers(ctx context.Context, client *gocloak.GoCloak, token *string, users []UserInfo) {
	for _, user := range users {
		createKcUser(ctx, client, token, user)
	}
}

// //////////////////////////////////////////////////////////////////////////
func getKcUsers(ctx context.Context, client *gocloak.GoCloak, token *string) map[string]struct{} {
	users := make(map[string]struct{})
	realm := os.Getenv("KEYCLOAK_REALM")

	kcusers, err := (*client).GetUsers(ctx, *token, realm, gocloak.GetUsersParams{})
	if err != nil {
		log.Panic("Get KC Users Query Failed: ", err)
	}

	// This is used to assign a essentially empty value to the map
	// to make it a SET type
	setExists := struct{}{}

	for _, user := range kcusers {
		users[*user.Username] = setExists
	}
	// log.Printf("KcUsers:\n%v", kcusers)
	return users
}

// //////////////////////////////////////////////////////////////////////////
var GET_NON_AUTH_USERS_GQl = `{
  users(showOnlyUsersWithoutAuthCreds: true) {
    id
    group
    firstName
    lastName
  }
}`

// //////////////////////////////////////////////////////////////////////////
type UserInfo struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Id          string `json:"id"`
	Group       string `json:"group"`
	Password    string `json:"password,omitempty"`
	CreatedTime string `json:"createdTime,omitempty"`
}

// //////////////////////////////////////////////////////////////////////////
type GetUsersResp struct {
	Data struct {
		Users []UserInfo `json:"users"`
	} `json:"data"`
}

// //////////////////////////////////////////////////////////////////////////
func getUsersWithoutAuthCreds(ctx context.Context) []UserInfo {
	rJSON, err := frgql.MakeGqlQuery(ctx, string(GET_NON_AUTH_USERS_GQl))
	if err != nil {
		log.Panic("Get Users GraphQL Query Failed: ", err)
	}

	resp := GetUsersResp{}
	if err := json.Unmarshal([]byte(rJSON), &resp); err != nil {
		log.Panic("Parsing results failed: ", err)
	}

	log.Printf("Get Users Resp:\n%#v ", resp)
	return resp.Data.Users
}

// //////////////////////////////////////////////////////////////////////////
func getUsersBackupDbEncKey() []byte {
	// op item get r7gfsypfbqr3zlim3pqso2r4a4 --field password
	cmd := exec.Command("op", "item", "get", "r7gfsypfbqr3zlim3pqso2r4a4", "--field", "password")
	// cmd.Stdin = strings.NewReader("and old falcon")
	var enckey bytes.Buffer
	cmd.Stdout = &enckey
	err := cmd.Run()
	if err != nil {
		log.Fatal("Failed running op cli", err)
	}
	return enckey.Bytes()[:32]
}

// //////////////////////////////////////////////////////////////////////////
func generatePassword() string {
	pw, err := password.Generate(24, 5, 5, false, false)
	if err != nil {
		log.Fatal("Failed to generate unique password", err)
	}
	pw = strings.ReplaceAll(pw, ",", "^")
	pw = strings.ReplaceAll(pw, "l", "k")
	pw = strings.ReplaceAll(pw, "|", "P")
	pw = strings.ReplaceAll(pw, "O", "r")
	pw = strings.ReplaceAll(pw, "o", "A")
	pw = strings.ReplaceAll(pw, "i", "b")
	pw = strings.ReplaceAll(pw, "I", "Q")
	pw = strings.ReplaceAll(pw, "0", "5")
	pw = strings.ReplaceAll(pw, "1", "9")
	pw = strings.ReplaceAll(pw, "\"", "*")
	pw = strings.ReplaceAll(pw, "\\", "S")
	log.Println("Generated Password: ", pw)
	return pw
}

// //////////////////////////////////////////////////////////////////////////
func encrypt(key, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

// //////////////////////////////////////////////////////////////////////////
func decrypt(key, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// //////////////////////////////////////////////////////////////////////////
func readUsersFromBackupDb(dbdir *string) map[string]UserInfo {
	users := make(map[string]UserInfo)

	// Get Backup DB Encryption Key
	enckey := getUsersBackupDbEncKey()

	buDbFile2UserMap := func(fname string) {
		log.Println("Reading Backup DB: ", fname)
		buDbDataIn, err := os.ReadFile(fname)
		if err != nil {
			log.Panic("Failed opening file: ", fname, " Err: ", err)
		}

		buDbDataIn, err = decrypt(enckey, buDbDataIn)
		if err != nil {
			log.Panic("Failed decrypting file: ", fname, " Err: ", err)
		}

		r := csv.NewReader(strings.NewReader(string(buDbDataIn)))
		// skip first line
		if _, err := r.Read(); err != nil {
			log.Fatal(err)
		}

		recs, err := r.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		for _, rec := range recs {
			userInfo := UserInfo{
				Id:          rec[0],
				LastName:    rec[1],
				FirstName:   rec[2],
				Group:       rec[3],
				Password:    rec[4],
				CreatedTime: rec[5],
			}
			users[userInfo.Id] = userInfo
		}
	}

	// Get DB Backup Directory
	log.Println("DB Backup Dir: ", *dbdir)

	files, err := os.ReadDir(*dbdir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".enc") {
			fname := fmt.Sprintf("%s/%s", *dbdir, file.Name())
			buDbFile2UserMap(fname)
		}
	}
	return users
}

// //////////////////////////////////////////////////////////////////////////
func userInfo2BuDbBytes(users []UserInfo) []byte {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	// Write Header
	if err := w.Write([]string{"uid", "last_name", "first_name", "group", "password", "created"}); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	for _, u := range users {
		rec := []string{u.Id, u.LastName, u.FirstName, u.Group, u.Password, u.CreatedTime}
		if err := w.Write(rec); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatalln("error flushing backup db", err)
	}
	return buf.Bytes()
}

// //////////////////////////////////////////////////////////////////////////
func saveUsersToBackupDb(dbdir *string, users []UserInfo) {
	fnameUuid, err := uuid.NewRandom()
	if err != nil {
		log.Fatal("Failed to generate uuid", err)
	}
	backupFn := fmt.Sprintf("%s/%d-%s.csv.enc", *dbdir, time.Now().Unix(), fnameUuid.String())
	log.Println("Saving BackupDB Fn: ", backupFn)

	buDbDataOut := userInfo2BuDbBytes(users)

	// Get Backup DB Encryption Key
	enckey := getUsersBackupDbEncKey()
	buDbDataOut, err = encrypt(enckey, buDbDataOut)
	if err != nil {
		log.Panic("Failed encrypting backup db file Err: ", err)
	}

	// Save File with buf.Bytes()
	f, err := os.Create(backupFn)
	if err != nil {
		log.Fatalln("failed to open backup file", err)
	}
	_, err = f.Write(buDbDataOut)
	if err != nil {
		log.Fatalln("failed to write data to file", err)
	}
	defer f.Close()
}

// //////////////////////////////////////////////////////////////////////////
func saveFullUsersTo1Password(userMap map[string]UserInfo) {
	// Convert to list of users so it can be sorted by creation time
	users := []UserInfo{}
	for _, user := range userMap {
		users = append(users, user)
	}
	sort.SliceStable(users, func(i, j int) bool {
		return users[i].CreatedTime < users[j].CreatedTime
	})

	buDbData := userInfo2BuDbBytes(users)

	// op document get u2dica6qcnv2aopudslkqrq6bq
	cmd := exec.Command("op", "document", "edit", "u2dica6qcnv2aopudslkqrq6bq", "--file-name", "FrUsersPasswords.csv", "-")
	cmd.Stdin = bytes.NewReader(buDbData)
	err := cmd.Run()
	if err != nil {
		log.Fatal("Failed running op cli", err)
	}
}

var usersGql = `
mutation {
  addOrUpdateUsers(users: [{
***USERS***
  }])
}`

// //////////////////////////////////////////////////////////////////////////
func updateFrDbUsersWithAuthCreds(ctx context.Context, users *[]UserInfo) {
	gqlUserEntries := []string{}

	for _, user := range *users {
		entry := fmt.Sprintf("        id: \"%s\"\n", user.Id)
		entry = entry + "        hasAuthCreds: true"
		gqlUserEntries = append(gqlUserEntries, entry)
	}
	gqlOut := strings.ReplaceAll(usersGql, "***USERS***", strings.Join(gqlUserEntries, "\n    },{\n"))
	log.Println(gqlOut)

	// Make GQL Query
	rJSON, err := frgql.MakeGqlQuery(ctx, gqlOut)
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

// //////////////////////////////////////////////////////////////////////////
func SyncKcUsers(ctx context.Context, dbdir *string) {
	nowTime := time.Now().UTC().Format(time.RFC3339)

	// Initialize Database Connection and Keycloak token
	if err := frgql.OpenDb(); err != nil {
		log.Panic("Failed to initialize db:", err)
	}
	defer frgql.CloseDb()

	client, token := LoginKcAdmin(ctx)
	ctx = context.WithValue(ctx, "T27FrAuthorization", token)

	// Get Users where the systems thinks they don't have Auth Creds
	users := getUsersWithoutAuthCreds(ctx)
	log.Printf("Users:\n%v", users)
	if len(users) == 0 {
		log.Println("There are no users in the system without credentials")
		return
	}

	// Read the existing records from BackupDB
	buDbUsers := readUsersFromBackupDb(dbdir)
	log.Printf("DbUsers:\n%v", buDbUsers)

	// Read the existing records from KeyCloak
	kcUsers := getKcUsers(ctx, client, &token)
	log.Printf("KcUsers:\n%v", kcUsers)

	toKcUsers := []UserInfo{}
	toBuDbUsers := []UserInfo{}
	for idx, user := range users {
		// Check if it is in KeyCloak
		_, prs := kcUsers[user.Id]
		if prs {
			log.Printf("User \"%s\" already exists in KeyCloak", user.Id)
			// If it is in KeyCloak but not in Backup DB then not sure we
			// care so continue
			continue
		}
		// Check if it is in the backup db
		buDbUserInfo, prs := buDbUsers[user.Id]
		if prs {
			log.Printf("User \"%s\" already exists in BackupDB but not KeyCloak", user.Id)
			// Already exists in BuDB but not keycloak so use credentials already generated in
			//  db
			toKcUsers = append(toKcUsers, buDbUserInfo)
		} else {
			// Not in  BuDB or KeyCloak so needs to be added to both
			newUser := users[idx]
			newUser.CreatedTime = nowTime
			newUser.Password = generatePassword()
			toKcUsers = append(toKcUsers, newUser)
			toBuDbUsers = append(toBuDbUsers, newUser)
			// We have to replace the 1Password copy with the complete list so add it to full list
			buDbUsers[newUser.Id] = newUser
		}
	}

	if len(toBuDbUsers) != 0 {
		// Write new entries to backup db
		saveUsersToBackupDb(dbdir, toBuDbUsers)
		// Write new entries to 1Password entry
		saveFullUsersTo1Password(buDbUsers)
	}

	if len(toKcUsers) != 0 {
		// Create users on KeyCloak
		createKcUsers(ctx, client, &token, toKcUsers)
	}

	// Update System setUsersAuthCreds to True
	updateFrDbUsersWithAuthCreds(ctx, &users)
}
