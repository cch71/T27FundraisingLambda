package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
var (
	dbMutex sync.Mutex
	Db      *pgxpool.Pool
	//mulchOrderFields bimap.BiMap
)

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////
//
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

////////////////////////////////////////////////////////////////////////////
//
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

////////////////////////////////////////////////////////////////////////////
//
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

////////////////////////////////////////////////////////////////////////////
//
func GetSummaryByOwnerId(ownerId string) (OwnerIdSummaryType, error) {
	log.Println("Getting Summary for onwerId: ", ownerId)
	return OwnerIdSummaryType{
		TotalDeliveryMinutes:             408,
		TotalNumBagsSold:                 24,
		TotalAmountCollectedForDonations: "52.44",
	}, nil
}

////////////////////////////////////////////////////////////////////////////
//
type TopSellerType struct {
	Name                 string
	TotalAmountCollected string
}

////////////////////////////////////////////////////////////////////////////
//
type GroupSummaryType struct {
	GroupId              string
	TotalAmountCollected string
}

////////////////////////////////////////////////////////////////////////////
//
type TroopSummaryType struct {
	TotalAmountCollected string
	GroupSummary         []GroupSummaryType
	TopSellers           []TopSellerType
}

////////////////////////////////////////////////////////////////////////////
//
func GetTroopSummary(numTopSellers int) (TroopSummaryType, error) {
	log.Println("Getting this many top sellers: ", numTopSellers)
	return TroopSummaryType{
		TotalAmountCollected: "66.75",
		GroupSummary:         []GroupSummaryType{GroupSummaryType{GroupId: "bears", TotalAmountCollected: "22.34"}, GroupSummaryType{GroupId: "lions", TotalAmountCollected: "42.34"}},
		TopSellers:           []TopSellerType{TopSellerType{Name: "John", TotalAmountCollected: "11.23"}},
	}, nil
}

////////////////////////////////////////////////////////////////////////////
//
type CustomerType struct {
	Addr1        string
	Addr2        string
	Phone        string
	Email        string
	Neighborhood string
	Name         string
}

////////////////////////////////////////////////////////////////////////////
//
type MulchProductsType struct {
	// Bags                      int // legacy
	//Spreading                 int // legacy
	BagsSold                  int    `json:"bags,omitempty" json:"bagsSold"`
	BagsToSpread              int    `json:"spreading,omitempty" json:"bagsToSpread"`
	AmountChargedForBags      string `json:"amountChargedForBags,omitempty"`
	AmountChargedForSpreading string `json:"amountChargedForSpreading,omitempty"`
}

////////////////////////////////////////////////////////////////////////////
//
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
	DeliveryId                   int    // Not in archived GraphQL
	YearOrdered                  string // Not in non archived GraphQL
}

////////////////////////////////////////////////////////////////////////////
//
type GetMulchOrdersParams struct {
	OwnerId       string
	GqlFields     []string
	IsFromArchive bool
	ArchiveYear   string
}

////////////////////////////////////////////////////////////////////////////
//
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
		case gqlField == "last_modified_time":
			inputs = append(inputs, &orderOutput.LastModifiedTime)
			sqlFields = append(sqlFields, "last_modified_time")
		case gqlField == "specialInstructions":
			inputs = append(inputs, &orderOutput.SpecialInstructions)
			sqlFields = append(sqlFields, "special_instructions")
		case gqlField == "amountFromDonationsCollected":
			inputs = append(inputs, &orderOutput.AmountFromDonationsCollected)
			sqlFields = append(sqlFields, "donation_amount_collected::string")
		case gqlField == "amountFromCashCollected":
			inputs = append(inputs, &orderOutput.AmountFromCashCollected)
			sqlFields = append(sqlFields, "cash_amount_collected::string")
		case gqlField == "amountFromChecksCollected":
			inputs = append(inputs, &orderOutput.AmountFromChecksCollected)
			sqlFields = append(sqlFields, "check_amount_collected::string")
		case gqlField == "checkNumbers":
			inputs = append(inputs, &orderOutput.WillCollectMoneyLater)
			sqlFields = append(sqlFields, "check_numbers::jsonb")
		case gqlField == "deliveryId":
			inputs = append(inputs, &orderOutput.DeliveryId)
			sqlFields = append(sqlFields, "delivery_id")
		case gqlField == "willCollectMoneyLater":
			inputs = append(inputs, &orderOutput.Purchases)
			sqlFields = append(sqlFields, "will_collect_money_later")
		case gqlField == "isVerified":
			inputs = append(inputs, &orderOutput.IsVerified)
			sqlFields = append(sqlFields, "is_verified")
		default:
			log.Println("Do not know how to handle GraphQL Field: ", gqlField)
		}

	}
	return sqlFields, inputs
}

