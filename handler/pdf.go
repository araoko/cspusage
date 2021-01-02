package handler

import (
	"fmt"
	"strconv"
	"time"

	myJson "github.com/araoko/cspusage/model/json"
	"github.com/araoko/cspusage/model/pdf"
)

func customerListJSON2PDFBytes(js []myJson.Customer) ([]byte, error) {
	//TODO Impliment

	billItemTable := customerListJSON2PDFTable(js)
	fPDF := pdf.GetInitPDF(true)
	pdf.PrintTable(fPDF, billItemTable, []float64{1.0, 2.0, 1.0}, nil, 0)
	return pdf.PDF2Bytes(fPDF)
}

func customerMonthyBillJSON2PDFBytes(js myJson.CustomerMonthlyBill) ([]byte, error) {
	//TODO Impliment
	custInfo := customerJSON2PDFTable(js.Owner)
	periodInfo := fmt.Sprintf("Year: %d\t\tMonth: %s", js.Date.Year, time.Month(js.Date.Month).String())
	billItemTable := billItemListJSON2PDFTable(js.LineItems)
	fPDF := pdf.GetInitPDF(true)
	pdf.PrintTable(fPDF, custInfo, []float64{1.0, 2.0}, nil, 0)
	fPDF.Ln(3.0)
	pdf.PrintText(fPDF, periodInfo)
	fPDF.Ln(-1)
	pdf.PrintTable(fPDF, billItemTable, []float64{3.0, 5.0, 1.0}, nil, 0)
	return pdf.PDF2Bytes(fPDF)
}

func customerMonthlyCostPerSubJSON2PDFBytes(js myJson.CustomerMonthlyCostPerSub) ([]byte, error) {
	custInfo := customerJSON2PDFTable(js.Owner)
	periodInfo := fmt.Sprintf("Year: %d\t\tMonth: %s", js.Date.Year, time.Month(js.Date.Month).String())
	subCostItemTable := subscriptionCostItemListJSON2PDFTable(js.LineItems)
	fPDF := pdf.GetInitPDF(true)
	pdf.PrintTable(fPDF, custInfo, []float64{1.0, 2.0}, nil, 0)
	fPDF.Ln(3.0)
	pdf.PrintText(fPDF, periodInfo)
	fPDF.Ln(-1)
	pdf.PrintTable(fPDF, subCostItemTable, []float64{4.0, 1.0}, nil, 0)
	return pdf.PDF2Bytes(fPDF)
}

func customerRangeBillJSON2PDFBytes(js myJson.CustomerRangeBill) ([]byte, error) {
	custInfo := customerJSON2PDFTable(js.Owner)
	startDateInfo := fmt.Sprintf("Start Date  -  Year: %d\t\tMonth: %s", js.DateRange.StartDate.Year, time.Month(js.DateRange.StartDate.Month).String())
	endDateInfo := fmt.Sprintf("End Date  -  Year: %d\t\tMonth: %s", js.DateRange.EndDate.Year, time.Month(js.DateRange.EndDate.Month).String())
	billItemTable := billItemListJSON2PDFTable(js.LineItems)
	fPDF := pdf.GetInitPDF(true)
	pdf.PrintTable(fPDF, custInfo, []float64{1.0, 2.0}, nil, 0)
	fPDF.Ln(3.0)
	pdf.PrintText(fPDF, startDateInfo)
	fPDF.Ln(-1)
	pdf.PrintText(fPDF, endDateInfo)
	fPDF.Ln(-1)
	pdf.PrintTable(fPDF, billItemTable, []float64{3.0, 5.0, 1.0}, nil, 0)
	return pdf.PDF2Bytes(fPDF)
}

