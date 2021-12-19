package main

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/joho/godotenv"
)

var summaryQueryGql = `
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

var frConfigQueryGql = `
{
  config {
    description
    kind
    isLocked
    neighborhoods {
      name
      distributionPoint
    }
    mulchDeliveryConfigs {
      id
      date
      newOrderCutoffDate
    }
    products {
      id
      label
      unitPrice
      minUnits
      priceBreaks {
        gt
        unitPrice
      }
    }
  }
}
`

var timecardQueryGql1 = `
{
    mulchTimeCards(id: "axell") {
        id
        deliveryId
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
        deliveryId
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

var archivedOrderQueryGql1 = `
{
  archivedMulchOrder(orderId: "04a4c1b4-6a6d-4e20-a09d-1dc28d7d5955") {
    ownerId
    amountTotalCollected
    yearOrdered
    customer {
        name
        addr1
        addr2
        phone
        email
        neighborhood
    }
    purchases {
        bagsSold
        bagsToSpread
        amountChargedForBags
        amountChargedForSpreading
    }
  }
}
`

var archivedOrdersQueryGql1 = `
{
  archivedMulchOrders {
    orderId
    ownerId
    amountTotalCollected
    yearOrdered
    customer {
        name
        addr1
        addr2
        phone
        email
        neighborhood
    }
    purchases {
        bagsSold
        bagsToSpread
        amountChargedForBags
        amountChargedForSpreading
    }
  }
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

func TestGraphQLArchivedOrder(t *testing.T) {
	{
		rJSON, err := MakeGqlQuery(archivedOrderQueryGql1)
		if err != nil {
			t.Fatal("GraphQL Query Failed: ", err)
		}
		t.Logf("%s \n", rJSON) // {"data":{"hello":"world"}}
	}
}

func TestGraphQLArchivedOrders(t *testing.T) {
	{
		rJSON, err := MakeGqlQuery(archivedOrdersQueryGql1)
		if err != nil {
			t.Fatal("GraphQL Query Failed: ", err)
		}
		t.Logf("%s \n", rJSON) // {"data":{"hello":"world"}}
	}
}
