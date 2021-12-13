package main

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/joho/godotenv"
)

var summaryConfigQueryGql = `
{
  summaryByOwnerId(ownerId: "Bob") {
  	totalDeliveryMinutes
  }
  troopSummary {
    totalAmountCollected
    topSellers(numTopSellers: 10) {
      totalAmountCollected
      name
    }
    groupSummary {
      groupId
      totalAmountCollected
    }
  }
}
`

var timecardQueryGql1 = `
{
    mulchTimeCards(id: "axell") {
        id
        lastModifiedTime
        timeIn
        timeOut
        timeTotal
    }
}
`

var timecardQueryGql2 = `
{
    mulchTimeCards {
        id
        lastModifiedTime
        timeIn
        timeOut
        timeTotal
    }
}
`

var createMulchOrderGql = `
mutation {
  createMulchOrder(order: {
    orderId: "24"
    ownerId: "Blogger1"
  })
}
`

// func TestGraphQL(t *testing.T) {
// 	h := handler.New(&handler.Config{
// 		Schema:     &FrSchema,
// 		Pretty:     true,
// 		GraphiQL:   false,
// 		Playground: true,
// 	})
//
// 	http.Handle("/graphql", h)
// 	http.ListenAndServe(":8080", nil)
// }

func TestMain(m *testing.M) {
	// Write code here to run before tests
	credentialsFile := path.Join(os.Getenv("HOME"), ".cockroachdb", "credentials.yml")
	_ = godotenv.Load(credentialsFile)
	if err := InitDb(); err != nil {
		log.Fatal("Failed to initialize db:", err)
	}
	defer Db.Close()

	// Run tests
	exitVal := m.Run()

	// Write code here to run after tests

	// Exit with exit value from tests
	os.Exit(exitVal)
}

func TestGraphQLTimeCards(t *testing.T) {
	{
		rJSON, err := MakeGqlQuery(timecardQueryGql1)
		if err != nil {
			t.Fatal("GraphQL Query Failed: ", err)
		}
		t.Logf("%s \n", rJSON) // {"data":{"hello":"world"}}
	}
	// {
	// 	rJSON, err := MakeGqlQuery(timecardQueryGql2)
	// 	if err != nil {
	// 		t.Fatal("GraphQL Query Failed: ", err)
	// 	}
	// 	t.Logf("%s \n", rJSON) // {"data":{"hello":"world"}}
	// }
}
