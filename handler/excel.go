package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/araoko/cspusage/model/excel"
	myJson "github.com/araoko/cspusage/model/json"
)

func customerListJSON2ExcelBytes(js []myJson.Customer) ([]byte, error) {
	sheetName := excel.FitSheetName("Account List")
	x := initExcel(sheetName)
	a := "B3"

	tRange, err := customerListJSON2Excel(x, sheetName, a, js)
	if err != nil {
		return nil, err
	}

	excel.FitColWidth(x, sheetName, tRange, 1.2)

	return excel2Bytes(x)

}

func customerMonthlyBillJSON2ExcelBytes(js myJson.CustomerMonthlyBill) ([]byte, error) {
	sheetName := excel.FitSheetName(js.Owner.CustomerCompanyName)
	x := initExcel(sheetName)
	a := "B3"
	_, err := customerJSON2Excel(x, sheetName, a, js.Owner)
	if err != nil {
		return nil, err
	}

	a, err = excel.OffsetJump(a, 0, 5)
	if err != nil {
		return nil, err
	}

	_, err = writeDate(x, sheetName, a, js.Date)
	if err != nil {
		return nil, err
	}

	a, err = excel.OffsetJump(a, 0, 2)

	tRange, err := billItemListJSON2Excel(x, sheetName, a, js.LineItems)
	if err != nil {
		return nil, err
	}

	excel.FitColWidth(x, sheetName, tRange, 1.2)

	return excel2Bytes(x)

}

func customerMonthlyCostPerSubJSON2ExcelBytes(js myJson.CustomerMonthlyCostPerSub) ([]byte, error) {
	sheetName := excel.FitSheetName(js.Owner.CustomerCompanyName)
	x := initExcel(sheetName)
	a := "B3"
	cRange, err := customerJSON2Excel(x, sheetName, a, js.Owner)
	if err != nil {
		return nil, err
	}

	a, err = excel.OffsetJump(a, 0, 5)
	if err != nil {
		return nil, err
	}

	_, err = writeDate(x, sheetName, a, js.Date)
	if err != nil {
		return nil, err
	}

	a, err = excel.OffsetJump(a, 0, 2)

	_, err = subscriptionCostItemListJSON2Excel(x, sheetName, a, js.LineItems)
	if err != nil {
		return nil, err
	}

	excel.FitColWidth(x, sheetName, cRange, 1.2)

	return excel2Bytes(x)

}

func customerRangeBillJSON2ExcelBytes(js myJson.CustomerRangeBill) ([]byte, error) {
	sheetName := excel.FitSheetName(js.Owner.CustomerCompanyName)
	x := initExcel(sheetName)
	a := "B3"
	_, err := customerJSON2Excel(x, sheetName, a, js.Owner)
	if err != nil {
		return nil, err
	}

	a, err = excel.OffsetJump(a, 0, 5)
	if err != nil {
		return nil, err
	}

	_, err = writeDateRange(x, sheetName, a, js.DateRange)
	if err != nil {
		return nil, err
	}

	a, err = excel.OffsetJump(a, 0, 3)

	tRange, err := billItemListJSON2Excel(x, sheetName, a, js.LineItems)
	if err != nil {
		return nil, err
	}

	excel.FitColWidth(x, sheetName, tRange, 1.2)

	return excel2Bytes(x)

}

func customerRangeCostPerSubJSON2ExcelBytes(js myJson.CustomerRangeCostPerSub) ([]byte, error) {
	sheetName := excel.FitSheetName(js.Owner.CustomerCompanyName)
	x := initExcel(sheetName)
	a := "B3"
	cRange, err := customerJSON2Excel(x, sheetName, a, js.Owner)
	if err != nil {
		return nil, err
	}

	a, err = excel.OffsetJump(a, 0, 5)
	if err != nil {
		return nil, err
	}

	_, err = writeDateRange(x, sheetName, a, js.DateRange)
	if err != nil {
		return nil, err
	}

	a, err = excel.OffsetJump(a, 0, 3)

	_, err = subscriptionCostItemListJSON2Excel(x, sheetName, a, js.LineItems)
	if err != nil {
		return nil, err
	}

	excel.FitColWidth(x, sheetName, cRange, 1.2)

	return excel2Bytes(x)

}

func monthlyBillJSON2ExcelBytes(js myJson.MonthlyBill) ([]byte, error) {
	sheetName := "Summary"
	x := initExcel(sheetName)
	a := "B3"

	_, err := writeDate(x, sheetName, a, js.Date)
	if err != nil {
		return nil, err
	}

	a, err = excel.OffsetJump(a, 0, 2)

	sRange, err := monthlySummaryCostJSON2Excel(x, sheetName, a, js.Summary.Customers)
	if err != nil {
		return nil, err
	}
	excel.FitColWidth(x, sheetName, sRange, 1.2)

	for _, cmb := range js.CustomerMonthlyBills {
		a = "B3"
		sheetName = excel.FitSheetName(cmb.Owner.CustomerCompanyName)
		x.NewSheet(sheetName)
		_, err := customerJSON2Excel(x, sheetName, a, cmb.Owner)
		if err != nil {
			return nil, err
		}

		a, err = excel.OffsetJump(a, 0, 5)
		if err != nil {
			return nil, err
		}
		tRange, err := billItemListJSON2Excel(x, sheetName, a, cmb.LineItems)
		if err != nil {
			return nil, err
		}

		excel.FitColWidth(x, sheetName, tRange, 1.2)

	}

	return excel2Bytes(x)
}

