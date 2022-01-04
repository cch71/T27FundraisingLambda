package main

import (
	"log"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/google/uuid"
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
    isLocked: false,
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
          id: 1,
          date: "3/13/2022",
          newOrderCutoffDate: "2/19/2022"
      },
      {
          id: 2,
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

var createMulchOrderGql1 = `
mutation {
  createMulchOrder(order: {
    orderId: "**UUID**"
    ownerId: "Blogger1"
    customer: {
        name: "John Ford"
        addr1: "123 Hola Gato Dr"
        addr2: "Apt 4c"
        phone: "211-234-5434"
        email: "my@noreply.com"
        neighborhood: "Brown Bear"
    }
    willCollectMoneyLater: true
    deliveryId: 2
    purchases: [{
        productId: "bags"
        numSold: 24
        amountCharged: "200.00"
    }]
    amountFromPurchases: "200.00"
  })
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

var updateMulchOrderGql1 = `
mutation {
  updateMulchOrder(order: {
    orderId: "**UUID**"
    ownerId: "Blogger1"
    customer: {
        name: "John Ford"
        addr1: "123 Hola Gato Dr"
        phone: "211-234-5434"
        neighborhood: "Brown Bear"
    }
    specialInstructions: "Don't leave it where I can see it"
    willCollectMoneyLater: false
    amountFromDonations: "25.00"
    amountFromCashCollected: "20.00"
    amountFromChecksCollected: "5.00"
    checkNumbers: "1234 1235"
    amountTotalCollected: "25.00"
  })
}
`
var deleteMulchOrderGql1 = `
mutation {
  deleteMulchOrder(orderId: "**UUID**")
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

var IntrospectionQuery1 = `
  query IntrospectionQuery {
    __schema {
      queryType { name }
      mutationType { name }
      subscriptionType { name }
      types {
        ...FullType
      }
      directives {
        name
        description
		locations
        args {
          ...InputValue
        }
        # deprecated, but included for coverage till removed
		onOperation
        onFragment
        onField
      }
    }
  }
  fragment FullType on __Type {
    kind
    name
    description
    fields(includeDeprecated: true) {
      name
      description
      args {
        ...InputValue
      }
      type {
        ...TypeRef
      }
      isDeprecated
      deprecationReason
    }
    inputFields {
      ...InputValue
    }
    interfaces {
      ...TypeRef
    }
    enumValues(includeDeprecated: true) {
      name
      description
      isDeprecated
      deprecationReason
    }
    possibleTypes {
      ...TypeRef
    }
  }
  fragment InputValue on __InputValue {
    name
    description
    type { ...TypeRef }
    defaultValue
  }
  fragment TypeRef on __Type {
    kind
    name
    ofType {
      kind
      name
      ofType {
        kind
        name
        ofType {
          kind
          name
          ofType {
            kind
            name
            ofType {
              kind
              name
              ofType {
                kind
                name
                ofType {
                  kind
                  name
                }
              }
            }
          }
        }
      }
    }
  }
`
var InstrospectionQuery2 = `
query allSchemaTypes {
 __schema {
    queryType {
      fields {
        name
        description
      }
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

func TestGraphQLSummaryCreateUpdateAndDeleteOrder(t *testing.T) {
	uuidStr := uuid.New().String()

	{
		gql := strings.ReplaceAll(createMulchOrderGql1, "**UUID**", uuidStr)
		rJSON, err := MakeGqlQuery(gql)
		if err != nil {
			t.Fatal("GraphQL Create Order Failed: ", err)
		}
		t.Logf("%s \n\n", rJSON)
	}
	{
		gql := strings.ReplaceAll(queryOrdersGql1, "**UUID**", uuidStr)
		rJSON, err := MakeGqlQuery(gql)
		if err != nil {
			t.Fatal("GraphQL Query1 Failed: ", err)
		}
		t.Logf("\n%s \n\n", rJSON)
	}
	{
		gql := strings.ReplaceAll(updateMulchOrderGql1, "**UUID**", uuidStr)
		rJSON, err := MakeGqlQuery(gql)
		if err != nil {
			t.Fatal("GraphQL Update Order Failed: ", err)
		}
		t.Logf("\n%s \n\n", rJSON)
	}
	{
		gql := strings.ReplaceAll(queryOrdersGql1, "**UUID**", uuidStr)
		rJSON, err := MakeGqlQuery(gql)
		if err != nil {
			t.Fatal("GraphQL Query2 Failed: ", err)
		}
		t.Logf("\n%s \n\n", rJSON)
	}
	{
		gql := strings.ReplaceAll(deleteMulchOrderGql1, "**UUID**", uuidStr)
		rJSON, err := MakeGqlQuery(gql)
		if err != nil {
			t.Fatal("GraphQL Delete Order Failed: ", err)
		}
		t.Logf("\n%s \n\n", rJSON)
	}
}

func TestGraphQLIntrospectionQuery(t *testing.T) {
	{
		rJSON, err := MakeGqlQuery(InstrospectionQuery2)
		if err != nil {
			t.Fatal("GraphQL Query Failed: ", err)
		}
		t.Logf("%s \n", rJSON) // {"data":{"hello":"world"}}
	}
}

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
