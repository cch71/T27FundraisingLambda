package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	dbMutex sync.Mutex
	Db      *pgxpool.Pool
	//mulchOrderFields bimap.BiMap
)

func InitDb() error {
	if Db == nil {
		dbMutex.Lock()
		defer dbMutex.Unlock()
		if Db == nil {
			cnxn, err := makeDbConnection()
			if err != nil {
				return err
			}
			Db = cnxn
		}
	}

	return nil
}

func makeDbConnection() (*pgxpool.Pool, error) {

	dbId := os.Getenv("DB_ID")
	dbToken := os.Getenv("DB_TOKEN")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbCaRoot := os.Getenv("DB_CA_ROOT_PATH")

	dbName := "defaultdb"
	cluster := "pushy-iguana-1562"

	dbOptions := url.PathEscape(fmt.Sprintf("--cluster=%s", cluster))
	dbParams := fmt.Sprintf("%s?sslmode=verify-full&sslrootcert=%s&options=%s", dbName, dbCaRoot, dbOptions)
	cnxnUri := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbId, dbToken, dbHost, dbPort, dbParams)
	// Attempt to connect
	// log.Println("\n\nCnxn String: ", cnxnUri, "\n")
	conn, err := pgxpool.Connect(context.Background(), cnxnUri)
	if err != nil {
		return nil, err
	}
	// defer conn.Close()
	return conn, nil
}

type OwnerIdSummaryType struct {
	TotalDeliveryMinutes                int
	TotalNumBagsSold                    int
	TotalNumBagsSoldToSpread            int
	TotalAmountCollectedForDonations    string
	TotalAmountCollectedForBags         string
	TotalAmountCollectedForBagsToSpread string
	TotalAmountCollected                string
	AllocationsFromDelivery             string
	AllocationsFromBagsSold             string
	AllocationsFromBagsSpread           string
	AllocationsTotal                    string
}

func GetOwnerIdSummary(ownerId string) OwnerIdSummaryType {
	log.Println("Getting Summary for onwerId: ", ownerId)
	return OwnerIdSummaryType{
		TotalDeliveryMinutes:             408,
		TotalNumBagsSold:                 24,
		TotalAmountCollectedForDonations: "52.44",
	}
}

type TopSellerType struct {
	Name                 string
	TotalAmountCollected string
}

type GroupSummaryType struct {
	GroupId              string
	TotalAmountCollected string
}

type TroopSummaryType struct {
	TotalAmountCollected string
	GroupSummary         []GroupSummaryType
	TopSellers           []TopSellerType
}

func GetTroopSummary(numSellers int) TroopSummaryType {
	log.Println("Getting this many top sellers: ", numSellers)
	return TroopSummaryType{
		TotalAmountCollected: "66.75",
		GroupSummary:         []GroupSummaryType{GroupSummaryType{GroupId: "bears", TotalAmountCollected: "22.34"}, GroupSummaryType{GroupId: "lions", TotalAmountCollected: "42.34"}},
		TopSellers:           []TopSellerType{TopSellerType{Name: "John", TotalAmountCollected: "11.23"}},
	}
}

type CustomerType struct {
	Addr1        string
	Addr2        string
	Phone        string
	Email        string
	Neighborhood string
	Name         string
}

type MulchProductsType struct {
	// Bags                      int // legacy
	//Spreading                 int // legacy
	BagsSold                  int    `json:"bags,omitempty" json:"bagsSold"`
	BagsToSpread              int    `json:"spreading,omitempty" json:"bagsToSpread"`
	AmountChargedForBags      string `json:"amountChargedForBags,omitempty"`
	AmountChargedForSpreading string `json:"amountChargedForSpreading,omitempty"`
}

type MulchOrderType struct {
	OrderId                      string
	OwnerId                      string
	LastModifiedTime             string
	SpecialInstructions          string
	AmountFromDonationsCollected string
	AmountFromCashCollected      string
	AmountFromChecksCollected    string
	AmountTotalCollected         string
	CheckNumbers                 []string
	WillCollectMoneyLater        bool
	IsVerified                   bool
	Customer                     CustomerType
	Purchases                    MulchProductsType
	YearOrdered                  string
}

type GetMulchOrdersParams struct {
	OwnerId       string
	GqlFields     []string
	IsFromArchive bool
	ArchiveYear   string
}