func rangeBillJSON2ExcelBytes(js myJson.RangeBill) ([]byte, error) {
	sheetName := "Summary"
	x := initExcel(sheetName)
	a := "B3"

	_, err := writeDateRange(x, sheetName, a, js.DateRange)
	if err != nil {
		return nil, err
	}

	a, err = excel.OffsetJump(a, 0, 3)
	sRange, err := monthlySummaryCostJSON2Excel(x, sheetName, a, js.Summary.Customers)
	if err != nil {
		return nil, err
	}
	excel.FitColWidth(x, sheetName, sRange, 1.2)

	for _, cmb := range js.CustomerRangeBills {
		a = "B3"
		sheetName = excel.FitSheetName(cmb.Owner.CustomerCompanyName)
		x.NewSheet(sheetName)
		_, err := customerJSON2Excel(x, sheetName, a, cmb.Owner)
		if err != nil {
			return nil, err
		}

		a, err = excel.OffsetJump(a, 0, 5)
		if err != nil {
			return nil, err
		}
		tRange, err := billItemListJSON2Excel(x, sheetName, a, cmb.LineItems)
		if err != nil {
			return nil, err
		}

		excel.FitColWidth(x, sheetName, tRange, 1.2)

	}

	return excel2Bytes(x)
}

func monthlyTrendJSON2ExcelBytes(js myJson.MonthlyTrend) ([]byte, error) {
	cht := chart{
		title:  fmt.Sprintf("CSP Monthly Trend %v %d to %v %d", time.Month(js.DateRange.StartDate.Month), js.DateRange.StartDate.Year, time.Month(js.DateRange.EndDate.Month), js.DateRange.EndDate.Year),
		cType:  "line",
		width:  1200,
		height: 480,
		lines:  make([]chartlet, 0),
	}
	sheetName := "Summary"
	x := initExcel(sheetName)
	a := "B3"
	chartSheetName := "Chart"
	x.NewSheet(chartSheetName)
	_, err := writeDateRange(x, sheetName, a, js.DateRange)
	if err != nil {
		return nil, err
	}

	a, err = excel.OffsetJump(a, 0, 3)
	sRange, err := dateCostItemListJSON2Excel(x, sheetName, a, js.Summary)
	if err != nil {
		return nil, err
	}

	top, bot := excel.R2c(sRange)
	topbuff, err := excel.OffsetJump(top, 0, 1)
	if err != nil {
		return nil, err
	}
	botbuff, err := excel.OffsetJump(bot, -1, -1)
	if err != nil {
		return nil, err
	}
	chartNameLabel, err := excel.OffsetJump(bot, -1, 0)
	if err != nil {
		return nil, err
	}
	chartNameLabel, err = excel.Dollarize(chartNameLabel, true, true)
	if err != nil {
		return nil, err
	}

	topbuff, err = excel.Dollarize(topbuff, true, true)
	if err != nil {
		return nil, err
	}

	botbuff, err = excel.Dollarize(botbuff, true, true)
	if err != nil {
		return nil, err
	}
	sChart := chartlet{
		name:     sheetName + "!" + chartNameLabel,
		catRange: sheetName + "!" + topbuff + ":" + botbuff,
	}

	topbuff, err = excel.OffsetJump(top, 2, 1)
	if err != nil {
		return nil, err
	}
	botbuff, err = excel.OffsetJump(bot, 0, -1)
	if err != nil {
		return nil, err
	}

	topbuff, err = excel.Dollarize(topbuff, true, true)
	if err != nil {
		return nil, err
	}

	botbuff, err = excel.Dollarize(botbuff, true, true)
	if err != nil {
		return nil, err
	}

	sChart.valRange = sheetName + "!" + topbuff + ":" + botbuff
	cht.lines = append(cht.lines, sChart)

	excel.FitColWidth(x, sheetName, sRange, 2.0)

	for _, cmb := range js.Trend {
		a = "B3"
		chartNameLabel, _ = excel.OffsetJump(a, 1, 1)
		sheetName = excel.FitSheetName(cmb.Owner.CustomerCompanyName)
		x.NewSheet(sheetName)
		cRange, err := customerJSON2Excel(x, sheetName, a, cmb.Owner)
		if err != nil {
			return nil, err
		}

		a, err = excel.OffsetJump(a, 0, 5)
		if err != nil {
			return nil, err
		}
		tRange, err := dateCostItemListJSON2Excel(x, sheetName, a, cmb.Trend)
		if err != nil {
			return nil, err
		}

		top, bot = excel.R2c(tRange)
		topbuff, err = excel.OffsetJump(top, 0, 1)
		if err != nil {
			return nil, err
		}
		botbuff, err = excel.OffsetJump(bot, -1, -1)
		if err != nil {
			return nil, err
		}

		chartNameLabel, err = excel.Dollarize(chartNameLabel, true, true)
		if err != nil {
			return nil, err
		}

		topbuff, err = excel.Dollarize(topbuff, true, true)
		if err != nil {
			return nil, err
		}

		botbuff, err = excel.Dollarize(botbuff, true, true)
		if err != nil {
			return nil, err
		}

		sChart = chartlet{
			name:     sheetName + "!" + chartNameLabel,
			catRange: sheetName + "!" + topbuff + ":" + botbuff,
		}

		topbuff, err = excel.OffsetJump(top, 2, 1)
		if err != nil {
			return nil, err
		}
		botbuff, err = excel.OffsetJump(bot, 0, -1)
		if err != nil {
			return nil, err
		}

		topbuff, err = excel.Dollarize(topbuff, true, true)
		if err != nil {
			return nil, err
		}

		botbuff, err = excel.Dollarize(botbuff, true, true)
		if err != nil {
			return nil, err
		}

		sChart.valRange = sheetName + "!" + topbuff + ":" + botbuff
		cht.lines = append(cht.lines, sChart)

		excel.FitColWidth(x, sheetName, tRange, 1.2)
		excel.FitColWidth(x, sheetName, cRange, 1.2)

	}

	err = x.AddChart(chartSheetName, "B3", cht.doString())
	if err != nil {
		fmt.Println("Chart Error:::", err.Error())
	}

	return excel2Bytes(x)
}

