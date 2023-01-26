package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/cch71/T27FundraisingLambda/frgql"
)

////////////////////////////////////////////////////////////////////////////
//
type LambdaRequestBody struct {
	Query string `json:"query"`
}

////////////////////////////////////////////////////////////////////////////
//
type LambdaRequest struct {
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers,omitempty"`
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
func HandleLambdaEvent(ctx context.Context, event LambdaRequest) (LambdaResponse, error) {
	//if dbconn not already established then RwLock to go ahead and try once sync.Once
	//check authorization
	//run query
	//return results
	if err := frgql.OpenDb(); err != nil {
		log.Println("Failed to initialize db:", err)
		return generateResp("", http.StatusInternalServerError), err
	}

	if bearerToken, prs := event.Headers["Authorization"]; prs {
		ctx = context.WithValue(ctx, "T27FrAuthorization", bearerToken[len("Bearer "):])
		// We don't want this printing out in the log
		delete(event.Headers, "Authorization")
	}

	log.Println("Rxed GraphQL Query: ", event)

	body := LambdaRequestBody{}
	json.Unmarshal([]byte(event.Body), &body)

	respBody, err := frgql.MakeGqlQuery(ctx, body.Query)
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
	frgql.CloseDb()
}