////////////////////////////////////////////////////////////////////////////
//
func GetMulchOrders(params GetMulchOrdersParams) []MulchOrderType {

	order := MulchOrderType{}
	sqlFields, _ := mulchOrderGql2SqlMap(params.GqlFields, &order)

	dbTable := "mulch_orders"
	if params.IsFromArchive {
		dbTable = "archived_mulch_orders"
	}

	if 0 == len(params.OwnerId) {
		log.Println("Retrieving mulch orders. ", "Is targeting archive: ", params.IsFromArchive)

	} else {
		log.Println("Retrieving mulch orders. ", "Is targeting archive: ", params.IsFromArchive, " OwnerId: ", params.OwnerId)

	}

	doQuery := func(id *string, dbTable *string, sqlFields []string) (pgx.Rows, error) {
		sqlCmd := fmt.Sprintf("select %s from %s", strings.Join(sqlFields, ","), *dbTable)
		if len(*id) == 0 {
			log.Println("SqlCmd: ", sqlCmd)
			return Db.Query(context.Background(), sqlCmd)
		} else {
			sqlCmd = sqlCmd + " where order_owner_id=$1"
			log.Println("SqlCmd: ", sqlCmd)
			return Db.Query(context.Background(), sqlCmd, *id)
		}
	}

	orders := []MulchOrderType{}
	rows, err := doQuery(&params.OwnerId, &dbTable, sqlFields)
	if err != nil {
		log.Println("Mulch Orders query failed", err)
		return orders
	}
	defer rows.Close()

	for rows.Next() {
		order := MulchOrderType{}
		_, inputs := mulchOrderGql2SqlMap(params.GqlFields, &order)
		err = rows.Scan(inputs...)
		if err != nil {
			log.Println("Reading mulch order row failed: ", err)
			continue
		}
		orders = append(orders, order)
	}

	if rows.Err() != nil {
		log.Println("Reading mulch order rows had an issue: ", err)
		return []MulchOrderType{}
	}
	return orders
}

////////////////////////////////////////////////////////////////////////////
//
type GetMulchOrderParams struct {
	OrderId       string
	GqlFields     []string
	IsFromArchive bool
}

////////////////////////////////////////////////////////////////////////////
//
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

////////////////////////////////////////////////////////////////////////////
//
func CreateMulchOrder(order MulchOrderType) string {
	log.Println("Creating Order: ", order)
	return order.OrderId
}

////////////////////////////////////////////////////////////////////////////
//
func UpdateMulchOrder(order MulchOrderType) string {
	log.Println("Updating Order: ", order)
	return order.OrderId
}

////////////////////////////////////////////////////////////////////////////
//
func DeleteMulchOrder(orderId string) string {
	log.Println("Deleteing OrderID: ", orderId)
	return orderId
}

////////////////////////////////////////////////////////////////////////////
//
type MulchDeliveryConfigType struct {
	Id                 string `json:"id"`
	Date               string `json:"date"`
	NewOrderCutoffDate string `json:"newOrderCutoffDate"`
}

////////////////////////////////////////////////////////////////////////////
//
type NeighborhoodsType struct {
	Name              string `json: name`
	DistributionPoint string `json:"distributionPt"`
}

////////////////////////////////////////////////////////////////////////////
//
type ProductPriceBreaks struct {
	Gt        int    `json:"gt"`
	UnitPrice string `json:"unitPrice"`
}

////////////////////////////////////////////////////////////////////////////
//
type ProductType struct {
	Id          string               `json:"id"`
	Label       string               `json:"label"`
	MinUnits    int                  `json:"minUnits"`
	UnitPrice   string               `json:"unitPrice"`
	PriceBreaks []ProductPriceBreaks `json:"priceBreaks"`
}

////////////////////////////////////////////////////////////////////////////
//
type FrConfigType struct {
	Kind                 string                    `json:"kind"`
	Description          string                    `json:"description"`
	LastModifiedTime     string                    `json:"lastModifiedTime"`
	IsLocked             *bool                     `json:"isLocked"`
	Neighborhoods        []NeighborhoodsType       `json:"neighborhoods"`
	MulchDeliveryConfigs []MulchDeliveryConfigType `json:"mulchDeliveryConfigs"`
	Products             []ProductType             `json:"products"`
}

////////////////////////////////////////////////////////////////////////////
//
func GetFundraiserConfig(gqlFields []string) (FrConfigType, error) {

	log.Println("Retrieving Fundraiser Config")

	frConfig := FrConfigType{}
	params := []interface{}{}
	sqlFields := []string{}

	for _, gqlField := range gqlFields {
		switch {
		case "kind" == gqlField:
			params = append(params, &frConfig.Kind)
			sqlFields = append(sqlFields, "kind")
		case "description" == gqlField:
			params = append(params, &frConfig.Description)
			sqlFields = append(sqlFields, "description")
		case "lastModifiedTime" == gqlField:
			params = append(params, &frConfig.LastModifiedTime)
			sqlFields = append(sqlFields, "last_modified_time::string")
		case "isLocked" == gqlField:
			params = append(params, &frConfig.IsLocked)
			sqlFields = append(sqlFields, "is_locked")
		case "neighborhoods" == gqlField:
			params = append(params, &frConfig.Neighborhoods)
			sqlFields = append(sqlFields, "neighborhoods::jsonb")
		case "mulchDeliveryConfigs" == gqlField:
			params = append(params, &frConfig.MulchDeliveryConfigs)
			sqlFields = append(sqlFields, "mulch_delivery_configs::jsonb")
		case "products" == gqlField:
			params = append(params, &frConfig.Products)
			sqlFields = append(sqlFields, "products::jsonb")
		default:
			return frConfig, errors.New(fmt.Sprintf("Unknown fundraiser config field: %s", gqlField))
		}

	}

	sqlCmd := fmt.Sprintf("select %s from fundraiser_config", strings.Join(sqlFields, ","))
	log.Println("SqlCmd: ", sqlCmd)
	err := Db.QueryRow(context.Background(), sqlCmd).Scan(params...)
	if err != nil {
		log.Println("Fundraiser config query failed", err)
		return FrConfigType{}, err
	}
	return frConfig, nil
}