func customerMonthlyTrendJSON2ExcelBytes(js myJson.CustomerMonthlyTrend) ([]byte, error) {
	cht := chart{
		title:  fmt.Sprintf("CSP Monthly Trend - %s - %v %d to %v %d", js.Owner.CustomerCompanyName, time.Month(js.DateRange.StartDate.Month), js.DateRange.StartDate.Year, time.Month(js.DateRange.EndDate.Month), js.DateRange.EndDate.Year),
		cType:  "line",
		width:  1200,
		height: 480,
		lines:  make([]chartlet, 0),
	}
	sheetName := "Data"
	x := initExcel(sheetName)
	a := "B3"
	chartSheetName := "Chart"
	x.NewSheet(chartSheetName)
	_, err := writeDateRange(x, sheetName, a, js.DateRange)
	if err != nil {
		return nil, err
	}

	a, err = excel.OffsetJump(a, 0, 3)
	sRange, err := dateCostItemListJSON2Excel(x, sheetName, a, js.Trend)
	if err != nil {
		return nil, err
	}

	top, bot := excel.R2c(sRange)
	topbuff, err := excel.OffsetJump(top, 0, 1)
	if err != nil {
		return nil, err
	}
	botbuff, err := excel.OffsetJump(bot, -1, -1)
	if err != nil {
		return nil, err
	}
	chartNameLabel, err := excel.OffsetJump(top, 2, 0)
	if err != nil {
		return nil, err
	}
	chartNameLabel, err = excel.Dollarize(chartNameLabel, true, true)
	if err != nil {
		return nil, err
	}

	topbuff, err = excel.Dollarize(topbuff, true, true)
	if err != nil {
		return nil, err
	}

	botbuff, err = excel.Dollarize(botbuff, true, true)
	if err != nil {
		return nil, err
	}
	sChart := chartlet{
		name:     sheetName + "!" + chartNameLabel,
		catRange: sheetName + "!" + topbuff + ":" + botbuff,
	}

	topbuff, err = excel.OffsetJump(top, 2, 1)
	if err != nil {
		return nil, err
	}
	botbuff, err = excel.OffsetJump(bot, 0, -1)
	if err != nil {
		return nil, err
	}

	topbuff, err = excel.Dollarize(topbuff, true, true)
	if err != nil {
		return nil, err
	}

	botbuff, err = excel.Dollarize(botbuff, true, true)
	if err != nil {
		return nil, err
	}

	sChart.valRange = sheetName + "!" + topbuff + ":" + botbuff
	cht.lines = append(cht.lines, sChart)

	excel.FitColWidth(x, sheetName, sRange, 2.0)

	err = x.AddChart(chartSheetName, "B3", cht.doString())
	if err != nil {
		fmt.Println("Chart Error:::", err.Error())
	}

	return excel2Bytes(x)
}

