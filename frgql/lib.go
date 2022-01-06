package frgql

import (
	"encoding/json"
	"log"

	"github.com/graphql-go/graphql"
)

////////////////////////////////////////////////////////////////////////////
//
func MakeGqlQuery(gql string) ([]byte, error) {
	params := graphql.Params{Schema: FrSchema, RequestString: gql}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Printf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, err := json.Marshal(r)
	if err != nil {
		log.Println("Error encoding JSON results: ", err, " for gql: ", gql)
		return nil, err
	}
	return rJSON, nil
}
