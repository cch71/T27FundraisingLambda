package frgql

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/joho/godotenv"
)

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

var queryOrdersGql1 = `
{
  mulchOrders(ownerId: "Blogger1") {
    orderId
    ownerId
    amountFromPurchases
    amountFromDonations
    amountTotalCollected
    willCollectMoneyLater
    isVerified
    customer {
        name
        addr1
        addr2
        phone
        email
        neighborhood
    }
    specialInstructions
    purchases {
        productId
        numSold
        amountCharged
    }
  }
}
`

// var createMulchOrderGql1 = `
// mutation {
//   createMulchOrder(order: {
//     orderId: "**UUID**"
//     ownerId: "Blogger1"
//     customer: {
//         name: "John Ford"
//         addr1: "123 Hola Gato Dr"
//         addr2: "Apt 4c"
//         phone: "211-234-5434"
//         email: "my@noreply.com"
//         neighborhood: "Brown Bear"
//     }
//     willCollectMoneyLater: true
//     deliveryId: 2
//     purchases: [{
//         productId: "bags"
//         numSold: 24
//         amountCharged: "200.00"
//     }]
//     amountFromPurchases: "200.00"
//   })
// }
// `
// var updateMulchOrderGql1 = `
// mutation {
//   updateMulchOrder(order: {
//     orderId: "**UUID**"
//     ownerId: "Blogger1"
//     customer: {
//         name: "John Ford"
//         addr1: "123 Hola Gato Dr"
//         phone: "211-234-5434"
//         neighborhood: "Brown Bear"
//     }
//     specialInstructions: "Don't leave it where I can see it"
//     willCollectMoneyLater: false
//     amountFromDonations: "25.00"
//     amountFromCashCollected: "20.00"
//     amountFromChecksCollected: "5.00"
//     checkNumbers: "1234 1235"
//     amountTotalCollected: "25.00"
//   })
// }
// `
// var deleteMulchOrderGql1 = `
// mutation {
//   deleteMulchOrder(orderId: "**UUID**")
// }
// `

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
        productId
        numSold
        amountCharged
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
        productId
        numSold
        amountCharged
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
	credentialsFile := path.Join(os.Getenv("HOME"), ".cockroachdb", "credentials")
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

// func TestGraphQLArchivedParvgOrders(t *testing.T) {
// 	{
// 		rJSON, err := MakeGqlQuery(archivedOrdersQueryGql1)
// 		if err != nil {
// 			t.Fatal("GraphQL Query Failed: ", err)
// 		}
// 		t.Logf("%s \n", rJSON) // {"data":{"hello":"world"}}
// 	}
// }

// func TestGraphQLSetActiveConfig(t *testing.T) {
// 	{
// 		rJSON, err := MakeGqlQuery(setActiveConfigMutationGql)
// 		if err != nil {
// 			t.Fatal("GraphQL Mutation Failed: ", err)
// 		}
// 		t.Logf("%s \n", rJSON)
// 	}
// }

func TestGraphQLActiveConfig(t *testing.T) {
	{
		rJSON, err := MakeGqlQuery(activeConfigQueryGql1)
		if err != nil {
			t.Fatal("GraphQL Query Failed: ", err)
		}
		t.Logf("%s \n", rJSON) // {"data":{"hello":"world"}}
	}
}

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
var summaryQueryGql2 = `
{
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
var summaryQueryGql3 = `
{
  summaryByOwnerId(ownerId: "fruser2") {
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
}
`

func TestGraphQLSummaryQuery(t *testing.T) {
	{
		rJSON, err := MakeGqlQuery(summaryQueryGql3)
		if err != nil {
			t.Fatal("GraphQL Query Failed: ", err)
		}
		t.Logf("%s \n", rJSON) // {"data":{"hello":"world"}}
	}
}

// func TestGraphQLCreateUpdateAndDeleteOrder(t *testing.T) {
// 	uuidStr := uuid.New().String()
// 	{
// 		gql := strings.ReplaceAll(createMulchOrderGql1, "**UUID**", uuidStr)
// 		rJSON, err := MakeGqlQuery(gql)
// 		if err != nil {
// 			t.Fatal("GraphQL Create Order Failed: ", err)
// 		}
// 		t.Logf("%s \n\n", rJSON)
// 	}
// 	{
// 		gql := strings.ReplaceAll(queryOrdersGql1, "**UUID**", uuidStr)
// 		rJSON, err := MakeGqlQuery(gql)
// 		if err != nil {
// 			t.Fatal("GraphQL Query1 Failed: ", err)
// 		}
// 		t.Logf("\n%s \n\n", rJSON)
// 	}
// 	{
// 		gql := strings.ReplaceAll(updateMulchOrderGql1, "**UUID**", uuidStr)
// 		rJSON, err := MakeGqlQuery(gql)
// 		if err != nil {
// 			t.Fatal("GraphQL Update Order Failed: ", err)
// 		}
// 		t.Logf("\n%s \n\n", rJSON)
// 	}
// 	{
// 		gql := strings.ReplaceAll(queryOrdersGql1, "**UUID**", uuidStr)
// 		rJSON, err := MakeGqlQuery(gql)
// 		if err != nil {
// 			t.Fatal("GraphQL Query2 Failed: ", err)
// 		}
// 		t.Logf("\n%s \n\n", rJSON)
// 	}
// 	{
// 		gql := strings.ReplaceAll(deleteMulchOrderGql1, "**UUID**", uuidStr)
// 		rJSON, err := MakeGqlQuery(gql)
// 		if err != nil {
// 			t.Fatal("GraphQL Delete Order Failed: ", err)
// 		}
// 		t.Logf("\n%s \n\n", rJSON)
// 	}
// }

var queryOrdersQuickReportGql1 = `
{
  mulchOrders(ownerId: "fradmin") {
    orderId
    deliveryId
    spreaders
    customer {
        name
    }
    purchases {
        productId
        numSold
        amountCharged
    }
  }
}
`

func TestGraphQLQueryQuickReportFrAdmin(t *testing.T) {
	rJSON, err := MakeGqlQuery(queryOrdersQuickReportGql1)
	if err != nil {
		t.Fatal("GraphQL Query Failed: ", err)
	}
	t.Logf("%s \n", rJSON) // {"data":{"hello":"world"}}
}

var queryUsersGql1 = `
{
  users {
    id
    group
    name
  }
}
`

func TestGraphQLQueryUsers(t *testing.T) {
	rJSON, err := MakeGqlQuery(queryUsersGql1)
	if err != nil {
		t.Fatal("GraphQL Query Failed: ", err)
	}
	t.Logf("%s \n", rJSON) // {"data":{"hello":"world"}}
}