func customerMonthlyPerSubTrendJSON2ExcelBytes(js myJson.CustomerMonthlyPerSubTrend) ([]byte, error) {
	subLen := len(js.Trend[0].CostPerSubs)
	cht := chart{
		title:  fmt.Sprintf("CSP Monthly Subscription Trend - %s - %v %d to %v %d", js.Owner.CustomerCompanyName, time.Month(js.DateRange.StartDate.Month), js.DateRange.StartDate.Year, time.Month(js.DateRange.EndDate.Month), js.DateRange.EndDate.Year),
		cType:  "line",
		width:  1200,
		height: 480,
		lines:  make([]chartlet, 0),
	}
	sheetName := "Data"
	x := initExcel(sheetName)
	a := "B3"
	chartSheetName := "Chart"
	x.NewSheet(chartSheetName)
	_, err := writeDateRange(x, sheetName, a, js.DateRange)
	if err != nil {
		return nil, err
	}

	a, err = excel.OffsetJump(a, 0, 3)
	sRange, err := dateCostSubItemListJSON2Excel(x, sheetName, a, js.Trend)
	if err != nil {
		return nil, err
	}

	top, bot := excel.R2c(sRange)
	topbuff, err := excel.OffsetJump(top, 0, 1)
	if err != nil {
		return nil, err
	}
	botbuff, err := excel.OffsetJump(bot, -(1 + subLen), -1)
	if err != nil {
		return nil, err
	}

	catRange := sheetName + "!" + topbuff + ":" + botbuff

	for i := 0; i <= subLen; i++ {
		chartNameLabel, err := excel.OffsetJump(top, 2+i, 0)
		if err != nil {
			return nil, err
		}

		sChart := chartlet{
			name:     sheetName + "!" + chartNameLabel,
			catRange: catRange,
		}
		topBuff, err := excel.OffsetJump(top, 2+i, 1)
		botBuff, err := excel.OffsetJump(bot, i-subLen, -1)

		sChart.valRange = sheetName + "!" + topBuff + ":" + botBuff
		cht.lines = append(cht.lines, sChart)
	}

	excel.FitColWidth(x, sheetName, sRange, 2.0)

	err = x.AddChart(chartSheetName, "B3", cht.doString())
	if err != nil {
		fmt.Println("Chart Error:::", err.Error())
	}

	return excel2Bytes(x)
}

func billItemListJSON2Excel(x *excelize.File, sheetName string, anchor string,
	bis []myJson.SubscriptionServiceCostItem) (string, error) {
	t, err := excel.NewColTable(nil)
	if err != nil {
		return "", err
	}

	totalCellStyle, err := x.
		NewStyle(`{"alignment":{"horizontal":"right"},"font":{"bold":true}}`)
	if err != nil {
		return "", err
	}

	colSubscripton, _ := excel.NewColumn("Subscription", nil, excel.StyleNone,
		excel.StyleNone)
	colService, _ := excel.NewColumn("Service Name - Type", nil,
		excel.StyleNone, excel.StyleNone)
	colCost, _ := excel.NewColumn("Cost($)", nil, excel.StyleNone,
		excel.StyleNone)

	for _, mbi := range bis {
		cSub, err := excel.NewCell(excel.CellTypeString, mbi.Suscription,
			excel.StyleNone)
		if err != nil {
			return "", err
		}
		colSubscripton.AppendCell(cSub)

		cSvc, err := excel.NewCell(excel.CellTypeString, mbi.ServiceNameAndType,
			excel.StyleNone)
		if err != nil {
			return "", err
		}
		colService.AppendCell(cSvc)

		cCst, err := excel.NewCell(excel.CellTypeFloat, mbi.Cost,
			excel.StyleNone)
		if err != nil {
			return "", err
		}
		colCost.AppendCell(cCst)
		//sum += mbi.Cost

	}

	totalLabelCell, err := excel.NewCell(excel.CellTypeString, "Total",
		totalCellStyle)
	if err != nil {
		return "", err
	}
	colService.AppendCell(totalLabelCell)

	t.AppendCol(colSubscripton)
	t.AppendCol(colService)
	t.AppendCol(colCost)

	tableFormat := `
	{"table_style":"TableStyleMedium4", "show_last_column":true,
	"show_row_stripes":true,"show_column_stripes":false}`

	t.SetTableFormat(tableFormat)
	si := x.GetSheetIndex(sheetName)
	if si == 0 {
		_ = x.NewSheet(sheetName)
	}
	tRange, err := t.ExcelizeWrite(x, sheetName, anchor)
	if err != nil {
		return "", err
	}
	//err = excel.FitColWidth(x, sheetName, tRange, 1.2)
	//if err != nil {
	//	return err
	//}
	_, lpos := excel.R2c(tRange)

	lposCn, lposRn, err := excelize.CellNameToCoordinates(lpos)
	if err != nil {
		return tRange, err
	}
	_, anchanchorRn, err := excelize.CellNameToCoordinates(anchor)
	if err != nil {
		return tRange, err
	}

	costColDataTop, err := excelize.CoordinatesToCellName(lposCn, anchanchorRn+1)
	if err != nil {
		return tRange, err
	}

	costColDataBottom, err := excelize.CoordinatesToCellName(lposCn, lposRn-1)
	if err != nil {
		return tRange, err
	}

	sumFormula := fmt.Sprintf("SUM(%s:%s)", costColDataTop, costColDataBottom)
	totalValueCell, err := excel.NewCell(excel.CellTypeFormula, sumFormula,
		totalCellStyle)
	if err != nil {
		return tRange, err
	}
	_, err = totalValueCell.ExcelizeWrite(x, sheetName, lpos)
	return tRange, err

}