func customerRangeCostPerSubJSON2PDFBytes(js myJson.CustomerRangeCostPerSub) ([]byte, error) {
	custInfo := customerJSON2PDFTable(js.Owner)
	startDateInfo := fmt.Sprintf("Start Date  -  Year: %d\t\tMonth: %s", js.DateRange.StartDate.Year, time.Month(js.DateRange.StartDate.Month).String())
	endDateInfo := fmt.Sprintf("End Date  -  Year: %d\t\tMonth: %s", js.DateRange.EndDate.Year, time.Month(js.DateRange.EndDate.Month).String())
	subscriptionCostItemTable := subscriptionCostItemListJSON2PDFTable(js.LineItems)
	fPDF := pdf.GetInitPDF(true)
	pdf.PrintTable(fPDF, custInfo, []float64{1.0, 2.0}, nil, 0)
	fPDF.Ln(3.0)
	pdf.PrintText(fPDF, startDateInfo)
	fPDF.Ln(-1)
	pdf.PrintText(fPDF, endDateInfo)
	fPDF.Ln(-1)
	pdf.PrintTable(fPDF, subscriptionCostItemTable, []float64{3.0, 5.0, 1.0}, nil, 0)
	return pdf.PDF2Bytes(fPDF)
}

func monthlyBillJSON2PDFBytes(js myJson.MonthlyBill) ([]byte, error) {

	periodInfo := fmt.Sprintf("Year: %d\t\tMonth: %s", js.Date.Year, time.Month(js.Date.Month).String())
	fPDF := pdf.GetInitPDF(true)
	pdf.PrintText(fPDF, periodInfo)
	fPDF.Ln(-1)

	summaryTable := monthlySummaryCostJSON2PDFTable(js.Summary.Customers)

	pdf.PrintTable(fPDF, summaryTable, []float64{3.0, 1.0}, nil, 0.0)

	for _, cmb := range js.CustomerMonthlyBills {
		fPDF.AddPage()
		custInfo := customerJSON2PDFTable(cmb.Owner)
		billItemTable := billItemListJSON2PDFTable(cmb.LineItems)
		pdf.PrintTable(fPDF, custInfo, []float64{1.0, 2.0}, nil, 0)
		fPDF.Ln(3.0)
		pdf.PrintText(fPDF, periodInfo)
		fPDF.Ln(-1)
		pdf.PrintTable(fPDF, billItemTable, []float64{3.0, 5.0, 1.0}, nil, 0)

	}
	return pdf.PDF2Bytes(fPDF)
}

func rangeBillJSON2PDFBytes(js myJson.RangeBill) ([]byte, error) {

	startDateInfo := fmt.Sprintf("Start Date  -  Year: %d\t\tMonth: %s", js.DateRange.StartDate.Year, time.Month(js.DateRange.StartDate.Month).String())
	fPDF := pdf.GetInitPDF(true)
	pdf.PrintText(fPDF, startDateInfo)
	fPDF.Ln(-1)
	endDateInfo := fmt.Sprintf("End Date  -  Year: %d\t\tMonth: %s", js.DateRange.EndDate.Year, time.Month(js.DateRange.EndDate.Month).String())
	pdf.PrintText(fPDF, endDateInfo)
	fPDF.Ln(-1)

	summaryTable := monthlySummaryCostJSON2PDFTable(js.Summary.Customers)

	pdf.PrintTable(fPDF, summaryTable, []float64{3.0, 1.0}, nil, 0.0)

	for _, cmb := range js.CustomerRangeBills {
		fPDF.AddPage()
		custInfo := customerJSON2PDFTable(cmb.Owner)
		billItemTable := billItemListJSON2PDFTable(cmb.LineItems)
		pdf.PrintTable(fPDF, custInfo, []float64{1.0, 2.0}, nil, 0)
		fPDF.Ln(3.0)
		pdf.PrintText(fPDF, startDateInfo)
		fPDF.Ln(-1)
		pdf.PrintText(fPDF, endDateInfo)
		fPDF.Ln(-1)

		pdf.PrintTable(fPDF, billItemTable, []float64{3.0, 5.0, 1.0}, nil, 0)

	}
	return pdf.PDF2Bytes(fPDF)
}

