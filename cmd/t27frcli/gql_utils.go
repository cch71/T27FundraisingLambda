package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/cch71/T27FundraisingLambda/frgql"
)

// //////////////////////////////////////////////////////////////////////////
func MakeGqlReq(ctx context.Context, gqlFn *string) {

	// Open File
	query, err := os.ReadFile(*gqlFn)
	if err != nil {
		log.Panic("Failed opening file: ", *gqlFn, " Err: ", err)
	}

	// Initialize Database Connection and Keycloak token
	if err := frgql.OpenDb(); err != nil {
		log.Panic("Failed to initialize db:", err)
	}
	defer frgql.CloseDb()

	_, token := LoginKcAdmin(ctx)
	ctx = context.WithValue(ctx, "T27FrAuthorization", token)

	// Make GQL Query
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