func dateCostItemListJSON2Excel(x *excelize.File, sheetName string,
	anchor string, dcis []myJson.DateCostItem) (string, error) {
	t, err := excel.NewColTable(nil)
	if err != nil {
		return "", err
	}

	totalCellStyle, err := x.
		NewStyle(`{"alignment":{"horizontal":"right"},"font":{"bold":true}}`)
	if err != nil {
		return "", err
	}

	colYear, _ := excel.NewColumn("Year", nil, excel.StyleNone,
		excel.StyleNone)
	colMonth, _ := excel.NewColumn("Month", nil, excel.StyleNone,
		excel.StyleNone)
	colCost, _ := excel.NewColumn("Cost($)", nil, excel.StyleNone,
		excel.StyleNone)

	for _, mbi := range dcis {
		cYr, err := excel.NewCell(excel.CellTypeInt, mbi.Date.Year,
			excel.StyleNone)
		if err != nil {
			return "", err
		}
		colYear.AppendCell(cYr)

		cMon, err := excel.NewCell(excel.CellTypeString,
			time.Month(mbi.Date.Month).String(),
			excel.StyleNone)
		if err != nil {
			return "", err
		}
		colMonth.AppendCell(cMon)

		cCst, err := excel.NewCell(excel.CellTypeFloat, mbi.Cost,
			excel.StyleNone)
		if err != nil {
			return "", err
		}
		colCost.AppendCell(cCst)
		//sum += mbi.Cost

	}

	totalLabelCell, err := excel.NewCell(excel.CellTypeString, "Total",
		totalCellStyle)
	if err != nil {
		return "", err
	}
	colMonth.AppendCell(totalLabelCell)

	t.AppendCol(colYear)
	t.AppendCol(colMonth)
	t.AppendCol(colCost)

	tableFormat := `
	{"table_style":"TableStyleMedium4", "show_last_column":true,
	"show_row_stripes":true,"show_column_stripes":false}`

	t.SetTableFormat(tableFormat)
	si := x.GetSheetIndex(sheetName)
	if si == 0 {
		_ = x.NewSheet(sheetName)
	}
	tRange, err := t.ExcelizeWrite(x, sheetName, anchor)
	if err != nil {
		return "", err
	}
	//err = excel.FitColWidth(x, sheetName, tRange, 1.2)
	//if err != nil {
	//	return err
	//}
	_, lpos := excel.R2c(tRange)

	lposCn, lposRn, err := excelize.CellNameToCoordinates(lpos)
	if err != nil {
		return tRange, err
	}
	_, anchanchorRn, err := excelize.CellNameToCoordinates(anchor)
	if err != nil {
		return tRange, err
	}

	costColDataTop, err := excelize.CoordinatesToCellName(lposCn,
		anchanchorRn+1)
	if err != nil {
		return tRange, err
	}

	costColDataBottom, err := excelize.CoordinatesToCellName(lposCn, lposRn-1)
	if err != nil {
		return tRange, err
	}

	sumFormula := fmt.Sprintf("SUM(%s:%s)", costColDataTop, costColDataBottom)
	totalValueCell, err := excel.NewCell(excel.CellTypeFormula, sumFormula,
		totalCellStyle)
	if err != nil {
		return tRange, err
	}
	_, err = totalValueCell.ExcelizeWrite(x, sheetName, lpos)
	return tRange, err

}

