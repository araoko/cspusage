package entity

import (
	"fmt"
	"strconv"

	"github.com/araoko/cspusage/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type SubscriptionServiceCostItem struct {
	Subscription string  `db:"Subscription"`
	Service      string  `db:"Service Name - Type"`
	Cost         float32 `db:"Cost"`
}

type SubscriptionCostItem struct {
	Subscription string `db:"Subscription"`
	//Service      string  `db:"Service Name - Type"`
	Cost float32 `db:"Cost"`
}

func GetCustomerMonthlyBill(db *sqlx.DB, cid string, yr, mo int) ([]SubscriptionServiceCostItem, error) {
	ss := []SubscriptionServiceCostItem{}
	q := "CALL sp_customer_monthly2(?, ?, ?)"
	err := db.Select(&ss, q, cid, yr, mo)
	if err != nil {
		return nil, err
	}

	return ss, nil

}

//GetCustomerMonthlyCostPerSub returns the azure usage cost for a cutomer per subscription
func GetCustomerMonthlyCostPerSub(db *sqlx.DB, cid string, yr, mo int) ([]SubscriptionCostItem, error) {
	ss := []SubscriptionCostItem{}
	q := "CALL sp_customer_monthly_per_sub(?, ?, ?)"
	err := doQuerySelect(db, &ss, q, cid, yr, mo)
	if err != nil {
		return nil, err
	}

	return ss, nil

}

func GetCustomerRangeBill(db *sqlx.DB, cid string, yr, smo, emo int) ([]SubscriptionServiceCostItem, error) {
	ss := []SubscriptionServiceCostItem{}
	q := "CALL sp_customer_monthlyR(?, ?, ?, ?)"
	err := db.Select(&ss, q, cid, yr, smo, emo)
	if err != nil {
		return nil, err
	}

	return ss, nil
}

func GetCustomerRangeCostPerSub(db *sqlx.DB, cid string, yr, smo, emo int) ([]SubscriptionCostItem, error) {
	ss := []SubscriptionCostItem{}
	q := "CALL sp_customer_range_per_sub(?, ?, ?, ?)"
	err := doQuerySelect(db, &ss, q, cid, yr, smo, emo)
	if err != nil {
		return nil, err
	}

	return ss, nil
}

func doQuerySelect(db *sqlx.DB, destination interface{}, queryString string, args ...interface{}) error {
	return db.Select(destination, queryString, args...)
}

type Customer struct {
	CustomerID          string `db:"CustomerId"`
	CustomerCompanyName string `db:"CustomerCompanyName"`
	FormerNames         string `db:"FormerNames"`
}

func GetCustomers(db *sqlx.DB) ([]Customer, error) {

	ss := []Customer{}
	q := "SELECT * FROM customer"
	err := db.Select(&ss, q)
	if err != nil {
		return nil, err
	}

	return ss, nil

}

func GetCustomerFromID(db *sqlx.DB, customerID string) (Customer, error) {
	s := Customer{}
	q := "SELECT * FROM customer WHERE CustomerId = ?"
	err := db.Get(&s, q, customerID)
	return s, err
}

type CustomerIDCostItem struct {
	CustomerID string  `db:"CustomerId"`
	Cost       float32 `db:"Cost"`
}

func GetMonthlyCostSummary(db *sqlx.DB, yr, mo int) ([]CustomerIDCostItem, error) {
	ss := []CustomerIDCostItem{}
	q := "CALL sp_monthly_cost_summary(?, ?)"
	err := db.Select(&ss, q, yr, mo)
	if err != nil {
		return nil, err
	}

	return ss, nil
}

type CustomerIDYearMonthCostItem struct {
	CustomerID string  `db:"CustomerId"`
	Year       int     `db:"Year"`
	Month      int     `db:"Month"`
	Cost       float32 `db:"Cost"`
}

type YearMonthCostItem struct {
	Year  int     `db:"Year"`
	Month int     `db:"Month"`
	Cost  float32 `db:"Cost"`
}

type YearMonthSubscriptionCostItem struct {
	Year         int     `db:"Year"`
	Month        int     `db:"Month"`
	Subscription string  `db:"Subscription"`
	Cost         float32 `db:"Cost"`
}

func GetMonthlyTrend(db *sqlx.DB, yr, smo, emo int) ([]CustomerIDYearMonthCostItem, error) {
	ss := []CustomerIDYearMonthCostItem{}
	q := "CALL sp_monthly_trend(?, ?, ?)"
	err := doQuerySelect(db, &ss, q, yr, smo, emo)
	if err != nil {
		return nil, err
	}

	return ss, nil
}

func GetCustomerMonthlyTrend(db *sqlx.DB, cid string, yr, smo, emo int) ([]YearMonthCostItem, error) {
	ss := []YearMonthCostItem{}
	q := "CALL sp_customer_monthly_trend(?, ?, ?, ?)"
	err := doQuerySelect(db, &ss, q, cid, yr, smo, emo)
	if err != nil {
		return nil, err
	}

	return ss, nil
}

func GetCustomerPerSubMonthlyTrend(db *sqlx.DB, cid string, yr, smo, emo int) ([]YearMonthSubscriptionCostItem, error) {
	ss := []YearMonthSubscriptionCostItem{}
	q := "CALL sp_customer_sub_monthly_trend(?, ?, ?, ?)"
	err := doQuerySelect(db, &ss, q, cid, yr, smo, emo)
	if err != nil {
		return nil, err
	}

	return ss, nil
}

type YearMonth struct {
	Year  int `db:"Year"`
	Month int `db:"Month"`
}

func GetStartYearMonth(db *sqlx.DB) (int, int, error) {
	q := "select MinYear() AS Year, MinMonth(MinYear()) As Month"

	ss := YearMonth{}
	err := db.Select(&ss, q)
	if err != nil {
		return 0, 0, err
	}

	return ss.Year, ss.Month, nil
}

func GetEndYearMonth(db *sqlx.DB) (int, int, error) {
	q := "select MaxYear() AS Year, MaxMonth(MaxYear()) As Month"

	ss := YearMonth{}
	err := db.Select(&ss, q)
	if err != nil {
		return 0, 0, err
	}

	return ss.Year, ss.Month, nil

}

func GetMySQLDB(c *config.Config) (db *sqlx.DB, err error) {

	return connect(c.DBUserName, c.DBPassword, c.DBHostName, strconv.Itoa(c.DBPort), c.DBName)
}

func connect(user, password, host, port, database string) (*sqlx.DB, error) {
	str := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true", user, password, host, port, database)
	return sqlx.Open("mysql", str)
}
