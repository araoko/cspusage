package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	myJson "github.com/araoko/cspusage/model/json"
	"github.com/araoko/cspusage/util"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

const (
	mediaTypeExcel = "excel"
	mediaTypePDF   = "pdf"
)

type ExportHandler struct {
	DB *sqlx.DB
}

func (h ExportHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	b64 := vars["b64"]
	var inBytes []byte
	var err error
	if len(b64) > 3 {
		inBytes, err = base64.URLEncoding.DecodeString(b64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var mType string
	hnd := req.FormValue("h")
	typ := req.FormValue("t")
	var resb []byte
	var fileName string

	switch hnd {
	case "cl":
		res, err := util.DoCustomerList(h.DB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resb, err = customerListJSON2Bytes(res, typ)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fileName = "CSP Customer List"

	case "mt":
		var in myJson.YearMonthRange
		err := json.Unmarshal(inBytes, &in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := util.DoTrend(h.DB, in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resb, err = monthlyTrendJSON2Bytes(res, typ)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fileName = fmt.Sprintf("Monthly Trend_%02d_%d_To_%02d_%d", res.DateRange.StartDate.Month, res.DateRange.StartDate.Year, res.DateRange.EndDate.Month, res.DateRange.EndDate.Year)

	case "cmt":
		var in myJson.CustomerIDAndDateRange
		err := json.Unmarshal(inBytes, &in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := util.DoCustomerTrend(h.DB, in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resb, err = customerMonthlyTrendJSON2Bytes(res, typ)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fileName = fmt.Sprintf("Monthly Trend_%s_%02d_%d_To_%02d_%d", res.Owner.CustomerCompanyName, res.DateRange.StartDate.Month, res.DateRange.StartDate.Year, res.DateRange.EndDate.Month, res.DateRange.EndDate.Year)

	case "cmtps":
		var in myJson.CustomerIDAndDateRange
		err := json.Unmarshal(inBytes, &in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := util.DoCustomerPerSubTrend(h.DB, in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resb, err = customerMonthlyPerSubTrendJSON2Bytes(util.SortCustomerPerSubTrend(res), typ)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fileName = fmt.Sprintf("Monthly Subscription Trend_%s_%02d_%d_To_%02d_%d", res.Owner.CustomerCompanyName, res.DateRange.StartDate.Month, res.DateRange.StartDate.Year, res.DateRange.EndDate.Month, res.DateRange.EndDate.Year)

	case "cmb":
		var in myJson.CustomerIDAndDate
		err := json.Unmarshal(inBytes, &in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := util.DoCustomerMonthlyBill(h.DB, in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resb, err = customerMonthyBillJSON2Bytes(res, typ)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fileName = fmt.Sprintf("Monthly Spending_%s_%02d_%d", res.Owner.CustomerCompanyName, res.Date.Month, res.Date.Year)
	case "cmbps":
		var in myJson.CustomerIDAndDate
		err := json.Unmarshal(inBytes, &in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := util.DoCustomerMonthlyCostPerSub(h.DB, in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resb, err = customerMonthlyCostPerSubJSON2Bytes(res, typ)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fileName = fmt.Sprintf("Monthly Per Subscription Spending_%s_%02d_%d", res.Owner.CustomerCompanyName, res.Date.Month, res.Date.Year)

	case "crb":
		var in myJson.CustomerIDAndDateRange
		err := json.Unmarshal(inBytes, &in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := util.DoCustomerRangeBill(h.DB, in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resb, err = customerRangeBillJSON2Bytes(res, typ)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fileName = fmt.Sprintf("Spending_%s_From_%02d_%d_To_%02d_%d", res.Owner.CustomerCompanyName, res.DateRange.StartDate.Month, res.DateRange.StartDate.Year, res.DateRange.EndDate.Month, res.DateRange.EndDate.Year)

	case "crbps":
		var in myJson.CustomerIDAndDateRange
		err := json.Unmarshal(inBytes, &in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := util.DoCustomerRangeCostPerSub(h.DB, in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resb, err = customerRangeCostPerSubJSON2Bytes(res, typ)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fileName = fmt.Sprintf("Spending Per Subscription_%s_From_%02d_%d_To_%02d_%d", res.Owner.CustomerCompanyName, res.DateRange.StartDate.Month, res.DateRange.StartDate.Year, res.DateRange.EndDate.Month, res.DateRange.EndDate.Year)

	case "mb":
		var in myJson.YearMonth
		err := json.Unmarshal(inBytes, &in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := util.DoMonthlyBill(h.DB, in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resb, err = monthyBillJSON2Bytes(res, typ)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fileName = fmt.Sprintf("Monthly Spending_%02d_%d", res.Date.Month, res.Date.Year)

	case "rb":
		var in myJson.YearMonthRange
		err := json.Unmarshal(inBytes, &in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := util.DoRangeBill(h.DB, in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resb, err = rangeBillJSON2Bytes(res, typ)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fileName = fmt.Sprintf("Spending_From_%02d_%d_To_%02d_%d", res.DateRange.StartDate.Month, res.DateRange.StartDate.Year, res.DateRange.EndDate.Month, res.DateRange.EndDate.Year)

	}

	switch typ {
	case mediaTypeExcel:
		mType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		fileName = fileName + ".xlsx"
	case mediaTypePDF:
		mType = "application/pdf"
		fileName = fileName + ".pdf"

	}

	sendFile(w, fileName, mType, resb)

}

func customerListJSON2Bytes(js []myJson.Customer, typ string) ([]byte, error) {
	switch typ {
	case mediaTypeExcel:
		return customerListJSON2ExcelBytes(js)
	case mediaTypePDF:
		return customerListJSON2PDFBytes(js)
	}
	return nil, fmt.Errorf("Unknown Media type: %s", typ)
}

func customerMonthyBillJSON2Bytes(js myJson.CustomerMonthlyBill, typ string) ([]byte, error) {
	switch typ {
	case mediaTypeExcel:
		return customerMonthlyBillJSON2ExcelBytes(js)
	case mediaTypePDF:
		return customerMonthyBillJSON2PDFBytes(js)
	}
	return nil, fmt.Errorf("Unknown Media type: %s", typ)
}

func customerMonthlyCostPerSubJSON2Bytes(js myJson.CustomerMonthlyCostPerSub, typ string) ([]byte, error) {
	switch typ {
	case mediaTypeExcel:
		return customerMonthlyCostPerSubJSON2ExcelBytes(js)
	case mediaTypePDF:
		return customerMonthlyCostPerSubJSON2PDFBytes(js)
	}
	return nil, fmt.Errorf("Unknown Media type: %s", typ)
}

func customerRangeBillJSON2Bytes(js myJson.CustomerRangeBill, typ string) ([]byte, error) {
	switch typ {
	case mediaTypeExcel:
		return customerRangeBillJSON2ExcelBytes(js)
	case mediaTypePDF:
		return customerRangeBillJSON2PDFBytes(js)
	}
	return nil, fmt.Errorf("Unknown Media type: %s", typ)
}

func customerRangeCostPerSubJSON2Bytes(js myJson.CustomerRangeCostPerSub, typ string) ([]byte, error) {
	switch typ {
	case mediaTypeExcel:
		return customerRangeCostPerSubJSON2ExcelBytes(js)
	case mediaTypePDF:
		return customerRangeCostPerSubJSON2PDFBytes(js)
	}
	return nil, fmt.Errorf("Unknown Media type: %s", typ)
}

func monthyBillJSON2Bytes(js myJson.MonthlyBill, typ string) ([]byte, error) {
	switch typ {
	case mediaTypeExcel:
		return monthlyBillJSON2ExcelBytes(js)
	case mediaTypePDF:
		return monthlyBillJSON2PDFBytes(js)
	}
	return nil, fmt.Errorf("Unknown Media type: %s", typ)
}

func rangeBillJSON2Bytes(js myJson.RangeBill, typ string) ([]byte, error) {
	switch typ {
	case mediaTypeExcel:
		return rangeBillJSON2ExcelBytes(js)
	case mediaTypePDF:
		return rangeBillJSON2PDFBytes(js)
	}
	return nil, fmt.Errorf("Unknown Media type: %s", typ)
}

func monthlyTrendJSON2Bytes(js myJson.MonthlyTrend, typ string) ([]byte, error) {
	switch typ {
	case mediaTypeExcel:
		return monthlyTrendJSON2ExcelBytes(js)
	case mediaTypePDF:
		return monthlyTrendJSON2PDFBytes(js)
	}
	return nil, fmt.Errorf("Unknown Media type: %s", typ)
}

func customerMonthlyTrendJSON2Bytes(js myJson.CustomerMonthlyTrend, typ string) ([]byte, error) {
	switch typ {
	case mediaTypeExcel:
		return customerMonthlyTrendJSON2ExcelBytes(js)
	case mediaTypePDF:
		return customerMonthlyTrendJSON2PDFBytes(js)
	}
	return nil, fmt.Errorf("Unknown Media type: %s", typ)
}

func customerMonthlyPerSubTrendJSON2Bytes(js myJson.CustomerMonthlyPerSubTrend, typ string) ([]byte, error) {
	switch typ {
	case mediaTypeExcel:
		return customerMonthlyPerSubTrendJSON2ExcelBytes(js)
	case mediaTypePDF:
		return customerMonthlyPerSubTrendJSON2PDFBytes(js)
	}
	return nil, fmt.Errorf("Unknown Media type: %s", typ)
}

func sendFile(w http.ResponseWriter, fileName string, mType string, fileContent []byte) {
	//"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Header().Set("Content-Type", mType)
	w.Header().Set("Content-Length", strconv.Itoa(len(fileContent)))
	w.Write(fileContent)
}
