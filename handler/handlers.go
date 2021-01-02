package handler

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/araoko/cspusage/entity"
	myJson "github.com/araoko/cspusage/model/json"
	"github.com/araoko/cspusage/util"
	"github.com/jmoiron/sqlx"
)

//CustomerListHandler ...
type CustomerListHandler struct {
	DB *sqlx.DB
}

func (h CustomerListHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//log.Println("got customer list request")

	customers, err := util.DoCustomerList(h.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, customers)
}

//CustomerMonthlyBillHandler ...
type CustomerMonthlyBillHandler struct {
	DB *sqlx.DB
}

func (h CustomerMonthlyBillHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	data := myJson.CustomerIDAndDate{}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := util.DoCustomerMonthlyBill(h.DB, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, result)

}

//CustomerMonthlyCostPerSubHandler handles requst for usage details
//for a specified account in a specified month summing the cost per subscription
type CustomerMonthlyCostPerSubHandler struct {
	DB *sqlx.DB
}

func (h CustomerMonthlyCostPerSubHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	data := myJson.CustomerIDAndDate{}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := util.DoCustomerMonthlyCostPerSub(h.DB, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, result)

}

type RangeBillHandler struct {
	DB *sqlx.DB
}

func (h RangeBillHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	data := myJson.YearMonthRange{}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := util.DoRangeBill(h.DB, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, resp)
}

//CustomerRangeCostPerSubHandler handles requst for usage details
//for a specified account for a specified date range summing the cost per subscription
type CustomerRangeCostPerSubHandler struct {
	DB *sqlx.DB
}

func (h CustomerRangeCostPerSubHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	data := myJson.CustomerIDAndDateRange{}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := util.DoCustomerRangeCostPerSub(h.DB, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, resp)
}

type CustomerRangeBillHandler struct {
	DB *sqlx.DB
}

func (h CustomerRangeBillHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	data := myJson.CustomerIDAndDateRange{}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := util.DoCustomerRangeBill(h.DB, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, resp)
}

type MonthlyBillHandler struct {
	DB *sqlx.DB
}

func (h MonthlyBillHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	data := myJson.YearMonth{}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	mb, err := util.DoMonthlyBill(h.DB, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSON(w, mb)

}

type MonthlySummaryHandler struct {
	DB *sqlx.DB
}

func (h MonthlySummaryHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	data := myJson.YearMonth{}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := util.DoMonthlySummary(h.DB, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, res)

}

type CustomerMonthlyTrendPerSubHandler struct {
	DB *sqlx.DB
}

func (h CustomerMonthlyTrendPerSubHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	data := myJson.CustomerIDAndDateRange{}

	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := util.DoCustomerPerSubTrend(h.DB, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, util.SortCustomerPerSubTrend(res))

}

type CustomerMonthlyTrendHandler struct {
	DB *sqlx.DB
}

func (h CustomerMonthlyTrendHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	data := myJson.CustomerIDAndDateRange{}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := util.DoCustomerTrend(h.DB, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, util.SortCustomerTrend(res))

}

type MonthlyTrendHandler struct {
	DB *sqlx.DB
}

func (h MonthlyTrendHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	data := myJson.YearMonthRange{}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := util.DoTrend(h.DB, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, res)

}

type ADLoginHandler struct {
	Auth *entity.ADAuthenticator
}

func (h ADLoginHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	data := myJson.Login{}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.LoggedIn = 0
	fmt.Printf("Auth cred recived: %v\n", data)
	h.Auth.LoadCred(data.UserName, data.Password)
	status, err := h.Auth.ADAuthenticate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data.LoggedIn = b2i(status)
	fmt.Printf("Post Auth cred : %v\n", data)
	data.Password = "**********"
	sendJSON(w, data)

}

func b2i(b bool) int {
	return (map[bool]int{true: 1, false: 0})[b]
}