func monthlyTrendJSON2PDFBytes(js myJson.MonthlyTrend) ([]byte, error) {

	startDateInfo := fmt.Sprintf("Start Date  -  Year: %d\t\tMonth: %s", js.DateRange.StartDate.Year, time.Month(js.DateRange.StartDate.Month).String())
	fPDF := pdf.GetInitPDF(true)
	pdf.PrintText(fPDF, startDateInfo)
	fPDF.Ln(-1)
	endDateInfo := fmt.Sprintf("End Date  -  Year: %d\t\tMonth: %s", js.DateRange.EndDate.Year, time.Month(js.DateRange.EndDate.Month).String())
	pdf.PrintText(fPDF, endDateInfo)
	fPDF.Ln(-1)

	summaryTable := dateCostItemListJSON2PDFTable(js.Summary)

	pdf.PrintTable(fPDF, summaryTable, []float64{1.0, 1.0, 2.0}, nil, 0.0)

	for _, cmb := range js.Trend {
		fPDF.AddPage()
		custInfo := customerJSON2PDFTable(cmb.Owner)
		billItemTable := dateCostItemListJSON2PDFTable(cmb.Trend)
		pdf.PrintTable(fPDF, custInfo, []float64{1.0, 2.0}, nil, 0)
		fPDF.Ln(3.0)
		pdf.PrintText(fPDF, startDateInfo)
		fPDF.Ln(-1)
		pdf.PrintText(fPDF, endDateInfo)
		fPDF.Ln(-1)

		pdf.PrintTable(fPDF, billItemTable, []float64{1.0, 1.0, 2.0}, nil, 0)

	}
	return pdf.PDF2Bytes(fPDF)
}

func customerMonthlyTrendJSON2PDFBytes(js myJson.CustomerMonthlyTrend) ([]byte, error) {

	startDateInfo := fmt.Sprintf("Start Date  -  Year: %d\t\tMonth: %s", js.DateRange.StartDate.Year, time.Month(js.DateRange.StartDate.Month).String())
	custInfo := customerJSON2PDFTable(js.Owner)
	fPDF := pdf.GetInitPDF(true)
	pdf.PrintTable(fPDF, custInfo, []float64{1.0, 2.0}, nil, 0)
	fPDF.Ln(-1)
	pdf.PrintText(fPDF, startDateInfo)
	fPDF.Ln(-1)
	endDateInfo := fmt.Sprintf("End Date  -  Year: %d\t\tMonth: %s", js.DateRange.EndDate.Year, time.Month(js.DateRange.EndDate.Month).String())
	pdf.PrintText(fPDF, endDateInfo)
	fPDF.Ln(-1)

	resultTable := dateCostItemListJSON2PDFTable(js.Trend)

	pdf.PrintTable(fPDF, resultTable, []float64{4.0, 10.0, 9.0}, nil, 0.0)

	return pdf.PDF2Bytes(fPDF)
}

func customerMonthlyPerSubTrendJSON2PDFBytes(js myJson.CustomerMonthlyPerSubTrend) ([]byte, error) {

	startDateInfo := fmt.Sprintf("Start Date  -  Year: %d\t\tMonth: %s", js.DateRange.StartDate.Year, time.Month(js.DateRange.StartDate.Month).String())
	custInfo := customerJSON2PDFTable(js.Owner)
	fPDF := pdf.GetInitPDF(true)
	pdf.PrintTable(fPDF, custInfo, []float64{1.0, 2.0}, nil, 0)
	fPDF.Ln(-1)
	pdf.PrintText(fPDF, startDateInfo)
	fPDF.Ln(-1)
	endDateInfo := fmt.Sprintf("End Date  -  Year: %d\t\tMonth: %s", js.DateRange.EndDate.Year, time.Month(js.DateRange.EndDate.Month).String())
	pdf.PrintText(fPDF, endDateInfo)
	fPDF.Ln(-1)

	resultTable := dateCostSubItemListJSON2PDFTable(js.Trend)

	pdf.PrintTable(fPDF, resultTable, nil, nil, 0.0)

	return pdf.PDF2Bytes(fPDF)
}

