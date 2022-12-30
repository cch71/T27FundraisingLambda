package frgql

import (
	"encoding/json"
	"errors"
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
	// Order Common Types
	customerType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "CustomerType",
		Description: "Customer contact information",
		Fields: graphql.Fields{
			"addr1":        &graphql.Field{Type: graphql.String},
			"addr2":        &graphql.Field{Type: graphql.String},
			"city":         &graphql.Field{Type: graphql.String},
			"zipcode":      &graphql.Field{Type: graphql.Int},
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
			"comments":                  &graphql.Field{Type: graphql.String},
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

	customerInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "CustomerInputType",
		Description: "Customer contact input information",
		Fields: graphql.InputObjectConfigFieldMap{
			"addr1":        &graphql.InputObjectFieldConfig{Type: graphql.String},
			"addr2":        &graphql.InputObjectFieldConfig{Type: graphql.String},
			"city":         &graphql.InputObjectFieldConfig{Type: graphql.String},
			"zipcode":      &graphql.InputObjectFieldConfig{Type: graphql.Int},
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
			"comments":                  &graphql.InputObjectFieldConfig{Type: graphql.String},
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
			return CreateMulchOrder(p.Context, newMulchOrder)
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
			return UpdateMulchOrder(p.Context, updatedMulchOrder)
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
			return DeleteMulchOrder(p.Context, p.Args["orderId"].(string))
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
				OrderId:   p.Args["orderId"].(string),
				GqlFields: getSelectedFields([]string{"mulchOrder"}, p),
			}
			return GetMulchOrder(params), nil
		},
	}

	queryFields["mulchOrders"] = &graphql.Field{
		Type:        graphql.NewList(mulchOrderType),
		Description: "Retrieves order associated with ownerId",
		Args: graphql.FieldConfigArgument{
			"ownerId": &graphql.ArgumentConfig{
				Description: "The owner id for which data should be returned.  If both params empty then all orders are returned",
				Type:        graphql.String,
			},
			"spreaderId": &graphql.ArgumentConfig{
				Description: "The spreader id for which data should be returned.  If both params empty then all orders are returned",
				Type:        graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// if is_mulch_order getMulchOrder Else get etc...
			params := GetMulchOrdersParams{
				GqlFields: getSelectedFields([]string{"mulchOrders"}, p),
			}
			if val, ok := p.Args["ownerId"]; ok {
				params.OwnerId = val.(string)
			}
			if val, ok := p.Args["spreaderId"]; ok {
				params.SpreaderId = val.(string)
			}
			return GetMulchOrders(params), nil
		},
	}

	//////////////////////////////////////////////////////////////////////////////
	// Mulch Timecard Common Types
	timecardType := graphql.NewObject(graphql.ObjectConfig{
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
	// Timecard Query Types
	queryFields["mulchTimecards"] = &graphql.Field{
		Type:        graphql.NewList(timecardType),
		Description: "Retrieves Timecards for Mulch Delivery",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Description: "The id for which data should be returned.  If empty then all timecards will be returned",
				Type:        graphql.String,
			},
			"deliveryId": &graphql.ArgumentConfig{
				Description: "The delivery id to return.  If empty then timecards from both deliveires will be returned",
				Type:        graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// if is_mulch_order getMulchOrder Else get etc...
			id := ""
			deliveryId := -1
			if val, ok := p.Args["id"]; ok {
				id = val.(string)
			}
			if val, ok := p.Args["deliveryId"]; ok {
				deliveryId = val.(int)
			}
			gqlFields := getSelectedFields([]string{"mulchTimecards"}, p)
			return GetMulchTimecards(id, deliveryId, gqlFields)
		},
	}

	timecardInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "MulchTimecardInputType",
		Description: "Mulch Timecard Input Entry",
		Fields: graphql.InputObjectConfigFieldMap{
			"id":         &graphql.InputObjectFieldConfig{Type: graphql.String},
			"deliveryId": &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"timeIn":     &graphql.InputObjectFieldConfig{Type: graphql.String},
			"timeOut":    &graphql.InputObjectFieldConfig{Type: graphql.String},
			"timeTotal":  &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})

	mutationFields["setMulchTimecards"] = &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Sets timecard record",
		Args: graphql.FieldConfigArgument{
			"timecards": &graphql.ArgumentConfig{
				Description: "List of timecards to record",
				Type:        graphql.NewList(timecardInputType),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// log.Println("Setting Config: ", p.Args["config"])

			jsonString, err := json.Marshal(p.Args["timecards"])
			if err != nil {
				log.Println("Error encoding JSON")
				return nil, nil
			}
			timecards := []MulchTimecardType{}
			if err := json.Unmarshal([]byte(jsonString), &timecards); err != nil {
				log.Println("Error decoding JSON to timecards")
				return nil, nil
			}
			return SetMulchTimecards(p.Context, timecards)
		},
	}

	//////////////////////////////////////////////////////////////////////////////
	// Neighborhood Query/Input Types
	neighborhoodInfoType := graphql.NewObject(graphql.ObjectConfig{
		Name: "NeighborhoodConfigType",
		Fields: graphql.Fields{
			"name":              &graphql.Field{Type: graphql.String},
			"zipcode":           &graphql.Field{Type: graphql.Int},
			"city":              &graphql.Field{Type: graphql.String},
			"isVisible":         &graphql.Field{Type: graphql.Boolean},
			"distributionPoint": &graphql.Field{Type: graphql.String},
		},
	})
	queryFields["neighborhoods"] = &graphql.Field{
		Type:        graphql.NewList(neighborhoodInfoType),
		Description: "Queries for list of neighborhoods",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			gqlFields := getSelectedFields([]string{"neighborhoods"}, p)
			return GetNeighborhoods(gqlFields)
		},
	}
	neighborhoodInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "NeighborhoodInfoInputType",
		Description: "Fundraiser Neighborhood Input",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":              &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"zipcode":           &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"city":              &graphql.InputObjectFieldConfig{Type: graphql.String},
			"isVisible":         &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
			"distributionPoint": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		},
	})

	mutationFields["addNeighborhoods"] = &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Add neighborhood(s) to system",
		Args: graphql.FieldConfigArgument{
			"neighborhoods": &graphql.ArgumentConfig{
				Description: "List of neighborhoods",
				Type:        graphql.NewList(neighborhoodInputType),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// log.Println("Setting Config: ", p.Args["config"])
			jsonString, err := json.Marshal(p.Args["neighborhoods"])
			if err != nil {
				log.Println("Error encoding JSON")
				return nil, nil
			}
			hoods := []NeighborhoodInfo{}
			if err := json.Unmarshal([]byte(jsonString), &hoods); err != nil {
				log.Println("Error decoding JSON to NeighborhoodInfo")
				return nil, nil
			}
			return AddNeighborhoods(p.Context, hoods)
		},
	}

	mutationFields["updateNeighborhoods"] = &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Add neighborhood(s) to system",
		Args: graphql.FieldConfigArgument{
			"neighborhoods": &graphql.ArgumentConfig{
				Description: "List of neighborhoods",
				Type:        graphql.NewList(neighborhoodInputType),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// log.Println("Setting Config: ", p.Args["config"])
			jsonString, err := json.Marshal(p.Args["neighborhoods"])
			if err != nil {
				log.Println("Error encoding JSON")
				return nil, nil
			}
			hoods := []NeighborhoodInfo{}
			if err := json.Unmarshal([]byte(jsonString), &hoods); err != nil {
				log.Println("Error decoding JSON to NeighborhoodInfo")
				return nil, nil
			}
			return UpdateNeighborhoods(p.Context, hoods)
		},
	}

	//////////////////////////////////////////////////////////////////////////////
	// User/Group Query/Input Types
	userInfoType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "UserInfoType",
		Description: "User Info Type",
		Fields: graphql.Fields{
			"firstName":    &graphql.Field{Type: graphql.String},
			"lastName":     &graphql.Field{Type: graphql.String},
			"name":         &graphql.Field{Type: graphql.String},
			"id":           &graphql.Field{Type: graphql.String},
			"group":        &graphql.Field{Type: graphql.String},
			"hasAuthCreds": &graphql.Field{Type: graphql.Boolean},
		},
	})

	queryFields["users"] = &graphql.Field{
		Type:        graphql.NewList(userInfoType),
		Description: "Queries for list of users",
		Args: graphql.FieldConfigArgument{
			"showOnlyUsersWithoutAuthCreds": &graphql.ArgumentConfig{
				Description: "If true will filter list to only show users with hasAuthCreds==false",
				Type:        graphql.Boolean,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			params := GetUsersParams{
				GqlFields:                 getSelectedFields([]string{"users"}, p),
				ShowUsersWithoutAuthCreds: false,
			}
			if val, ok := p.Args["showOnlyUsersWithoutAuthCreds"]; ok {
				params.ShowUsersWithoutAuthCreds = val.(bool)
			}
			return GetUsers(params)
		},
	}

	userInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "UserInfoInputType",
		Description: "Fundraiser user",
		Fields: graphql.InputObjectConfigFieldMap{
			"firstName":    &graphql.InputObjectFieldConfig{Type: graphql.String},
			"lastName":     &graphql.InputObjectFieldConfig{Type: graphql.String},
			"id":           &graphql.InputObjectFieldConfig{Type: graphql.String},
			"group":        &graphql.InputObjectFieldConfig{Type: graphql.String},
			"hasAuthCreds": &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
		},
	})

	mutationFields["addOrUpdateUsers"] = &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Add user(s) to system",
		Args: graphql.FieldConfigArgument{
			"users": &graphql.ArgumentConfig{
				Description: "List of users",
				Type:        graphql.NewList(userInputType),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// log.Println("Setting Config: ", p.Args["config"])
			jsonString, err := json.Marshal(p.Args["users"])
			if err != nil {
				log.Println("Error encoding JSON")
				return nil, nil
			}
			users := []UserInfo{}
			if err := json.Unmarshal([]byte(jsonString), &users); err != nil {
				log.Println("Error decoding JSON to userinfo")
				return nil, nil
			}
			return AddOrUpdateUsers(p.Context, users)
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
	finalizationDataConfigType := graphql.NewObject(graphql.ObjectConfig{
		Name: "finalizationDataConfigType",
		Fields: graphql.Fields{
			"bankDeposited":              &graphql.Field{Type: graphql.String},
			"mulchCost":                  &graphql.Field{Type: graphql.String},
			"perBagCost":                 &graphql.Field{Type: graphql.String},
			"profitsFromBags":            &graphql.Field{Type: graphql.String},
			"mulchSalesGross":            &graphql.Field{Type: graphql.String},
			"moneyPoolForTroop":          &graphql.Field{Type: graphql.String},
			"moneyPoolForScoutsSubPools": &graphql.Field{Type: graphql.String},
			"moneyPoolForScoutsSales":    &graphql.Field{Type: graphql.String},
			"moneyPoolForScoutsDelivery": &graphql.Field{Type: graphql.String},
			"perBagAvgEarnings":          &graphql.Field{Type: graphql.String},
			"deliveryEarningsPerMinute":  &graphql.Field{Type: graphql.String},
		},
	})
	configType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "ConfigType",
		Description: "Fundraiser config information",
		Fields: graphql.Fields{
			"kind":                 &graphql.Field{Type: graphql.String},
			"description":          &graphql.Field{Type: graphql.String},
			"lastModifiedTime":     &graphql.Field{Type: graphql.String},
			"isLocked":             &graphql.Field{Type: graphql.Boolean},
			"mulchDeliveryConfigs": &graphql.Field{Type: graphql.NewList(mulchDeliveryConfigType)},
			"products":             &graphql.Field{Type: graphql.NewList(productConfigType)},
			"finalizationData":     &graphql.Field{Type: finalizationDataConfigType},
			"neighborhoods":        queryFields["neighborhoods"],
			"users":                queryFields["users"],
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
	finalizationDataInputConfigType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "finalizationDataInputConfigType",
		Fields: graphql.InputObjectConfigFieldMap{
			"bankDeposited":              &graphql.InputObjectFieldConfig{Type: graphql.String},
			"mulchCost":                  &graphql.InputObjectFieldConfig{Type: graphql.String},
			"perBagCost":                 &graphql.InputObjectFieldConfig{Type: graphql.String},
			"profitsFromBags":            &graphql.InputObjectFieldConfig{Type: graphql.String},
			"mulchSalesGross":            &graphql.InputObjectFieldConfig{Type: graphql.String},
			"moneyPoolForTroop":          &graphql.InputObjectFieldConfig{Type: graphql.String},
			"moneyPoolForScoutsSubPools": &graphql.InputObjectFieldConfig{Type: graphql.String},
			"moneyPoolForScoutsSales":    &graphql.InputObjectFieldConfig{Type: graphql.String},
			"moneyPoolForScoutsDelivery": &graphql.InputObjectFieldConfig{Type: graphql.String},
			"perBagAvgEarnings":          &graphql.InputObjectFieldConfig{Type: graphql.String},
			"deliveryEarningsPerMinute":  &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})
	configInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "ConfigType",
		Description: "Fundraiser config information",
		Fields: graphql.InputObjectConfigFieldMap{
			"kind":                 &graphql.InputObjectFieldConfig{Type: graphql.String},
			"description":          &graphql.InputObjectFieldConfig{Type: graphql.String},
			"lastModifiedTime":     &graphql.InputObjectFieldConfig{Type: graphql.String},
			"isLocked":             &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
			"mulchDeliveryConfigs": &graphql.InputObjectFieldConfig{Type: graphql.NewList(mulchDeliveryInputConfigType)},
			"products":             &graphql.InputObjectFieldConfig{Type: graphql.NewList(productInputConfigType)},
			"finalizationData":     &graphql.InputObjectFieldConfig{Type: finalizationDataInputConfigType},
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
			return SetFundraiserConfig(p.Context, frConfig)
		},
	}
	mutationFields["updateConfig"] = &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Update existing config values",
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
			return UpdateFundraiserConfig(p.Context, frConfig)
		},
	}

	//////////////////////////////////////////////////////////////////////////////
	// Adds Spreaders to order
	mutationFields["setSpreaders"] = &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Sets spreader information for an order",
		Args: graphql.FieldConfigArgument{
			"orderId": &graphql.ArgumentConfig{
				Description: "The id of the order associated with the spreaders",
				Type:        graphql.NewNonNull(graphql.String),
			},
			"spreaders": &graphql.ArgumentConfig{
				Description: "list of userids that performed the spreading, can be empty",
				Type:        graphql.NewNonNull(graphql.NewList(graphql.String)),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// log.Println("Setting Config: ", p.Args["config"])
			orderId := p.Args["orderId"].(string)
			jsonString, err := json.Marshal(p.Args["spreaders"])
			if err != nil {
				return false, errors.New("spreaders param not formatted correctly")
			}
			spreaders := []string{}
			if err := json.Unmarshal([]byte(jsonString), &spreaders); err != nil {
				return false, errors.New("spreaders could not be decoded")
			}
			return SetSpreaders(orderId, spreaders)
		},
	}

	//////////////////////////////////////////////////////////////////////////////
	// Creates Issue Report
	newIssueInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "NewIssueType",
		Description: "Issue Report Information",
		Fields: graphql.InputObjectConfigFieldMap{
			"id":    &graphql.InputObjectFieldConfig{Type: graphql.String},
			"title": &graphql.InputObjectFieldConfig{Type: graphql.String},
			"body":  &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})
	mutationFields["createIssue"] = &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Creates a new issue report",
		Args: graphql.FieldConfigArgument{
			"input": &graphql.ArgumentConfig{
				Description: "Issue Report Information",
				Type:        newIssueInputType,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// log.Println("Setting Config: ", p.Args["config"])
			jsonString, err := json.Marshal(p.Args["input"])
			if err != nil {
				fmt.Println("Error encoding JSON")
				return nil, nil
			}

			issue := NewIssue{}
			json.Unmarshal([]byte(jsonString), &issue)
			return CreateIssue(issue)
		},
	}

	//////////////////////////////////////////////////////////////////////////////
	// Set Fundraiser closeout allocations
	allocationsInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "AllocationsInputType",
		Description: "Allocations",
		Fields: graphql.InputObjectConfigFieldMap{
			"uid":                       &graphql.InputObjectFieldConfig{Type: graphql.String},
			"bagsSold":                  &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"bagsSpread":                &graphql.InputObjectFieldConfig{Type: graphql.String},
			"deliveryMinutes":           &graphql.InputObjectFieldConfig{Type: graphql.String},
			"totalDonations":            &graphql.InputObjectFieldConfig{Type: graphql.String},
			"allocationsFromBagsSold":   &graphql.InputObjectFieldConfig{Type: graphql.String},
			"allocationsFromBagsSpread": &graphql.InputObjectFieldConfig{Type: graphql.String},
			"allocationsFromDelivery":   &graphql.InputObjectFieldConfig{Type: graphql.String},
			"allocationsTotal":          &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})
	mutationFields["setFundraiserCloseoutAllocations"] = &graphql.Field{
		Type:        graphql.Boolean,
		Description: "SetFundraiserCloseoutAllocations",
		Args: graphql.FieldConfigArgument{
			"allocations": &graphql.ArgumentConfig{
				Description: "Allocations",
				Type:        graphql.NewList(allocationsInputType),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// log.Println("Setting Config: ", p.Args["config"])
			jsonString, err := json.Marshal(p.Args["allocations"])
			if err != nil {
				fmt.Println("Error encoding JSON")
				return nil, nil
			}

			allocations := []AllocationItemType{}
			json.Unmarshal([]byte(jsonString), &allocations)
			return SetFrCloseoutAllocations(p.Context, allocations)
		},
	}
	//////////////////////////////////////////////////////////////////////////////
	// Summary Query Types
	ownerIdSummaryType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "OwnerIdSummaryType",
		Description: "Summary information for a specfic OnwerID",
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

	orderOwnerSummary := graphql.Field{
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
		Description: "Summary information for the different patrols",
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
		Description: "Summary information for the troop",
		Fields: graphql.Fields{
			"totalAmountCollected": &graphql.Field{Type: graphql.String},
			"groupSummary":         &graphql.Field{Type: graphql.NewList(troopSummaryByGroupType)},
			"topSellers":           &graphql.Field{Type: graphql.NewList(troopSummaryTopSellersType)},
		},
	})

	troopSummary := graphql.Field{
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

	neighborhoodSummaryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "NeighborhoodSummaryType",
		Fields: graphql.Fields{
			"neighborhood": &graphql.Field{Type: graphql.String},
			"numOrders":    &graphql.Field{Type: graphql.Int},
		},
	})

	neighborhoodsSummary := graphql.Field{
		Type:        graphql.NewList(neighborhoodSummaryType),
		Description: "Summary information for neighborhoods",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return GetNeighborhoodSummary()
		},
	}

	summaryType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "SummaryType",
		Description: "Summary information",
		Fields: graphql.Fields{
			"neighborhoods": &neighborhoodsSummary,
			"orderOwner":    &orderOwnerSummary,
			"troop":         &troopSummary,
		},
	})

	queryFields["summary"] = &graphql.Field{
		Type: summaryType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// graphql-go requires this shim to do sublevel queries.  Without this
			// the sub resolves wouldn't trigger
			type Shimmer struct {
				Troop         TroopSummaryType
				OrderOwner    OwnerIdSummaryType
				Neighborhoods []NeighborhoodSummaryType
			}
			return Shimmer{}, nil
		},
	}

	//////////////////////////////////////////////////////////////////////////////
	// Resets Fundraiser User and Order Data
	mutationFields["resetFundraisingData"] = &graphql.Field{
		Type:        graphql.Boolean,
		Description: "ResetFundraiser",
		Args: graphql.FieldConfigArgument{
			"doResetUsers": &graphql.ArgumentConfig{
				Description: "Resets users list",
				Type:        graphql.Boolean,
			},
			"doResetOrders": &graphql.ArgumentConfig{
				Description: "Resets mulch orders, mulch spreaders, allocation summary, mulch delivery config, finalization data, and time cards data",
				Type:        graphql.Boolean,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			doResetUsers, doResetOrders := false, false
			if val, ok := p.Args["doResetUsers"]; ok {
				doResetUsers = val.(bool)
			}
			if val, ok := p.Args["doResetOrders"]; ok {
				doResetOrders = val.(bool)
			}
			return ResetFundraisingData(p.Context, doResetUsers, doResetOrders)
		},
	}

	//////////////////////////////////////////////////////////////////////////////
	// Geolocation Address Type
	addressType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "AddressType",
		Description: "Address Information",
		Fields: graphql.Fields{
			"zipcode":     &graphql.Field{Type: graphql.Int},
			"city":        &graphql.Field{Type: graphql.String},
			"houseNumber": &graphql.Field{Type: graphql.String},
			"street":      &graphql.Field{Type: graphql.String},
		},
	})
	queryFields["getAddress"] = &graphql.Field{
		Type:        addressType,
		Description: "Retrieves Address from geo information",
		Args: graphql.FieldConfigArgument{
			"lat": &graphql.ArgumentConfig{
				Description: "Latitude",
				Type:        graphql.NewNonNull(graphql.Float),
			},
			"lng": &graphql.ArgumentConfig{
				Description: "Longitude",
				Type:        graphql.NewNonNull(graphql.Float),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var lat, lng float64
			if val, ok := p.Args["lat"]; ok {
				lat = val.(float64)
			}
			if val, ok := p.Args["lng"]; ok {
				lng = val.(float64)
			}
			return GetAddrFromLatLng(p.Context, lat, lng)
		},
	}

	// queryFields["testApi"] = &graphql.Field{
	// 	Type:        graphql.Boolean,
	// 	Description: "",
	// 	Args: graphql.FieldConfigArgument{
	// 		"param1": &graphql.ArgumentConfig{
	// 			Description: "",
	// 			Type:        graphql.NewNonNull(graphql.String),
	// 		},
	// 	},
	// 	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
	// 		return AdminTestApi(p.Context, p.Args["param1"].(string))
	// 	},
	// }

	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: graphql.Fields(queryFields)}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{Name: "Mutation", Fields: graphql.Fields(mutationFields)}),
	}

	FrSchema, _ = graphql.NewSchema(schemaConfig)
}
