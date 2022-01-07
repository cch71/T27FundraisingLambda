package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/cch71/T27FundraisingLambda/frgql"
	"github.com/joho/godotenv"
)

////////////////////////////////////////////////////////////////////////////
//
func main() {

	credentialsFile := path.Join(os.Getenv("HOME"), ".t27fr", "credentials")
	_ = godotenv.Load(credentialsFile)

	gqlFn := flag.String("gql", "", "GraphGQ File")
	flag.Parse()

	query, err := ioutil.ReadFile(*gqlFn)
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
