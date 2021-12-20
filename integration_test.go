package main

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/joho/godotenv"
)

var summaryQueryGql1 = `
{
  config {
    kind
    lastModifiedTime
    isLocked
  }
  summaryByOwnerId(ownerId: "Bob") {
    totalDeliveryMinutes
    totalNumBagsSold
    totalNumBagsSoldToSpread
    totalAmountCollectedForDonations
    totalAmountCollectedForBags
    totalAmountCollectedForBagsToSpread
    totalAmountCollected
    allocationsFromDelivery
    allocationsFromBagsSold
    allocationsFromBagsSpread
    allocationsTotal
  }
  troopSummary(numTopSellers: 10) {
    totalAmountCollected
    topSellers {
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

var activeConfigQueryGql1 = `
{
  config {
    description
    kind
    isLocked
    lastModifiedTime
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

var setActiveConfigMutationGql = `
mutation {
  setConfig(config: {
    kind: "mulch",
    description: "Mulch",
    isLocked: true,
    neighborhoods: [
      {
         name: "Avery Ranch",
         distributionPoint: "Walsh"
      },{

         name: "Behrens Ranch",
         distributionPoint: "Walsh"
      }
    ],
    mulchDeliveryConfigs: [
      {
          id: "1",
          date: "3/13/2022",
          newOrderCutoffDate: "2/19/2022"
      },
      {
          id: "2",
          date: "4/10/2022",
          newOrderCutoffDate: "3/27/2022"
      }
    ],
    products: [
      {
          id: "bags",
          label: "Bags of Mulch",
          unitPrice: "4.15",
          minUnits: 5,
          priceBreaks: [
              {
                  gt: 15,
                  unitPrice: "4.00"
              },{
                  gt: 35,
                  unitPrice: "3.85"
              },{
                  gt: 64,
                  unitPrice: "3.75"
              }
          ]
      },{
          id: "spreading",
          label: "Bags to Spread",
          minUnits: 5,
          unitPrice: "2.00"
      }
    ]
  })
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
  archivedMulchOrders(ownerId: "parvg") {
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

func TestGraphQLArchivedSingleOrder(t *testing.T) {
	{
		rJSON, err := MakeGqlQuery(archivedOrderQueryGql1)
		if err != nil {
			t.Fatal("GraphQL Query Failed: ", err)
		}
		t.Logf("%s \n", rJSON) // {"data":{"hello":"world"}}
	}
}

func TestGraphQLArchivedParvgOrders(t *testing.T) {
	{
		rJSON, err := MakeGqlQuery(archivedOrdersQueryGql1)
		if err != nil {
			t.Fatal("GraphQL Query Failed: ", err)
		}
		t.Logf("%s \n", rJSON) // {"data":{"hello":"world"}}
	}
}

func TestGraphQLSetActiveConfig(t *testing.T) {
	{
		rJSON, err := MakeGqlQuery(setActiveConfigMutationGql)
		if err != nil {
			t.Fatal("GraphQL Mutation Failed: ", err)
		}
		t.Logf("%s \n", rJSON)
	}
}

func TestGraphQLActiveConfig(t *testing.T) {
	{
		rJSON, err := MakeGqlQuery(activeConfigQueryGql1)
		if err != nil {
			t.Fatal("GraphQL Query Failed: ", err)
		}
		t.Logf("%s \n", rJSON) // {"data":{"hello":"world"}}
	}
}

func TestGraphQLSummaryQuery(t *testing.T) {
	{
		rJSON, err := MakeGqlQuery(summaryQueryGql1)
		if err != nil {
			t.Fatal("GraphQL Query Failed: ", err)
		}
		t.Logf("%s \n", rJSON) // {"data":{"hello":"world"}}
	}
}
