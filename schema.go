package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

var (
	FrSchema graphql.Schema
)

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
			"firstName":    &graphql.Field{Type: graphql.String},
			"lastName":     &graphql.Field{Type: graphql.String},
		},
	})

	mulchProductType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "MulchProductType",
		Description: "Mulch Products Record Type",
		Fields: graphql.Fields{
			"bagsSold":                  &graphql.Field{Type: graphql.Int},
			"bagsToSpread":              &graphql.Field{Type: graphql.Int},
			"AmountChargedForBags":      &graphql.Field{Type: graphql.String},
			"AmountChargedForSpreading": &graphql.Field{Type: graphql.String},
		},
	})

	mulchOrderType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "MulchOrderType",
		Description: "Mulch Order Record Type",
		Fields: graphql.Fields{
			"orderId":                      &graphql.Field{Type: graphql.String},
			"ownerId":                      &graphql.Field{Type: graphql.String},
			"lastModifiedTime":             &graphql.Field{Type: graphql.String},
			"specialInstructions":          &graphql.Field{Type: graphql.String},
			"amountFromDonationsCollected": &graphql.Field{Type: graphql.String},
			"amountFromCashCollected":      &graphql.Field{Type: graphql.String},
			"amountFromChecksCollected":    &graphql.Field{Type: graphql.String},
			"amountTotalCollected":         &graphql.Field{Type: graphql.String},
			"checkNumbers":                 &graphql.Field{Type: graphql.NewList(graphql.String)},
			"willCollectMoneyLater":        &graphql.Field{Type: graphql.Boolean},
			"isVerified":                   &graphql.Field{Type: graphql.Boolean},
			"customer":                     &graphql.Field{Type: customerType},
			"purchases":                    &graphql.Field{Type: mulchProductType},
			"spreaders":                    &graphql.Field{Type: graphql.NewList(graphql.String)},
			"deliveryId":                   &graphql.Field{Type: graphql.String},
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
			"firstName":    &graphql.InputObjectFieldConfig{Type: graphql.String},
			"lastName":     &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})

	mulchProductInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "MulchProductInputType",
		Description: "Mulch Products Input Record Type",
		Fields: graphql.InputObjectConfigFieldMap{
			"bagsSold":                  &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"bagsToSpread":              &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"AmountChargedForBags":      &graphql.InputObjectFieldConfig{Type: graphql.String},
			"AmountChargedForSpreading": &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})

	mulchOrderInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "MulchOrderInputType",
		Description: "Mulch Order Input Record Type",
		Fields: graphql.InputObjectConfigFieldMap{
			"orderId":                      &graphql.InputObjectFieldConfig{Type: graphql.String},
			"ownerId":                      &graphql.InputObjectFieldConfig{Type: graphql.String},
			"lastModifiedTime":             &graphql.InputObjectFieldConfig{Type: graphql.String},
			"specialInstructions":          &graphql.InputObjectFieldConfig{Type: graphql.String},
			"amountFromDonationsCollected": &graphql.InputObjectFieldConfig{Type: graphql.String},
			"amountFromCashCollected":      &graphql.InputObjectFieldConfig{Type: graphql.String},
			"amountFromChecksCollected":    &graphql.InputObjectFieldConfig{Type: graphql.String},
			"amountTotalCollected":         &graphql.InputObjectFieldConfig{Type: graphql.String},
			"checkNumbers":                 &graphql.InputObjectFieldConfig{Type: graphql.NewList(graphql.String)},
			"willCollectMoneyLater":        &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
			"isVerified":                   &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
			"customer":                     &graphql.InputObjectFieldConfig{Type: customerInputType},
			"purchases":                    &graphql.InputObjectFieldConfig{Type: mulchProductInputType},
			"spreaders":                    &graphql.InputObjectFieldConfig{Type: graphql.NewList(graphql.String)},
			"deliveryId":                   &graphql.InputObjectFieldConfig{Type: graphql.String},
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
			return CreateMulchOrder(newMulchOrder), nil
		},
	}

	mutationFields["updateMulchOrder"] = &graphql.Field{
		Type:        graphql.String,
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
			return UpdateMulchOrder(updatedMulchOrder), nil
		},
	}

	mutationFields["deleteMulchOrder"] = &graphql.Field{
		Type:        graphql.String,
		Description: "Deletes order associated with orderId",
		Args: graphql.FieldConfigArgument{
			"orderId": &graphql.ArgumentConfig{
				Description: "The id of the order that should be deleted",
				Type:        graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return DeleteMulchOrder(p.Args["orderId"].(string)), nil
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
			return GetMulchOrder(p.Args["orderId"].(string)), nil
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
			id := ""
			if val, ok := p.Args["ownerId"]; ok {
				id = val.(string)
			}
			return GetMulchOrders(id), nil
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
			"id":                 &graphql.Field{Type: graphql.String},
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
			"mulchDeliveryConfigs": &graphql.Field{Type: graphql.NewList(mulchDeliveryConfigType)},
			"products":             &graphql.Field{Type: graphql.NewList(productConfigType)},
			"isLocked":             &graphql.Field{Type: graphql.Boolean},
			"neighborhoods": &graphql.Field{
				Type: graphql.NewList(neighborhoodConfigType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					type NeighborhoodItem struct {
						Name              string
						DistributionPoint string
					}
					hoodList := []NeighborhoodItem{}
					theHoods := p.Source.(FrConfigType).Neighborhoods
					for name, obj := range theHoods {
						hoodList = append(hoodList,
							NeighborhoodItem{
								Name:              name,
								DistributionPoint: obj.DistributionPoint,
							})
					}
					return hoodList, nil
				},
			},
		},
	})
	queryFields["config"] = &graphql.Field{
		Type:        configType,
		Description: "Queries for Summary information based on Owner ID",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return GetFundraisingConfig()
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
			return GetOwnerIdSummary(p.Args["ownerId"].(string)), nil
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
			"topSellers": &graphql.Field{
				Type: graphql.NewList(troopSummaryTopSellersType),
				Args: graphql.FieldConfigArgument{
					"numTopSellers": &graphql.ArgumentConfig{
						Description: "The number of top sellers to return",
						Type:        graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// log.Println("ResolveParams: topSellers ", p.Args["numTopSellers"].(int), " Src: ", p.Source)
					return p.Source.(TroopSummaryType).TopSellers, nil
				},
			},
		},
	})

	queryFields["troopSummary"] = &graphql.Field{
		Type:        troopSummaryType,
		Description: "Queries for Summary information for the entire troop",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			troopSummary := GetTroopSummary(11)
			// log.Println("Troop Summary: ", troopSummary)
			return troopSummary, nil
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
