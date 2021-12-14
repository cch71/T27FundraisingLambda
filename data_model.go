package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	//"github.com/vishalkuo/bimap"
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

	// GraphQL<->SQL Table Mapping
	// 	mulchOrderFields = bimap.NewBiMap()
	// 	mulchOrderFields.Insert("orderId", "order_id")
	// 	mulchOrderFields.Insert("ownerId", "owner_id")
	// 	mulchOrderFields.Insert("lastModifiedTime", "last_modified_time")
	// 	mulchOrderFields.Insert("specialInstructions", "special_instructions")
	// 	mulchOrderFields.Insert("amountFromDonationsCollected", "donation_amount_collected")
	// 	mulchOrderFields.Insert("amountFromCashCollected", "cash_amount_collected")
	// 	mulchOrderFields.Insert("amountFromChecksCollected", "check_amount_collected")
	// 	mulchOrderFields.Insert("amountTotalCollected", "total_amount_collected")
	// 	mulchOrderFields.Insert("checkNumbers", "check_numbers")
	// 	mulchOrderFields.Insert("willCollectMoneyLater", "will_collect_money_later")
	// 	mulchOrderFields.Insert("isVerified", "is_verified")
	// 	mulchOrderFields.Insert("deliveryId", "delivery_id")
	// 	mulchOrderFields.Insert("yearOrdered", "year_ordered")

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
	Neighborhood string
	FirstName    string
	LastName     string
}

type MulchProductsType struct {
	BagsSold                  int
	BagsToSpread              int
	AmountChargedForBags      string
	AmountChargedForSpreading string
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
	ArchiveYear                  string
}

type GetMulchOrdersParams struct {
	OwnerId       string
	GqlFields     []string
	IsFromArchive bool
	ArchiveYear   string
}

func mulchOrderGql2SqlMap(gqlFields []string) []string {
	sqlFields := []string{}
	for gqlField := range gqlFields {
		log.Println(gqlField)
		// mulchOrderFields.Get(gqlField)

	}
	return sqlFields
}

func GetMulchOrders(params GetMulchOrdersParams) []MulchOrderType {
	log.Println("Retrieving OwnerId: ", params.OwnerId)
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
				FirstName:    "James",
				LastName:     "Something",
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
	log.Println("Retrieving OrderID: ", params.OrderId)
	return MulchOrderType{
		OrderId:                 "24",
		OwnerId:                 "BobbyJo",
		LastModifiedTime:        "Now+",
		AmountFromCashCollected: "22.11",
		AmountTotalCollected:    "22.11",
		Customer: CustomerType{
			Addr1:        "192 Subway",
			Neighborhood: "Sesame",
			Phone:        "444.444.4442",
			FirstName:    "James",
			LastName:     "Something",
		},
		Purchases: MulchProductsType{
			BagsSold:             24,
			AmountChargedForBags: "22.11",
		},
	}
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