func dateCostSubItemListJSON2Excel(x *excelize.File, sheetName string,
	anchor string, dcis []myJson.DateCostSubItem) (string, error) {
	t, err := excel.NewColTable(nil)
	if err != nil {
		return "", err
	}

	totalCellStyle, err := x.
		NewStyle(`{"alignment":{"horizontal":"right"},"font":{"bold":true}}`)
	if err != nil {
		return "", err
	}

	colYear, _ := excel.NewColumn("Year", nil, excel.StyleNone,
		excel.StyleNone)
	colMonth, _ := excel.NewColumn("Month", nil, excel.StyleNone,
		excel.StyleNone)

	subCols := make([]excel.Column, len(dcis[0].CostPerSubs))
	for i, v := range dcis[0].CostPerSubs {
		colSub, _ := excel.NewColumn(v.Subscription, nil, excel.StyleNone, excel.StyleNone)
		subCols[i] = colSub
	}
	colTotal, _ := excel.NewColumn("Total", nil, excel.StyleNone,
		excel.StyleNone)

	for _, mbi := range dcis {
		cYr, err := excel.NewCell(excel.CellTypeInt, mbi.Date.Year,
			excel.StyleNone)
		if err != nil {
			return "", err
		}
		colYear.AppendCell(cYr)

		cMon, err := excel.NewCell(excel.CellTypeString,
			time.Month(mbi.Date.Month).String(),
			excel.StyleNone)
		if err != nil {
			return "", err
		}
		colMonth.AppendCell(cMon)
		cstSubs := mbi.CostPerSubs
		var rowTotal float32
		for i, v := range cstSubs {
			cst := v.Cost
			cCst, err := excel.NewCell(excel.CellTypeFloat, cst,
				excel.StyleNone)
			if err != nil {
				return "", err
			}
			subCols[i].AppendCell(cCst)
			rowTotal += cst
		}

		rCst, err := excel.NewCell(excel.CellTypeFloat, rowTotal,
			excel.StyleNone)
		if err != nil {
			return "", err
		}
		colTotal.AppendCell(rCst)
		//sum += mbi.Cost

	}

	totalLabelCell, err := excel.NewCell(excel.CellTypeString, "Total",
		totalCellStyle)
	if err != nil {
		return "", err
	}
	colMonth.AppendCell(totalLabelCell)

	t.AppendCol(colYear)
	t.AppendCol(colMonth)
	for _, v := range subCols {
		t.AppendCol(v)
	}
	t.AppendCol(colTotal)

	tableFormat := `
	{"table_style":"TableStyleMedium4", "show_last_column":true,
	"show_row_stripes":true,"show_column_stripes":false}`

	t.SetTableFormat(tableFormat)
	si := x.GetSheetIndex(sheetName)
	if si == 0 {
		_ = x.NewSheet(sheetName)
	}
	tRange, err := t.ExcelizeWrite(x, sheetName, anchor)
	if err != nil {
		return "", err
	}

	_, brs := excel.R2c(tRange)

	_, totalRn, err := excelize.CellNameToCoordinates(brs)
	tlCn, tlRn, err := excelize.CellNameToCoordinates(anchor)
	if err != nil {
		return tRange, err
	}

	startCn := tlCn + 2
	topRn := tlRn + 1
	botomRn := totalRn - 1

	for i := range subCols {
		cn := startCn + i
		topDataCellName, err := excelize.CoordinatesToCellName(cn, topRn)
		if err != nil {
			return tRange, err
		}

		bottomDataCellName, err := excelize.CoordinatesToCellName(cn, botomRn)
		if err != nil {
			return tRange, err
		}

		sumFormula := fmt.Sprintf("SUM(%s:%s)", topDataCellName, bottomDataCellName)

		totalValueCell, err := excel.NewCell(excel.CellTypeFormula, sumFormula,
			totalCellStyle)
		if err != nil {
			return tRange, err
		}

		totalCellName, err := excelize.CoordinatesToCellName(cn, totalRn)
		if err != nil {
			return tRange, err
		}

		_, err = totalValueCell.ExcelizeWrite(x, sheetName, totalCellName)

	}

	totalColN := startCn + len(subCols)
	topTotalCellName, err := excelize.CoordinatesToCellName(totalColN, topRn)
	if err != nil {
		return tRange, err
	}

	bottomTotalCellName, err := excelize.CoordinatesToCellName(totalColN, botomRn)
	if err != nil {
		return tRange, err
	}

	sumFormula := fmt.Sprintf("SUM(%s:%s)", topTotalCellName, bottomTotalCellName)

	totalValueCell, err := excel.NewCell(excel.CellTypeFormula, sumFormula,
		totalCellStyle)
	if err != nil {
		return tRange, err
	}

	_, err = totalValueCell.ExcelizeWrite(x, sheetName, brs)

	return tRange, err

}

func subscriptionCostItemListJSON2Excel(x *excelize.File, sheetName string,
	anchor string, bis []myJson.SubscriptionCostItem) (string, error) {
	t, err := excel.NewColTable(nil)
	if err != nil {
		return "", err
	}

	totalCellStyle, err := x.
		NewStyle(`{"alignment":{"horizontal":"right"},"font":{"bold":true}}`)
	if err != nil {
		return "", err
	}

	colSubscripton, _ := excel.
		NewColumn("Subscription", nil, excel.StyleNone, excel.StyleNone)
	colCost, _ := excel.
		NewColumn("Cost($)", nil, excel.StyleNone, excel.StyleNone)

	for _, mbi := range bis {
		cSub, err := excel.
			NewCell(excel.CellTypeString, mbi.Suscription, excel.StyleNone)
		if err != nil {
			return "", err
		}
		colSubscripton.AppendCell(cSub)

		cCst, err := excel.
			NewCell(excel.CellTypeFloat, mbi.Cost, excel.StyleNone)
		if err != nil {
			return "", err
		}
		colCost.AppendCell(cCst)
		//sum += mbi.Cost

	}

	totalLabelCell, err := excel.
		NewCell(excel.CellTypeString, "Total", totalCellStyle)
	if err != nil {
		return "", err
	}
	colSubscripton.AppendCell(totalLabelCell)

	t.AppendCol(colSubscripton)
	t.AppendCol(colCost)

	tableFormat := `
	{"table_style":"TableStyleMedium4", "show_last_column":true,
	"show_row_stripes":true,"show_column_stripes":false}`

	t.SetTableFormat(tableFormat)
	si := x.GetSheetIndex(sheetName)
	if si == 0 {
		_ = x.NewSheet(sheetName)
	}
	tRange, err := t.ExcelizeWrite(x, sheetName, anchor)
	if err != nil {
		return "", err
	}

	_, lpos := excel.R2c(tRange)

	lposCn, lposRn, err := excelize.CellNameToCoordinates(lpos)
	if err != nil {
		return tRange, err
	}
	_, anchanchorRn, err := excelize.CellNameToCoordinates(anchor)
	if err != nil {
		return tRange, err
	}

	costColDataTop, err := excelize.CoordinatesToCellName(lposCn,
		anchanchorRn+1)
	if err != nil {
		return tRange, err
	}

	costColDataBottom, err := excelize.CoordinatesToCellName(lposCn, lposRn-1)
	if err != nil {
		return tRange, err
	}

	sumFormula := fmt.Sprintf("SUM(%s:%s)", costColDataTop, costColDataBottom)
	totalValueCell, err := excel.NewCell(excel.CellTypeFormula, sumFormula,
		totalCellStyle)
	if err != nil {
		return tRange, err
	}
	_, err = totalValueCell.ExcelizeWrite(x, sheetName, lpos)
	return tRange, err

}

