package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/graphql-go/graphql"
)

// var invokeCount = 0
// var myObjects []*s3.Object
// func init() {
//         svc := s3.New(session.New())
//         input := &s3.ListObjectsV2Input{
//                 Bucket: aws.String("examplebucket"),
//         }
//         result, _ := svc.ListObjectsV2(input)
//         myObjects = result.Contents
// }

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

////////////////////////////////////////////////////////////////////////////
//
type LambdaRequestBody struct {
	Query string `json:"query"`
}

////////////////////////////////////////////////////////////////////////////
//
type LambdaRequest struct {
	Body string `json:"body"`
}

////////////////////////////////////////////////////////////////////////////
//
type LambdaResponse struct {
	Body       string            `json:"body"`
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
}

////////////////////////////////////////////////////////////////////////////
//
func generateResp(body string, statusCode int) LambdaResponse {
	return LambdaResponse{
		StatusCode: statusCode,
		Body:       body,
		Headers: map[string]string{
			// "Access-Control-Allow-Methods": "POST,OPTIONS",
			// "Access-Control-Allow-Headers": "X-Amz-Date,X-Api-Key,X-Amz-Security-Token,X-Requested-With,X-Auth-Token,Referer,User-Agent,Origin,Content-Type,Authorization,Accept,Access-Control-Allow-Methods,Access-Control-Allow-Origin,Access-Control-Allow-Headers",
			"Access-Control-Allow-Origin": "*",
		},
	}

}

////////////////////////////////////////////////////////////////////////////
//
func generateOkResp(body string) LambdaResponse {
	return generateResp(body, http.StatusOK)
}

////////////////////////////////////////////////////////////////////////////
//
func HandleLambdaEvent(_ctx context.Context, event LambdaRequest) (LambdaResponse, error) {
	//if dbconn not already established then RwLock to go ahead and try once sync.Once
	//check authorization
	//run query
	//return results
	if err := InitDb(); err != nil {
		log.Println("Failed to initialize db:", err)
		return generateResp("", http.StatusInternalServerError), err
	}
	log.Println("Rxed GraphQL Query: ", event)

	body := LambdaRequestBody{}
	json.Unmarshal([]byte(event.Body), &body)

	respBody, err := MakeGqlQuery(body.Query)
	if err != nil {
		log.Println("GraphQL Query Failed: ", err)
		return generateResp("", http.StatusBadRequest), err
	}

	return generateOkResp(string(respBody)), nil
}

////////////////////////////////////////////////////////////////////////////
//
func main() {
	lambda.Start(HandleLambdaEvent)
	if Db != nil {
		Db.Close()
	}
}
