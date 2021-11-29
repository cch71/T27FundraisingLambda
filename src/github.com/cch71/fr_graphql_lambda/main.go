package main

import (
	"net/http"

	"github.com/cch71/fr_graphql_lambda/frsvc"

	"github.com/graphql-go/handler"
)

func main() {
	h := handler.New(&handler.Config{
		Schema:     &frsvc.FrSchema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})

	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", nil)
}