func customerJSON2Excel(x *excelize.File, sheetName string, anchor string,
	c myJson.Customer) (string, error) {
	t, err := excel.NewColTable(nil)
	if err != nil {
		return "", err
	}
	names := []string{"Customer ID", "Customer Company Name", "Other Name(s)"}
	values := []string{c.CustomerId, c.CustomerCompanyName, c.FormerNames}

	colName, err := excel.NewColumn("", nil, excel.StyleNone, excel.StyleNone)
	if err != nil {
		return "", err
	}
	colValue, err := excel.NewColumn("", nil, excel.StyleNone, excel.StyleNone)
	if err != nil {
		return "", err
	}

	for i := 0; i < 3; i++ {
		n, err := excel.NewCell(excel.CellTypeString, names[i], excel.StyleNone)
		if err != nil {
			return "", err
		}
		colName.AppendCell(n)

		v, err := excel.NewCell(excel.CellTypeString, values[i], excel.StyleNone)
		if err != nil {
			return "", err
		}
		colValue.AppendCell(v)

	}

	t.AppendCol(colName)
	t.AppendCol(colValue)

	return t.ExcelizeWrite(x, sheetName, anchor)
}

func customerListJSON2Excel(x *excelize.File, sheetName string, anchor string,
	cs []myJson.Customer) (string, error) {
	t, err := excel.NewColTable(nil)
	if err != nil {
		return "", err
	}

	colID, err := excel.NewColumn("Customer ID", nil, excel.StyleNone,
		excel.StyleNone)
	if err != nil {
		return "", err
	}
	colName, err := excel.NewColumn("Customer Company Name", nil,
		excel.StyleNone, excel.StyleNone)
	if err != nil {
		return "", err
	}
	colOther, err := excel.NewColumn("Other Name(s)", nil, excel.StyleNone,
		excel.StyleNone)
	if err != nil {
		return "", err
	}

	for _, c := range cs {
		ci, err := excel.NewCell(excel.CellTypeString, c.CustomerId,
			excel.StyleNone)
		if err != nil {
			return "", err
		}
		colID.AppendCell(ci)

		cn, err := excel.NewCell(excel.CellTypeString, c.CustomerCompanyName,
			excel.StyleNone)
		if err != nil {
			return "", err
		}
		colName.AppendCell(cn)

		co, err := excel.NewCell(excel.CellTypeString, c.FormerNames,
			excel.StyleNone)
		if err != nil {
			return "", err
		}
		colOther.AppendCell(co)
	}

	t.AppendCol(colID)
	t.AppendCol(colName)
	t.AppendCol(colOther)
	tableFormat := `
	{"table_style":"TableStyleMedium4", "show_last_column":true,
	"show_row_stripes":true,"show_column_stripes":false}`
	t.SetTableFormat(tableFormat)

	return t.ExcelizeWrite(x, sheetName, anchor)

}