func mulchOrderGql2SqlMap(gqlFields []string, orderOutput *MulchOrderType) ([]string, []interface{}) {

	sqlFields := []string{}
	inputs := []interface{}{}
	for _, gqlField := range gqlFields {
		// log.Println(gqlField)
		switch {
		case gqlField == "orderId":
			inputs = append(inputs, &orderOutput.OrderId)
			sqlFields = append(sqlFields, "order_id")
		case gqlField == "ownerId":
			inputs = append(inputs, &orderOutput.OwnerId)
			sqlFields = append(sqlFields, "order_owner_id")
		case gqlField == "amountTotalCollected":
			inputs = append(inputs, &orderOutput.AmountTotalCollected)
			sqlFields = append(sqlFields, "total_amount_collected::string")
		case gqlField == "yearOrdered":
			inputs = append(inputs, &orderOutput.YearOrdered)
			sqlFields = append(sqlFields, "year_ordered::string")
		case gqlField == "customer":
			inputs = append(inputs, &orderOutput.Customer.Name)
			sqlFields = append(sqlFields, "customer_name")
			inputs = append(inputs, &orderOutput.Customer.Addr1)
			sqlFields = append(sqlFields, "customer_addr1")
			inputs = append(inputs, &orderOutput.Customer.Addr2)
			sqlFields = append(sqlFields, "customer_addr2")
			inputs = append(inputs, &orderOutput.Customer.Phone)
			sqlFields = append(sqlFields, "customer_phone")
			inputs = append(inputs, &orderOutput.Customer.Email)
			sqlFields = append(sqlFields, "customer_email")
			inputs = append(inputs, &orderOutput.Customer.Neighborhood)
			sqlFields = append(sqlFields, "customer_neighborhood")
		case gqlField == "purchases":
			inputs = append(inputs, &orderOutput.Purchases)
			sqlFields = append(sqlFields, "purchases::jsonb")
		}
		// GraphQL<->SQL Table Mapping
		// 	mulchOrderFields = bimap.NewBiMap()
		// 	mulchOrderFields.Insert("lastModifiedTime", "last_modified_time")
		// 	mulchOrderFields.Insert("specialInstructions", "special_instructions")
		// 	mulchOrderFields.Insert("amountFromDonationsCollected", "donation_amount_collected")
		// 	mulchOrderFields.Insert("amountFromCashCollected", "cash_amount_collected")
		// 	mulchOrderFields.Insert("amountFromChecksCollected", "check_amount_collected")
		// 	mulchOrderFields.Insert("checkNumbers", "check_numbers")
		// 	mulchOrderFields.Insert("willCollectMoneyLater", "will_collect_money_later")
		// 	mulchOrderFields.Insert("isVerified", "is_verified")
		// 	mulchOrderFields.Insert("deliveryId", "delivery_id")
		// mulchOrderFields.Get(gqlField)

	}
	return sqlFields, inputs
}

func GetMulchOrders(params GetMulchOrdersParams) []MulchOrderType {
	if 0 == len(params.OwnerId) {
		log.Println("Retrieving mulch orders. ", "Is targeting archive: ", params.IsFromArchive)
		log.Println("Selections: ", params.GqlFields)

	} else {
		log.Println("Retrieving mulch orders. ", "Is targeting archive: ", params.IsFromArchive, " OwnerId: ", params.OwnerId)

	}

	// order := MulchOrderType{}
	// _sqlFields, _inputs := mulchOrderGql2SqlMap(params.GqlFields, &order)
	return []MulchOrderType{
		MulchOrderType{
			OrderId:                 "24",
			OwnerId:                 "BobbyJo",
			LastModifiedTime:        "Now+",
			AmountFromCashCollected: "22.11",
			AmountTotalCollected:    "22.11",
			Customer: CustomerType{
				Addr1:        "192 Subway",
				Neighborhood: "Sesame",
				Phone:        "444.444.4442",
				Name:         "James",
			},
			Purchases: MulchProductsType{
				BagsSold:             24,
				AmountChargedForBags: "22.11",
			},
		},
	}
}

type GetMulchOrderParams struct {
	OrderId       string
	GqlFields     []string
	IsFromArchive bool
}