func customerJSON2PDFTable(c myJson.Customer) pdf.PDFTable {
	b := make([][]pdf.PDFCell, 3)
	b[0] = pdf.NewPDFROW([]string{"Customer ID", c.CustomerId}, "", "")
	b[1] = pdf.NewPDFROW([]string{"Customer Company Name", c.CustomerCompanyName}, "", "")
	b[2] = pdf.NewPDFROW([]string{"Other Name(s)", c.FormerNames}, "", "")

	return pdf.PDFTable{
		Body: b,
	}
}

func customerListJSON2PDFTable(cs []myJson.Customer) pdf.PDFTable {
	h := pdf.NewPDFROW([]string{"Customer ID", "Customer Company Name", "Other Name(s)"}, "C", "B")
	b := make([][]pdf.PDFCell, len(cs))
	for i, c := range cs {
		b[i] = pdf.NewPDFROW([]string{c.CustomerId, c.CustomerCompanyName, c.FormerNames}, "L", "")
	}
	return pdf.PDFTable{
		Headers: h,
		Body:    b,
	}
}

func monthlySummaryCostJSON2PDFTable(ccs []myJson.CustomerCostItem) pdf.PDFTable {
	h := pdf.NewPDFROW([]string{"Customer Company Name", "Cost($)"}, "C", "B")
	b := make([][]pdf.PDFCell, len(ccs)+1)
	var total float32
	for i, cc := range ccs {
		bb := make([]pdf.PDFCell, 2)
		bb[0] = pdf.NewPDFCell(cc.Owner.CustomerCompanyName, "", "")
		bb[1] = pdf.NewPDFCell(strconv.FormatFloat(float64(cc.Cost), 'f', 2, 32), "R", "")
		b[i] = bb
		total += cc.Cost
	}

	bb := make([]pdf.PDFCell, 2)
	bb[0] = pdf.NewPDFCell("Total", "R", "B")
	bb[1] = pdf.NewPDFCell(strconv.FormatFloat(float64(total), 'f', 2, 32), "R", "B")
	b[len(ccs)] = bb

	return pdf.PDFTable{
		Headers: h,
		Body:    b,
	}

}

func billItemListJSON2PDFTable(bis []myJson.SubscriptionServiceCostItem) pdf.PDFTable {
	h := pdf.NewPDFROW([]string{"Subscription", "ServiceName - Type", "Cost($)"}, "C", "B")
	b := make([][]pdf.PDFCell, len(bis)+1)
	var total float32
	for i, bi := range bis {
		bb := make([]pdf.PDFCell, 3)
		bb[0] = pdf.NewPDFCell(bi.Suscription, "", "")
		bb[1] = pdf.NewPDFCell(bi.ServiceNameAndType, "", "")
		bb[2] = pdf.NewPDFCell(strconv.FormatFloat(float64(bi.Cost), 'f', 2, 32), "R", "")
		b[i] = bb
		total += bi.Cost
	}
	bb := make([]pdf.PDFCell, 3)
	bb[0] = pdf.NewPDFCell("", "", "")
	bb[1] = pdf.NewPDFCell("Total", "R", "B")
	bb[2] = pdf.NewPDFCell(strconv.FormatFloat(float64(total), 'f', 2, 32), "R", "")
	b[len(bis)] = bb
	return pdf.PDFTable{
		Headers: h,
		Body:    b,
	}
}