func monthlySummaryCostJSON2Excel(x *excelize.File, sheetName string,
	anchor string, ccs []myJson.CustomerCostItem) (string, error) {
	t, err := excel.NewColTable(nil)
	if err != nil {
		return "", err
	}

	totalCellStyle, err := x.
		NewStyle(`{"alignment":{"horizontal":"right"},"font":{"bold":true}}`)
	if err != nil {
		return "", err
	}

	colName, err := excel.
		NewColumn("Customer Company Name", nil, excel.StyleNone, excel.StyleNone)
	if err != nil {
		return "", err
	}

	colCost, err := excel.
		NewColumn("Cost($)", nil, excel.StyleNone, excel.StyleNone)
	if err != nil {
		return "", err
	}

	for _, cc := range ccs {
		ccn, err := excel.
			NewCell(excel.CellTypeString, cc.Owner.CustomerCompanyName, excel.StyleNone)
		if err != nil {
			return "", err
		}
		colName.AppendCell(ccn)

		ccc, err := excel.NewCell(excel.CellTypeFloat, cc.Cost, excel.StyleNone)
		if err != nil {
			return "", err
		}
		colCost.AppendCell(ccc)
	}
	totalLabelCell, err := excel.
		NewCell(excel.CellTypeString, "Total", totalCellStyle)
	if err != nil {
		return "", err
	}
	totalRowOffset := colName.AppendCell(totalLabelCell)
	t.AppendCol(colName)
	costColOffset := t.AppendCol(colCost)

	totalValueCellName, err := excel.
		OffsetJump(anchor, costColOffset, totalRowOffset)
	if err != nil {
		return "", err
	}
	sumTopCellName, err := excel.OffsetJump(anchor, costColOffset, 1)
	if err != nil {
		return "", err
	}
	sumLastCellName, err := excel.OffsetJump(totalValueCellName, 0, -1)
	if err != nil {
		return "", err
	}

	sumFormula := fmt.Sprintf("SUM(%s:%s)", sumTopCellName, sumLastCellName)

	totalValueCell, err := excel.
		NewCell(excel.CellTypeFormula, sumFormula, totalCellStyle)
	if err != nil {
		return "", err
	}
	err = t.AddCell(totalValueCell, costColOffset, totalRowOffset, false)
	if err != nil {
		return "", err
	}
	tableFormat := `
	{"table_style":"TableStyleMedium4", "show_last_column":true,
	"show_row_stripes":true,"show_column_stripes":false}`
	t.SetTableFormat(tableFormat)
	return t.ExcelizeWrite(x, sheetName, anchor)
}

func initExcel(sheetName string) *excelize.File {
	x := excelize.NewFile()
	sn := x.GetSheetName(1)
	x.SetSheetName(sn, sheetName)
	return x
}

func excelWriteText(x *excelize.File, sheetName string, anchor string,
	value string, h, v int) (string, error) {
	endAxis := anchor
	if v > 1 || h > 1 {
		endAxis, err := excel.OffsetJump(anchor, h-1, v-1)
		if err != nil {
			return "", err
		}
		err = x.MergeCell(sheetName, anchor, endAxis)
	}

	err := x.SetCellStr(sheetName, anchor, value)
	if err != nil {
		return "", err
	}
	return anchor + ":" + endAxis, nil
}

func writeDateRange(x *excelize.File, sheetName string, anchor string,
	dateRange myJson.YearMonthRange) (string, error) {

	_, err := excelWriteText(x, sheetName, anchor, "Start Date", 1, 1)
	if err != nil {
		return "", err
	}
	currAnchor, _ := excel.OffsetJump(anchor, 1, 0)
	_, err = writeDate(x, sheetName, currAnchor, dateRange.StartDate)
	if err != nil {
		return "", err
	}
	currAnchor, _ = excel.OffsetJump(anchor, 0, 1)
	_, err = excelWriteText(x, sheetName, currAnchor, "End Date", 1, 1)
	if err != nil {
		return "", err
	}
	currAnchor, _ = excel.OffsetJump(currAnchor, 1, 0)
	r, err := writeDate(x, sheetName, currAnchor, dateRange.EndDate)
	if err != nil {
		return "", err
	}

	_, e := excel.R2c(r)
	return anchor + ":" + e, nil
}

func writeDate(x *excelize.File, sheetName string, anchor string,
	date myJson.YearMonth) (string, error) {

	_, err := excelWriteText(x, sheetName, anchor,
		fmt.Sprintf("Year: %d", date.Year), 1, 1)
	if err != nil {
		return "", err
	}

	monthCellName, err := excel.OffsetJump(anchor, 1, 0)
	_, err = excelWriteText(x, sheetName, monthCellName,
		fmt.Sprintf("Month: %s", time.Month(date.Month).String()), 1, 1)
	if err != nil {
		return "", err
	}
	return anchor + ":" + monthCellName, nil
}

func excel2Bytes(x *excelize.File) ([]byte, error) {
	b, err := x.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

type chartlet struct { //sheetName+"!"+nameLabel,sheetName+"!"+categoryRange,sheetName+"!"+valueRange
	name     string
	catRange string
	valRange string
}

func (c chartlet) doString() string {
	s := fmt.Sprintf(chartValTpl2, c.name, c.catRange, c.valRange)
	return s

}

type chart struct { //cType,width,height,title
	title  string
	cType  string
	width  int
	height int
	lines  []chartlet
}

func (c chart) doString() string {
	head := fmt.Sprintf(chartValTpl1, c.cType, c.width, c.height, c.title)
	tail := `]}`
	mids := make([]string, len(c.lines))
	for i, v := range c.lines {
		mids[i] = fmt.Sprintf(v.doString())
	}
	mid := strings.Join(mids, ",")
	s := head + mid + tail
	return s
}

const (
	chartValTpl1 = `{"type":"%s","dimension":{"width":%d,"height":%d},"legend":{"position":"bottom","show_legend_key":false},"title":{"name":"%s"},"series":[`
	chartValTpl2 = `{"name":"%s","categories":"%s","values":"%s"}`
	tblFmt       = `
	{"table_style":"TableStyleMedium4", "show_last_column":true,
	"show_row_stripes":true,"show_column_stripes":false}`
)
