package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
	ast "github.com/graphql-go/graphql/language/ast"
)

var (
	FrSchema graphql.Schema
)

////////////////////////////////////////////////////////////////////////////
// Function for retrieving selected fields
func getSelectedFields(selectionPath []string, resolveParams graphql.ResolveParams) []string {
	fields := resolveParams.Info.FieldASTs

	for _, propName := range selectionPath {
		found := false
		for _, field := range fields {
			if field.Name.Value == propName {
				selections := field.SelectionSet.Selections
				fields = make([]*ast.Field, 0)
				for _, selection := range selections {
					fields = append(fields, selection.(*ast.Field))
				}
				found = true
				break
			}
		}
		if !found {
			return []string{}
		}
	}
	var collect []string
	for _, field := range fields {
		collect = append(collect, field.Name.Value)
	}
	return collect
}

////////////////////////////////////////////////////////////////////////////
//
func init() {

	queryFields := make(map[string]*graphql.Field)
	mutationFields := make(map[string]*graphql.Field)

	//////////////////////////////////////////////////////////////////////////////
	// Mulch Timecard Common Types
	mulchTimecardType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "MulchTimecardType",
		Description: "Mulch Timecard Record Type",
		Fields: graphql.Fields{
			"id":               &graphql.Field{Type: graphql.String},
			"lastModifiedTime": &graphql.Field{Type: graphql.String},
			"deliveryId":       &graphql.Field{Type: graphql.Int},
			"timeIn":           &graphql.Field{Type: graphql.String},
			"timeOut":          &graphql.Field{Type: graphql.String},
			"timeTotal":        &graphql.Field{Type: graphql.String},
		},
	})

	//////////////////////////////////////////////////////////////////////////////
	// Order Common Types
	customerType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "CustomerType",
		Description: "Customer contact information",
		Fields: graphql.Fields{
			"addr1":        &graphql.Field{Type: graphql.String},
			"addr2":        &graphql.Field{Type: graphql.String},
			"phone":        &graphql.Field{Type: graphql.String},
			"email":        &graphql.Field{Type: graphql.String},
			"neighborhood": &graphql.Field{Type: graphql.String},
			"name":         &graphql.Field{Type: graphql.String},
		},
	})

	productType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "ProductType",
		Description: "Products Record Type",
		Fields: graphql.Fields{
			"productId":     &graphql.Field{Type: graphql.String},
			"numSold":       &graphql.Field{Type: graphql.Int},
			"amountCharged": &graphql.Field{Type: graphql.String},
		},
	})

	mulchOrderType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "MulchOrderType",
		Description: "Mulch Order Record Type",
		Fields: graphql.Fields{
			"orderId":                   &graphql.Field{Type: graphql.String},
			"ownerId":                   &graphql.Field{Type: graphql.String},
			"lastModifiedTime":          &graphql.Field{Type: graphql.String},
			"specialInstructions":       &graphql.Field{Type: graphql.String},
			"amountFromDonations":       &graphql.Field{Type: graphql.String},
			"amountFromPurchases":       &graphql.Field{Type: graphql.String},
			"amountFromCashCollected":   &graphql.Field{Type: graphql.String},
			"amountFromChecksCollected": &graphql.Field{Type: graphql.String},
			"amountTotalCollected":      &graphql.Field{Type: graphql.String},
			"checkNumbers":              &graphql.Field{Type: graphql.String},
			"willCollectMoneyLater":     &graphql.Field{Type: graphql.Boolean},
			"isVerified":                &graphql.Field{Type: graphql.Boolean},
			"customer":                  &graphql.Field{Type: customerType},
			"purchases":                 &graphql.Field{Type: graphql.NewList(productType)},
			"spreaders":                 &graphql.Field{Type: graphql.NewList(graphql.String)},
			"deliveryId":                &graphql.Field{Type: graphql.Int},
		},
	})
	archivedMulchOrderType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "ArchivedMulchOrderType",
		Description: "Archive Mulch Order Record Type",
		Fields: graphql.Fields{
			"orderId":                   &graphql.Field{Type: graphql.String},
			"ownerId":                   &graphql.Field{Type: graphql.String},
			"lastModifiedTime":          &graphql.Field{Type: graphql.String},
			"specialInstructions":       &graphql.Field{Type: graphql.String},
			"amountFromDonations":       &graphql.Field{Type: graphql.String},
			"amountFromPurchases":       &graphql.Field{Type: graphql.String},
			"amountFromCashCollected":   &graphql.Field{Type: graphql.String},
			"amountFromChecksCollected": &graphql.Field{Type: graphql.String},
			"amountTotalCollected":      &graphql.Field{Type: graphql.String},
			"checkNumbers":              &graphql.Field{Type: graphql.String},
			"willCollectMoneyLater":     &graphql.Field{Type: graphql.Boolean},
			"isVerified":                &graphql.Field{Type: graphql.Boolean},
			"customer":                  &graphql.Field{Type: customerType},
			"purchases":                 &graphql.Field{Type: graphql.NewList(productType)},
			"spreaders":                 &graphql.Field{Type: graphql.NewList(graphql.String)},
			"yearOrdered":               &graphql.Field{Type: graphql.String},
		},
	})

	customerInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "CustomerInputType",
		Description: "Customer contact input information",
		Fields: graphql.InputObjectConfigFieldMap{
			"addr1":        &graphql.InputObjectFieldConfig{Type: graphql.String},
			"addr2":        &graphql.InputObjectFieldConfig{Type: graphql.String},
			"phone":        &graphql.InputObjectFieldConfig{Type: graphql.String},
			"email":        &graphql.InputObjectFieldConfig{Type: graphql.String},
			"neighborhood": &graphql.InputObjectFieldConfig{Type: graphql.String},
			"name":         &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})

	productInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "ProductInputType",
		Description: "Products Input Record Type",
		Fields: graphql.InputObjectConfigFieldMap{
			"productId":     &graphql.InputObjectFieldConfig{Type: graphql.String},
			"numSold":       &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"amountCharged": &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})

	mulchOrderInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "MulchOrderInputType",
		Description: "Mulch Order Input Record Type",
		Fields: graphql.InputObjectConfigFieldMap{
			"orderId":                   &graphql.InputObjectFieldConfig{Type: graphql.String},
			"ownerId":                   &graphql.InputObjectFieldConfig{Type: graphql.String},
			"specialInstructions":       &graphql.InputObjectFieldConfig{Type: graphql.String},
			"amountFromDonations":       &graphql.InputObjectFieldConfig{Type: graphql.String},
			"amountFromPurchases":       &graphql.InputObjectFieldConfig{Type: graphql.String},
			"amountFromCashCollected":   &graphql.InputObjectFieldConfig{Type: graphql.String},
			"amountFromChecksCollected": &graphql.InputObjectFieldConfig{Type: graphql.String},
			"amountTotalCollected":      &graphql.InputObjectFieldConfig{Type: graphql.String},
			"checkNumbers":              &graphql.InputObjectFieldConfig{Type: graphql.String},
			"willCollectMoneyLater":     &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
			"isVerified":                &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
			"customer":                  &graphql.InputObjectFieldConfig{Type: customerInputType},
			"purchases":                 &graphql.InputObjectFieldConfig{Type: graphql.NewList(productInputType)},
			"spreaders":                 &graphql.InputObjectFieldConfig{Type: graphql.NewList(graphql.String)},
			"deliveryId":                &graphql.InputObjectFieldConfig{Type: graphql.Int},
		},
	})

	//////////////////////////////////////////////////////////////////////////////
	// Order Mutation Types
	mutationFields["createMulchOrder"] = &graphql.Field{
		Type:        graphql.String,
		Description: "Creates order",
		Args: graphql.FieldConfigArgument{
			"order": &graphql.ArgumentConfig{
				Description: "The order entry",
				Type:        mulchOrderInputType,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			log.Println("Creating Order: ", p.Args["order"])
			jsonString, err := json.Marshal(p.Args["order"])
			if err != nil {
				fmt.Println("Error encoding JSON")
				return nil, nil
			}

			newMulchOrder := MulchOrderType{}
			json.Unmarshal([]byte(jsonString), &newMulchOrder)
			return CreateMulchOrder(newMulchOrder)
		},
	}

	mutationFields["updateMulchOrder"] = &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Update order",
		Args: graphql.FieldConfigArgument{
			"order": &graphql.ArgumentConfig{
				Description: "The order entry",
				Type:        graphql.NewNonNull(mulchOrderInputType),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			jsonString, err := json.Marshal(p.Args["order"])
			if err != nil {
				fmt.Println("Error encoding JSON")
				return nil, nil
			}

			updatedMulchOrder := MulchOrderType{}
			json.Unmarshal([]byte(jsonString), &updatedMulchOrder)
			return UpdateMulchOrder(updatedMulchOrder)
		},
	}

	mutationFields["deleteMulchOrder"] = &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Deletes order associated with orderId",
		Args: graphql.FieldConfigArgument{
			"orderId": &graphql.ArgumentConfig{
				Description: "The id of the order that should be deleted",
				Type:        graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return DeleteMulchOrder(p.Args["orderId"].(string))
		},
	}

	//////////////////////////////////////////////////////////////////////////////
	// Order Query Types
	queryFields["mulchOrder"] = &graphql.Field{
		Type:        mulchOrderType,
		Description: "Retrieves order associated with orderId",
		Args: graphql.FieldConfigArgument{
			"orderId": &graphql.ArgumentConfig{
				Description: "The id of the order that should be returned",
				Type:        graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// if is_mulch_order getMulchOrder Else get etc...
			params := GetMulchOrderParams{
				OrderId:       p.Args["orderId"].(string),
				GqlFields:     getSelectedFields([]string{"mulchOrder"}, p),
				IsFromArchive: false,
			}
			return GetMulchOrder(params), nil
		},
	}

	queryFields["mulchOrders"] = &graphql.Field{
		Type:        graphql.NewList(mulchOrderType),
		Description: "Retrieves order associated with ownerId",
		Args: graphql.FieldConfigArgument{
			"ownerId": &graphql.ArgumentConfig{
				Description: "The owner id for which data should be returned.  If empty then all orders are returned",
				Type:        graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// if is_mulch_order getMulchOrder Else get etc...
			params := GetMulchOrdersParams{
				IsFromArchive: false,
				GqlFields:     getSelectedFields([]string{"mulchOrders"}, p),
			}
			if val, ok := p.Args["ownerId"]; ok {
				params.OwnerId = val.(string)
			}
			return GetMulchOrders(params), nil
		},
	}

	queryFields["archivedMulchOrder"] = &graphql.Field{
		Type:        archivedMulchOrderType,
		Description: "Retrieves order associated with orderId",
		Args: graphql.FieldConfigArgument{
			"orderId": &graphql.ArgumentConfig{
				Description: "The id of the order that should be returned",
				Type:        graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// if is_mulch_order getMulchOrder Else get etc...
			params := GetMulchOrderParams{
				OrderId:       p.Args["orderId"].(string),
				GqlFields:     getSelectedFields([]string{"archivedMulchOrder"}, p),
				IsFromArchive: true,
			}
			return GetMulchOrder(params), nil
		},
	}

	queryFields["archivedMulchOrders"] = &graphql.Field{
		Type:        graphql.NewList(archivedMulchOrderType),
		Description: "Retrieves order associated with ownerId",
		Args: graphql.FieldConfigArgument{
			"ownerId": &graphql.ArgumentConfig{
				Description: "The owner id for which data should be returned.  If empty then all orders are returned",
				Type:        graphql.String,
			},
			"archiveYear": &graphql.ArgumentConfig{
				Description: "If specified then the year (YYYY) from the archive when the order was made",
				Type:        graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// if is_mulch_order getMulchOrder Else get etc...
			params := GetMulchOrdersParams{
				IsFromArchive: true,
				GqlFields:     getSelectedFields([]string{"archivedMulchOrders"}, p),
			}
			if val, ok := p.Args["ownerId"]; ok {
				params.OwnerId = val.(string)
			}
			if val, ok := p.Args["archiveYear"]; ok {
				params.ArchiveYear = val.(string)
			}
			return GetMulchOrders(params), nil
		},
	}
	//////////////////////////////////////////////////////////////////////////////
	// Timecard Query Types
	queryFields["mulchTimeCards"] = &graphql.Field{
		Type:        graphql.NewList(mulchTimecardType),
		Description: "Retrieves Timecards for Mulch Delivery",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Description: "The id for which data should be returned.  If empty then all orders are returned",
				Type:        graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// if is_mulch_order getMulchOrder Else get etc...
			id := ""
			if val, ok := p.Args["id"]; ok {
				id = val.(string)
			}
			return GetMulchTimeCards(id), nil
		},
	}

	//////////////////////////////////////////////////////////////////////////////
	// Config Query Types
	mulchDeliveryConfigType := graphql.NewObject(graphql.ObjectConfig{
		Name: "MulchDeliveryConfigType",
		Fields: graphql.Fields{
			"id":                 &graphql.Field{Type: graphql.Int},
			"date":               &graphql.Field{Type: graphql.String},
			"newOrderCutoffDate": &graphql.Field{Type: graphql.String},
		},
	})
	neighborhoodConfigType := graphql.NewObject(graphql.ObjectConfig{
		Name: "NeighborhoodConfigType",
		Fields: graphql.Fields{
			"name":              &graphql.Field{Type: graphql.String},
			"distributionPoint": &graphql.Field{Type: graphql.String},
		},
	})
	productPriceBreakConfigType := graphql.NewObject(graphql.ObjectConfig{
		Name: "ProductPriceBreakConfigType",
		Fields: graphql.Fields{
			"gt":        &graphql.Field{Type: graphql.Int},
			"unitPrice": &graphql.Field{Type: graphql.String},
		},
	})
	productConfigType := graphql.NewObject(graphql.ObjectConfig{
		Name: "ProductConfigType",
		Fields: graphql.Fields{
			"id":          &graphql.Field{Type: graphql.String},
			"label":       &graphql.Field{Type: graphql.String},
			"minUnits":    &graphql.Field{Type: graphql.Int},
			"unitPrice":   &graphql.Field{Type: graphql.String},
			"priceBreaks": &graphql.Field{Type: graphql.NewList(productPriceBreakConfigType)},
		},
	})
	configType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "ConfigType",
		Description: "Fundraiser config inforamation",
		Fields: graphql.Fields{
			"kind":                 &graphql.Field{Type: graphql.String},
			"description":          &graphql.Field{Type: graphql.String},
			"lastModifiedTime":     &graphql.Field{Type: graphql.String},
			"isLocked":             &graphql.Field{Type: graphql.Boolean},
			"mulchDeliveryConfigs": &graphql.Field{Type: graphql.NewList(mulchDeliveryConfigType)},
			"products":             &graphql.Field{Type: graphql.NewList(productConfigType)},
			"neighborhoods":        &graphql.Field{Type: graphql.NewList(neighborhoodConfigType)},
		},
	})
	queryFields["config"] = &graphql.Field{
		Type:        configType,
		Description: "Queries for Summary information based on Owner ID",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			gqlFields := getSelectedFields([]string{"config"}, p)
			return GetFundraiserConfig(gqlFields)
		},
	}

	mulchDeliveryInputConfigType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "MulchDeliveryInputConfigType",
		Fields: graphql.InputObjectConfigFieldMap{
			"id":                 &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"date":               &graphql.InputObjectFieldConfig{Type: graphql.String},
			"newOrderCutoffDate": &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})
	neighborhoodInputConfigType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "NeighborhoodInputConfigType",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":              &graphql.InputObjectFieldConfig{Type: graphql.String},
			"distributionPoint": &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})
	productPriceBreakInputConfigType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "ProductPriceBreakInputConfigType",
		Fields: graphql.InputObjectConfigFieldMap{
			"gt":        &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"unitPrice": &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})
	productInputConfigType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "ProductInputConfigType",
		Fields: graphql.InputObjectConfigFieldMap{
			"id":          &graphql.InputObjectFieldConfig{Type: graphql.String},
			"label":       &graphql.InputObjectFieldConfig{Type: graphql.String},
			"minUnits":    &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"unitPrice":   &graphql.InputObjectFieldConfig{Type: graphql.String},
			"priceBreaks": &graphql.InputObjectFieldConfig{Type: graphql.NewList(productPriceBreakInputConfigType)},
		},
	})
	configInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "ConfigType",
		Description: "Fundraiser config inforamation",
		Fields: graphql.InputObjectConfigFieldMap{
			"kind":                 &graphql.InputObjectFieldConfig{Type: graphql.String},
			"description":          &graphql.InputObjectFieldConfig{Type: graphql.String},
			"lastModifiedTime":     &graphql.InputObjectFieldConfig{Type: graphql.String},
			"isLocked":             &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
			"mulchDeliveryConfigs": &graphql.InputObjectFieldConfig{Type: graphql.NewList(mulchDeliveryInputConfigType)},
			"products":             &graphql.InputObjectFieldConfig{Type: graphql.NewList(productInputConfigType)},
			"neighborhoods":        &graphql.InputObjectFieldConfig{Type: graphql.NewList(neighborhoodInputConfigType)},
		},
	})
	mutationFields["setConfig"] = &graphql.Field{
		Type:        graphql.Boolean,
		Description: "SetConfig",
		Args: graphql.FieldConfigArgument{
			"config": &graphql.ArgumentConfig{
				Description: "The config entry",
				Type:        configInputType,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// log.Println("Setting Config: ", p.Args["config"])
			jsonString, err := json.Marshal(p.Args["config"])
			if err != nil {
				fmt.Println("Error encoding JSON")
				return nil, nil
			}

			frConfig := FrConfigType{}
			json.Unmarshal([]byte(jsonString), &frConfig)
			return SetFundraiserConfig(frConfig)
		},
	}

	//////////////////////////////////////////////////////////////////////////////
	// Summary Query Types
	ownerIdSummaryType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "OwnerIdSummaryType",
		Description: "Summary inforamation for a specfic OnwerID",
		Fields: graphql.Fields{
			"totalDeliveryMinutes":                &graphql.Field{Type: graphql.Int},
			"totalNumBagsSold":                    &graphql.Field{Type: graphql.Int},
			"totalNumBagsSoldToSpread":            &graphql.Field{Type: graphql.Int},
			"totalAmountCollectedForDonations":    &graphql.Field{Type: graphql.String},
			"totalAmountCollectedForBags":         &graphql.Field{Type: graphql.String},
			"totalAmountCollectedForBagsToSpread": &graphql.Field{Type: graphql.String},
			"totalAmountCollected":                &graphql.Field{Type: graphql.String},
			"allocationsFromDelivery":             &graphql.Field{Type: graphql.String},
			"allocationsFromBagsSold":             &graphql.Field{Type: graphql.String},
			"allocationsFromBagsSpread":           &graphql.Field{Type: graphql.String},
			"allocationsTotal":                    &graphql.Field{Type: graphql.String},
		},
	})
	queryFields["summaryByOwnerId"] = &graphql.Field{
		Type:        ownerIdSummaryType,
		Description: "Queries for Summary information based on Owner ID",
		Args: graphql.FieldConfigArgument{
			"ownerId": &graphql.ArgumentConfig{
				Description: "The owner id for which data should be returned",
				Type:        graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return GetSummaryByOwnerId(p.Args["ownerId"].(string))
		},
	}

	troopSummaryByGroupType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "TroopSummaryByGroupType",
		Description: "Summary inforamation for the different patrols",
		Fields: graphql.Fields{
			"groupId":              &graphql.Field{Type: graphql.String},
			"totalAmountCollected": &graphql.Field{Type: graphql.String},
		},
	})

	troopSummaryTopSellersType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "TroopSummaryTopSellerType",
		Description: "List of top sellers",
		Fields: graphql.Fields{
			"name":                 &graphql.Field{Type: graphql.String},
			"totalAmountCollected": &graphql.Field{Type: graphql.String},
		},
	})

	troopSummaryType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "TroopSummaryType",
		Description: "Summary inforamation for the troop",
		Fields: graphql.Fields{
			"totalAmountCollected": &graphql.Field{Type: graphql.String},
			"groupSummary":         &graphql.Field{Type: graphql.NewList(troopSummaryByGroupType)},
			"topSellers":           &graphql.Field{Type: graphql.NewList(troopSummaryTopSellersType)},
		},
	})

	queryFields["troopSummary"] = &graphql.Field{
		Type:        troopSummaryType,
		Description: "Queries for Summary information for the entire troop",
		Args: graphql.FieldConfigArgument{
			"numTopSellers": &graphql.ArgumentConfig{
				Description: "The number of top sellers to return",
				Type:        graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return GetTroopSummary(p.Args["numTopSellers"].(int))
		},
	}

	// rootQuery := graphql.ObjectConfig{
	// 	Name:   "Query",
	// 	Fields: graphql.Fields(queryFields),
	// }

	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: graphql.Fields(queryFields)}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{Name: "Mutation", Fields: graphql.Fields(mutationFields)}),
	}

	FrSchema, _ = graphql.NewSchema(schemaConfig)
}