func dateCostItemListJSON2PDFTable(dcis []myJson.DateCostItem) pdf.PDFTable {
	h := pdf.NewPDFROW([]string{"Year", "Month", "Cost($)"}, "C", "B")
	b := make([][]pdf.PDFCell, len(dcis)+1)
	var total float32
	for i, dci := range dcis {
		bb := make([]pdf.PDFCell, 3)
		bb[0] = pdf.NewPDFCell(strconv.Itoa(dci.Date.Year), "", "")
		bb[1] = pdf.NewPDFCell(time.Month(dci.Date.Month).String(), "", "")
		bb[2] = pdf.NewPDFCell(strconv.FormatFloat(float64(dci.Cost), 'f', 2, 32), "R", "")
		b[i] = bb
		total += dci.Cost
	}
	bb := make([]pdf.PDFCell, 3)
	bb[0] = pdf.NewPDFCell("", "", "")
	bb[1] = pdf.NewPDFCell("Total", "R", "B")
	bb[2] = pdf.NewPDFCell(strconv.FormatFloat(float64(total), 'f', 2, 32), "R", "")
	b[len(dcis)] = bb
	return pdf.PDFTable{
		Headers: h,
		Body:    b,
	}
}

func dateCostSubItemListJSON2PDFTable(dcis []myJson.DateCostSubItem) pdf.PDFTable {
	cps := dcis[0].CostPerSubs
	hdrs := make([]string, len(cps)+3)
	hdrs[0] = "Year"
	hdrs[1] = "Month"
	for i, v := range cps {
		hdrs[2+i] = v.Subscription
	}
	hdrs[len(hdrs)-1] = "Total"

	h := pdf.NewPDFROW(hdrs, "C", "B")
	b := make([][]pdf.PDFCell, len(dcis)+1)
	var total float32
	colTotal := make([]float32, len(cps))
	for i, dci := range dcis {
		var rowTolal float32
		bb := make([]pdf.PDFCell, len(hdrs))
		bb[0] = pdf.NewPDFCell(strconv.Itoa(dci.Date.Year), "", "")
		bb[1] = pdf.NewPDFCell(time.Month(dci.Date.Month).String(), "", "")
		for j, v := range dci.CostPerSubs {
			cost := v.Cost
			bb[2+j] = pdf.NewPDFCell(strconv.FormatFloat(float64(cost), 'f', 2, 32), "R", "")
			rowTolal += cost
			colTotal[j] += cost
		}
		bb[len(bb)-1] = pdf.NewPDFCell(strconv.FormatFloat(float64(rowTolal), 'f', 2, 32), "R", "")
		total += rowTolal
		b[i] = bb
	}
	bb := make([]pdf.PDFCell, len(hdrs))
	bb[0] = pdf.NewPDFCell("", "", "")
	bb[1] = pdf.NewPDFCell("Total", "R", "B")
	for i, ct := range colTotal {
		bb[2+i] = pdf.NewPDFCell(strconv.FormatFloat(float64(ct), 'f', 2, 32), "R", "")
	}
	bb[len(bb)-1] = pdf.NewPDFCell(strconv.FormatFloat(float64(total), 'f', 2, 32), "R", "")
	b[len(dcis)] = bb
	return pdf.PDFTable{
		Headers: h,
		Body:    b,
	}
}

func subscriptionCostItemListJSON2PDFTable(bis []myJson.SubscriptionCostItem) pdf.PDFTable {
	h := pdf.NewPDFROW([]string{"Subscription", "Cost($)"}, "C", "B")
	b := make([][]pdf.PDFCell, len(bis)+1)
	var total float32
	for i, bi := range bis {
		bb := make([]pdf.PDFCell, 2)
		bb[0] = pdf.NewPDFCell(bi.Suscription, "", "")
		bb[1] = pdf.NewPDFCell(strconv.FormatFloat(float64(bi.Cost), 'f', 2, 32), "R", "")
		b[i] = bb
		total += bi.Cost
	}
	bb := make([]pdf.PDFCell, 2)
	bb[0] = pdf.NewPDFCell("Total", "R", "B")
	bb[1] = pdf.NewPDFCell(strconv.FormatFloat(float64(total), 'f', 2, 32), "R", "")
	b[len(bis)] = bb
	return pdf.PDFTable{
		Headers: h,
		Body:    b,
	}
}