////////////////////////////////////////////////////////////////////////////
//
func SetFundraiserConfig(frConfig FrConfigType) (bool, error) {
	frConfig.LastModifiedTime = time.Now().UTC().Format(time.RFC3339)
	log.Println("Setting Fundraiding Config: ", frConfig)

	// Reality is they need to set the entire row every time right now so
	//  this should probably be all required fields
	//  doing it this way for future when it doesn't
	values := []interface{}{}
	valIdxs := []string{}
	valIdx := 1
	sqlFields := []string{}
	if len(frConfig.Kind) != 0 {
		sqlFields = append(sqlFields, "kind")
		values = append(values, frConfig.Kind)
		valIdxs = append(valIdxs, fmt.Sprintf("$%d::string", valIdx))
		valIdx++
	}
	if len(frConfig.Description) != 0 {
		sqlFields = append(sqlFields, "description")
		values = append(values, frConfig.Description)
		valIdxs = append(valIdxs, fmt.Sprintf("$%d::string", valIdx))
		valIdx++
	}
	if len(frConfig.Neighborhoods) != 0 {
		sqlFields = append(sqlFields, "neighborhoods")
		values = append(values, frConfig.Neighborhoods)
		valIdxs = append(valIdxs, fmt.Sprintf("$%d::jsonb", valIdx))
		valIdx++
	}
	if len(frConfig.Products) != 0 {
		sqlFields = append(sqlFields, "products")
		values = append(values, frConfig.Products)
		valIdxs = append(valIdxs, fmt.Sprintf("$%d::jsonb", valIdx))
		valIdx++
	}
	if len(frConfig.MulchDeliveryConfigs) != 0 {
		sqlFields = append(sqlFields, "mulch_delivery_configs")
		values = append(values, frConfig.MulchDeliveryConfigs)
		valIdxs = append(valIdxs, fmt.Sprintf("$%d::jsonb", valIdx))
		valIdx++
	}
	if nil != frConfig.IsLocked {
		// Unfortunately hard to detect if this is set or not
		sqlFields = append(sqlFields, "is_locked")
		values = append(values, *frConfig.IsLocked)
		valIdxs = append(valIdxs, fmt.Sprintf("$%d::bool", valIdx))
		valIdx++
	}

	// Always do timestamp
	sqlFields = append(sqlFields, "last_modified_time")
	values = append(values, frConfig.LastModifiedTime)
	valIdxs = append(valIdxs, fmt.Sprintf("$%d::timestamp", valIdx))

	sqlCmd := fmt.Sprintf("insert into fundraiser_config(%s) values (%s)",
		strings.Join(sqlFields, ","), strings.Join(valIdxs, ","))

	// Start Database Operations
	trxn, err := Db.Begin(context.Background())

	log.Println("Deleting existing record")
	_, err = trxn.Exec(context.Background(), "delete from fundraiser_config")
	if err != nil {
		return false, err
	}
	log.Println("Setting Config SqlCmd: ", sqlCmd)
	_, err = trxn.Exec(context.Background(), sqlCmd, values...)
	if err != nil {
		return false, err
	}

	log.Println("About to make a commitment")
	err = trxn.Commit(context.Background())
	if err != nil {
		return false, err
	}
	return true, nil
}

////////////////////////////////////////////////////////////////////////////
//
type MulchTimecardType struct {
	Id               string
	DeliveryId       int
	LastModifiedTime time.Time
	TimeIn           string
	TimeOut          string
	TimeTotal        string
}

////////////////////////////////////////////////////////////////////////////
//
func GetMulchTimeCards(id string) []MulchTimecardType {
	doQuery := func(id string) (pgx.Rows, error) {
		if len(id) == 0 {
			log.Println("Retrieving All Timecards")
			sqlCmd := `select uid, delivery_id, last_modified_time, time_in::string, time_out::string, time_total::string from mulch_delivery_timecards`
			return Db.Query(context.Background(), sqlCmd)
		} else {
			log.Println("Retrieving Timecards for: ", id)
			sqlCmd := `select uid, delivery_id, last_modified_time, time_in::string, time_out::string, time_total::string from mulch_delivery_timecards where uid=$1`
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
