package frgql

import (
	"context"
	"encoding/json"
	"log"

	"github.com/graphql-go/graphql"
)

////////////////////////////////////////////////////////////////////////////
//
func MakeGqlQuery(ctx context.Context, gql string) ([]byte, error) {
	params := graphql.Params{Schema: FrSchema, RequestString: gql, Context: ctx}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Printf("failed to execute graphql operation:\n%s\n, errors: %+v", gql, r.Errors)
		return nil, r.Errors[0]
	}

	rJSON, err := json.Marshal(r)
	if err != nil {
		log.Println("Error encoding JSON results: ", err, " for gql: ", gql)
		return nil, err
	}
	return rJSON, nil
}