func GetMulchOrder(params GetMulchOrderParams) MulchOrderType {
	log.Println("Retrieving mulch order. ", "Is targeting archive: ", params.IsFromArchive, " OrderId: ", params.OrderId)

	order := MulchOrderType{}
	sqlFields, inputs := mulchOrderGql2SqlMap(params.GqlFields, &order)

	dbTable := "mulch_orders"
	if params.IsFromArchive {
		dbTable = "archived_mulch_orders"
	}
	sqlCmd := fmt.Sprintf("select %s from %s where order_id=$1", strings.Join(sqlFields, ","), dbTable)
	log.Println("SqlCmd: ", sqlCmd)
	err := Db.QueryRow(context.Background(), sqlCmd, params.OrderId).Scan(inputs...)
	if err != nil {
		log.Println("Mulch order query for: ", params.OrderId, " failed", err)
	}
	// log.Println("Purchases: ", order.Purchases)
	return order
}

func CreateMulchOrder(order MulchOrderType) string {
	log.Println("Creating Order: ", order)
	return order.OrderId
}

func UpdateMulchOrder(order MulchOrderType) string {
	log.Println("Updating Order: ", order)
	return order.OrderId
}

func DeleteMulchOrder(orderId string) string {
	log.Println("Deleteing OrderID: ", orderId)
	return orderId
}

type MulchDeliveryConfigType struct {
	Id                 string `json:"id"`
	Date               string `json:"date"`
	NewOrderCutoffDate string `json:"newOrderCutoffDate"`
}

type NeighborhoodsType struct {
	DistributionPoint string `json:"distPt"`
}

type ProductPriceBreaks struct {
	Gt        int    `json:"gt"`
	UnitPrice string `json:"unitPrice"`
}

type ProductType struct {
	Id          string               `json:"id"`
	Label       string               `json:"label"`
	MinUnits    int                  `json:"minUnits"`
	UnitPrice   string               `json:"unitPrice"`
	PriceBreaks []ProductPriceBreaks `json:"priceBreaks"`
}

type FrConfigType struct {
	Kind                 string                       `json:"kind"`
	Description          string                       `json:"description"`
	IsLocked             bool                         `json:"isLocked"`
	Neighborhoods        map[string]NeighborhoodsType `json:"neighborhoods"`
	MulchDeliveryConfigs []MulchDeliveryConfigType    `json:"mulchDeliveryConfigs"`
	Products             []ProductType                `json:"products"`
}

func GetFundraisingConfig() (FrConfigType, error) {
	jsonFile, err := os.Open("/Users/chamilton/wip/t27/t27utils/v2/T27FundraiserConfig.json")
	if err != nil {
		return FrConfigType{}, err
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	frConfig := FrConfigType{}
	json.Unmarshal(byteValue, &frConfig)
	return frConfig, nil

}

type MulchTimecardType struct {
	Id               string
	DeliveryId       int
	LastModifiedTime time.Time
	TimeIn           string
	TimeOut          string
	TimeTotal        string
}

func GetMulchTimeCards(id string) []MulchTimecardType {
	doQuery := func(id string) (pgx.Rows, error) {
		if len(id) == 0 {
			log.Println("Retrieving All Timecards")
			sqlCmd := `select uid, delivery_id, last_modify_time, time_in::string, time_out::string, time_total::string from mulch_delivery_timecards`
			return Db.Query(context.Background(), sqlCmd)
		} else {
			log.Println("Retrieving Timecards for: ", id)
			sqlCmd := `select uid, delivery_id, last_modify_time, time_in::string, time_out::string, time_total::string from mulch_delivery_timecards where uid=$1`
			return Db.Query(context.Background(), sqlCmd, id)
		}
	}

	timecards := []MulchTimecardType{}
	rows, err := doQuery(id)
	if err != nil {
		log.Println("Timecard query Failed", err)
		return timecards
	}
	defer rows.Close()

	for rows.Next() {
		tc := MulchTimecardType{}
		inputs := []interface{}{&tc.Id, &tc.DeliveryId, &tc.LastModifiedTime, &tc.TimeIn, &tc.TimeOut, &tc.TimeTotal}
		err = rows.Scan(inputs...)
		if err != nil {
			log.Println("Reading timecard row failed: ", err)
			continue
		}
		timecards = append(timecards, tc)
	}

	if rows.Err() != nil {
		log.Println("Reading timecard rows had an issue: ", err)
		return []MulchTimecardType{}
	}
	return timecards
}
