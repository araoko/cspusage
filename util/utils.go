package util

import (
	"database/sql"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/araoko/cspusage/entity"
	myJson "github.com/araoko/cspusage/model/json"
	"github.com/jmoiron/sqlx"
)

type monthRangePerYr struct {
	yr int
	sm int
	em int
}

func DoCustomerList(db *sqlx.DB) ([]myJson.Customer, error) {
	list, err := entity.GetCustomers(db)
	if err != nil {
		return nil, err
	}

	customers := make([]myJson.Customer, len(list))

	for i, item := range list {
		customers[i] = myJson.Customer{
			CustomerId:          item.CustomerID,
			CustomerCompanyName: item.CustomerCompanyName,
			FormerNames:         item.FormerNames,
		}

	}
	return customers, nil
}

func DoMonthlyBill(db *sqlx.DB, input myJson.YearMonth) (myJson.MonthlyBill, error) {
	mb := myJson.MonthlyBill{Date: input,
		CustomerMonthlyBills: make([]myJson.CustomerBillNoDate, 0),
	}

	customerMap, err := GetCustomerMap(db)
	if err != nil {
		return mb, err
	}
	//var total float32
	summary := myJson.MonthlyCostSummary{
		Date:      input,
		Customers: make([]myJson.CustomerCostItem, 0),
	}
	for id, cust := range customerMap {

		mbits, err := entity.GetCustomerMonthlyBill(db, id, input.Year, input.Month)
		if err != nil {
			return mb, err
		}
		if len(mbits) == 0 {
			continue
		}
		cmb := myJson.CustomerBillNoDate{
			Owner:     cust,
			LineItems: make([]myJson.SubscriptionServiceCostItem, len(mbits)),
		}
		var totalPerCustomer float32
		customerCost := myJson.CustomerCostItem{
			Owner: cust,
		}
		for i, mbit := range mbits {
			cmb.LineItems[i].Suscription = mbit.Subscription
			cmb.LineItems[i].ServiceNameAndType = mbit.Service
			cmb.LineItems[i].Cost = mbit.Cost
			totalPerCustomer += mbit.Cost
		}
		customerCost.Cost = totalPerCustomer
		//total += totalPerCustomer
		mb.CustomerMonthlyBills = append(mb.CustomerMonthlyBills, cmb)
		summary.Customers = append(summary.Customers, customerCost)

	}
	mb.Summary = summary
	return mb, nil
}

func DoCustomerMonthlyBill(db *sqlx.DB, input myJson.CustomerIDAndDate) (myJson.CustomerMonthlyBill, error) {
	result := myJson.CustomerMonthlyBill{Date: input.Date}

	custEntity, err := entity.GetCustomerFromID(db, input.CustomerId)
	if err == sql.ErrNoRows {

		return result, err
	}
	if err != nil {
		return result, err
	}

	custJson := myJson.Customer{
		CustomerId:          custEntity.CustomerID,
		CustomerCompanyName: custEntity.CustomerCompanyName,
		FormerNames:         custEntity.FormerNames,
	}

	result.Owner = custJson
	ss, err := entity.GetCustomerMonthlyBill(db, input.CustomerId, input.Date.Year, input.Date.Month)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return result, err
	}
	//log.Println("got lineitems. Count: ", len(ss))
	lineItems := make([]myJson.SubscriptionServiceCostItem, len(ss))
	for i, l := range ss {
		lineItems[i] = myJson.SubscriptionServiceCostItem{
			Suscription:        l.Subscription,
			ServiceNameAndType: l.Service,
			Cost:               l.Cost,
		}
	}
	result.LineItems = lineItems
	return result, nil
}

//doCustomerMonthlyBillPerSub queries the database for usage of a an account in a month
//summing the cost per subscription
func DoCustomerMonthlyCostPerSub(db *sqlx.DB, input myJson.CustomerIDAndDate) (myJson.CustomerMonthlyCostPerSub, error) {
	result := myJson.CustomerMonthlyCostPerSub{Date: input.Date}

	custEntity, err := entity.GetCustomerFromID(db, input.CustomerId)
	if err == sql.ErrNoRows {

		return result, err
	}
	if err != nil {
		return result, err
	}

	custJson := myJson.Customer{
		CustomerId:          custEntity.CustomerID,
		CustomerCompanyName: custEntity.CustomerCompanyName,
		FormerNames:         custEntity.FormerNames,
	}

	result.Owner = custJson
	ss, err := entity.GetCustomerMonthlyCostPerSub(db, input.CustomerId, input.Date.Year, input.Date.Month)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return result, err
	}
	//log.Println("got lineitems. Count: ", len(ss))
	lineItems := make([]myJson.SubscriptionCostItem, len(ss))
	for i, l := range ss {
		lineItems[i] = myJson.SubscriptionCostItem{
			Suscription: l.Subscription,
			Cost:        l.Cost,
		}
	}
	result.LineItems = lineItems
	return result, nil
}

func DoMonthlySummary(db *sqlx.DB, input myJson.YearMonth) (myJson.MonthlyCostSummary, error) {
	res := myJson.MonthlyCostSummary{
		Date: input,
	}
	customerMap, err := GetCustomerMap(db)
	if err != nil {
		return res, err
	}

	summaryLineItems, err := entity.GetMonthlyCostSummary(db, input.Year, input.Month)
	if err != nil {
		return res, err
	}
	res.Customers = make([]myJson.CustomerCostItem, len(summaryLineItems))

	for i, s := range summaryLineItems {
		cust := customerMap[s.CustomerID]
		res.Customers[i].Owner = cust
		res.Customers[i].Cost = s.Cost

	}

	return res, nil

}

func DoRangeBill(db *sqlx.DB, input myJson.YearMonthRange) (myJson.RangeBill, error) {
	rb := myJson.RangeBill{
		DateRange:          input,
		CustomerRangeBills: make([]myJson.CustomerBillNoDate, 0),
	}

	customerMap, err := GetCustomerMap(db)
	if err != nil {
		return rb, err
	}

	summary := myJson.RangeCostSummary{
		DateRange: input,
		Customers: make([]myJson.CustomerCostItem, 0),
	}

	for id := range customerMap {
		cidadr := myJson.CustomerIDAndDateRange{
			CustomerId: id,
			DateRange:  input,
		}
		crb, err := DoCustomerRangeBill(db, cidadr)
		if err != nil {
			return rb, err
		}
		var summaryCost float32
		for _, li := range crb.LineItems {
			summaryCost += li.Cost
		}
		cc := myJson.CustomerCostItem{
			Owner: crb.Owner,
			Cost:  summaryCost,
		}

		cbnd := myJson.CustomerBillNoDate{
			Owner:     crb.Owner,
			LineItems: crb.LineItems,
		}
		summary.Customers = append(summary.Customers, cc)
		rb.CustomerRangeBills = append(rb.CustomerRangeBills, cbnd)
	}
	rb.Summary = summary
	return rb, nil

}

func DoCustomerRangeCostPerSub(db *sqlx.DB, input myJson.CustomerIDAndDateRange) (myJson.CustomerRangeCostPerSub, error) {
	resp := myJson.CustomerRangeCostPerSub{
		DateRange: input.DateRange,
	}

	custEntity, err := entity.GetCustomerFromID(db, input.CustomerId)
	if err == sql.ErrNoRows {
		return resp, fmt.Errorf("Customer Id: %s does not exist in database", input.CustomerId)
	}
	if err != nil {
		return resp, err
	}

	custJson := myJson.Customer{
		CustomerId:          custEntity.CustomerID,
		CustomerCompanyName: custEntity.CustomerCompanyName,
		FormerNames:         custEntity.FormerNames,
	}

	resp.Owner = custJson

	dateRanges := GetRange(input.DateRange.StartDate.Year, input.DateRange.StartDate.Month, input.DateRange.EndDate.Year, input.DateRange.EndDate.Month)
	subMap := make(map[string]float32)

	for _, dateRange := range dateRanges {
		ss, err := entity.GetCustomerRangeCostPerSub(db, input.CustomerId, dateRange.yr, dateRange.sm, dateRange.em)
		if err != nil {
			return resp, err
		}
		for _, l := range ss {
			subMap[l.Subscription] += l.Cost

		}

	}

	lineItems := make([]myJson.SubscriptionCostItem, len(subMap))
	subCount := 0
	for k, v := range subMap {
		lineItem := myJson.SubscriptionCostItem{
			Suscription: k,
			Cost:        v,
		}
		lineItems[subCount] = lineItem
		subCount++
	}
	lineItems2 := make([]myJson.SubscriptionCostItem, 0)

	//had to do these loops to dedup subscription
OUTER:
	for _, li1 := range lineItems {
		for j, li2 := range lineItems2 {
			if strings.EqualFold(li1.Suscription, li2.Suscription) {
				lineItems2[j].Cost += li1.Cost
				continue OUTER
			}
		}
		lineItems2 = append(lineItems2, li1)
	}
	resp.LineItems = lineItems2

	return resp, nil
}

func DoCustomerRangeBill(db *sqlx.DB, input myJson.CustomerIDAndDateRange) (myJson.CustomerRangeBill, error) {
	resp := myJson.CustomerRangeBill{
		DateRange: input.DateRange,
	}

	custEntity, err := entity.GetCustomerFromID(db, input.CustomerId)
	if err == sql.ErrNoRows {
		return resp, fmt.Errorf("Customer Id: %s does not exist in database", input.CustomerId)
	}
	if err != nil {
		return resp, err
	}

	custJson := myJson.Customer{
		CustomerId:          custEntity.CustomerID,
		CustomerCompanyName: custEntity.CustomerCompanyName,
		FormerNames:         custEntity.FormerNames,
	}

	resp.Owner = custJson

	dateRanges := GetRange(input.DateRange.StartDate.Year, input.DateRange.StartDate.Month, input.DateRange.EndDate.Year, input.DateRange.EndDate.Month)
	subMap := make(map[string]map[string]float32)

	for _, dateRange := range dateRanges {
		ss, err := entity.GetCustomerRangeBill(db, input.CustomerId, dateRange.yr, dateRange.sm, dateRange.em)
		if err != nil {
			return resp, err
		}
		for _, l := range ss {
			lineItem := myJson.SubscriptionServiceCostItem{
				Suscription:        l.Subscription,
				ServiceNameAndType: l.Service,
				Cost:               l.Cost,
			}
			AddLineItem2Map(subMap, lineItem)

		}

	}

	resp.LineItems = map2LineItems(subMap)

	return resp, nil
}

func DoCustomerPerSubTrend(db *sqlx.DB, input myJson.CustomerIDAndDateRange) (myJson.CustomerMonthlyPerSubTrend, error) {
	//TODO In porgress

	resp := myJson.CustomerMonthlyPerSubTrend{
		DateRange: input.DateRange,
	}

	custEntity, err := entity.GetCustomerFromID(db, input.CustomerId)
	if err == sql.ErrNoRows {
		return resp, fmt.Errorf("Customer Id: %s does not exist in database", input.CustomerId)
	}
	if err != nil {
		return resp, err
	}

	custJson := myJson.Customer{
		CustomerId:          custEntity.CustomerID,
		CustomerCompanyName: custEntity.CustomerCompanyName,
		FormerNames:         custEntity.FormerNames,
	}

	resp.Owner = custJson

	dateRanges := GetRange(input.DateRange.StartDate.Year,
		input.DateRange.StartDate.Month,
		input.DateRange.EndDate.Year,
		input.DateRange.EndDate.Month)

	subCol := make(map[string]struct{})
	yearArr := make([]map[int][]myJson.CostSub, 0)
	for _, dateRange := range dateRanges {
		yMC, err := entity.GetCustomerPerSubMonthlyTrend(db, custJson.CustomerId, dateRange.yr, dateRange.sm, dateRange.em)
		if err != nil {
			return resp, err
		}
		monthSubCostHash := make(map[int][]myJson.CostSub)
		for _, yMCLI := range yMC {
			subCol[yMCLI.Subscription] = struct{}{}
			cs := myJson.CostSub{
				Subscription: yMCLI.Subscription,
				Cost:         yMCLI.Cost,
			}
			m := yMCLI.Month
			if csa, exist := monthSubCostHash[m]; exist {
				csa = append(csa, cs)
				monthSubCostHash[m] = csa
			} else {
				monthSubCostHash[m] = []myJson.CostSub{cs}
			}

		}
		yearArr = append(yearArr, monthSubCostHash)

	}
	arrSub := make([]string, 0)

	for sub := range subCol {
		arrSub = append(arrSub, sub)
	}
	dcsArr := make([]myJson.DateCostSubItem, 0)
	for i := 0; i < len(yearArr); i++ {
		yr := dateRanges[i].yr
		monthSubCostHash := yearArr[i]

		for m, sca := range monthSubCostHash {
			ym := myJson.YearMonth{
				Year:  yr,
				Month: m,
			}
			cps := GetCostSubs(arrSub, sca)
			dcs := myJson.DateCostSubItem{
				Date:        ym,
				CostPerSubs: cps,
			}
			dcsArr = append(dcsArr, dcs)
		}

	}

	resp.Trend = dcsArr

	return resp, nil
}

func GetCostSubs(subList []string, data []myJson.CostSub) []myJson.CostSub {
	result := make([]myJson.CostSub, len(subList))

OUTERLOOP:
	for i := 0; i < len(subList); i++ {
		sub := subList[i]
		//fmt.Println("Checking match for:", sub)
		for _, cs := range data {
			if strings.EqualFold(sub, cs.Subscription) {
				//fmt.Println("Mathched:", cs)
				result[i] = cs
				continue OUTERLOOP
			}
		}
		result[i] = myJson.CostSub{
			Subscription: sub,
			Cost:         0.0,
		}
	}
	return result
}

func DoCustomerTrend(db *sqlx.DB, input myJson.CustomerIDAndDateRange) (myJson.CustomerMonthlyTrend, error) {
	//TODO In porgress

	resp := myJson.CustomerMonthlyTrend{
		DateRange: input.DateRange,
	}

	custEntity, err := entity.GetCustomerFromID(db, input.CustomerId)
	if err == sql.ErrNoRows {
		return resp, fmt.Errorf("Customer Id: %s does not exist in database", input.CustomerId)
	}
	if err != nil {
		return resp, err
	}

	custJson := myJson.Customer{
		CustomerId:          custEntity.CustomerID,
		CustomerCompanyName: custEntity.CustomerCompanyName,
		FormerNames:         custEntity.FormerNames,
	}

	resp.Owner = custJson

	lis := make([]myJson.DateCostItem, 0)

	dateRanges := GetRange(input.DateRange.StartDate.Year,
		input.DateRange.StartDate.Month,
		input.DateRange.EndDate.Year,
		input.DateRange.EndDate.Month)

	// cMap, err := getCustomerMap(db)
	// if err != nil {
	// 	return resp, err
	// }

	//load lineitems into map
	for _, dateRange := range dateRanges {
		yMC, err := entity.GetCustomerMonthlyTrend(db, custJson.CustomerId, dateRange.yr, dateRange.sm, dateRange.em)
		if err != nil {
			return resp, err
		}
		for _, yMCLI := range yMC {
			li := myJson.DateCostItem{
				Date: myJson.YearMonth{
					Year:  yMCLI.Year,
					Month: yMCLI.Month,
				},
				Cost: yMCLI.Cost,
			}

			lis = append(lis, li)

		}

	}

	resp.Trend = lis

	return resp, nil
}

func DoTrend(db *sqlx.DB, input myJson.YearMonthRange) (myJson.MonthlyTrend, error) {
	resp := myJson.MonthlyTrend{
		DateRange: input,
		Summary:   make([]myJson.DateCostItem, 0),
		Trend:     make([]myJson.CustomerMonthlyTrendNoDateRange, 0),
	}
	gMap := make(map[string]map[int]map[int]float32)
	//cIDMap := make(map[string]*myJson.CustomerMonthlyTrendNoDateRange)
	//sMap := make(map[string]*myJson.CustomerCostItem)
	sMap := make(map[int]map[int]float32)
	mrprs := GetRange(input.StartDate.Year, input.StartDate.Month, input.EndDate.Year, input.EndDate.Month)
	cMap, err := GetCustomerMap(db)
	if err != nil {
		return resp, err
	}

	//load lineitems into map
	for _, mrpr := range mrprs {
		idYMC, err := entity.GetMonthlyTrend(db, mrpr.yr, mrpr.sm, mrpr.em)
		if err != nil {
			return resp, err
		}
		yMap := make(map[int]float32)
		for _, idYMCLI := range idYMC {
			cid := idYMCLI.CustomerID
			yr := idYMCLI.Year
			mo := idYMCLI.Month
			c := idYMCLI.Cost
			idMap, exist := gMap[cid]
			if !exist {
				idMap = make(map[int]map[int]float32)
				gMap[cid] = idMap
			}

			yrMap, exist := idMap[yr]
			if !exist {
				yrMap = make(map[int]float32)
				idMap[yr] = yrMap
			}
			yrMap[mo] = c
			yMap[mo] += c
		}
		sMap[mrpr.yr] = yMap
	}

	//load map to CustomerMonthlyTrendNoDateRange array
	for cid, idMap := range gMap {
		l1 := myJson.CustomerMonthlyTrendNoDateRange{
			Owner: cMap[cid],
			Trend: make([]myJson.DateCostItem, 0),
		}

		for _, mrpr := range mrprs {
			yr := mrpr.yr
			for m := mrpr.sm; m <= mrpr.em; m++ {
				l2 := myJson.DateCostItem{
					Date: myJson.YearMonth{
						Year:  yr,
						Month: m,
					},
				}
				var c float32
				yrMap, exist := idMap[yr]
				if exist {
					c = yrMap[m]
				}
				l2.Cost = c
				l1.Trend = append(l1.Trend, l2)
			}

		}
		resp.Trend = append(resp.Trend, l1)
	}

	for _, mrpr := range mrprs {
		yr := mrpr.yr
		for m := mrpr.sm; m <= mrpr.em; m++ {
			ci := myJson.DateCostItem{
				Date: myJson.YearMonth{
					Year:  yr,
					Month: m,
				},
			}
			var c float32
			yMap, exist := sMap[yr]
			if exist {
				c = yMap[m]
			}
			ci.Cost = c

			resp.Summary = append(resp.Summary, ci)
		}
	}

	return resp, nil
}

func Customer2CustomerMap(cs []entity.Customer) map[string]myJson.Customer {
	customerMap := make(map[string]myJson.Customer)
	for _, c := range cs {
		customerMap[c.CustomerID] = myJson.Customer{
			CustomerId:          c.CustomerID,
			CustomerCompanyName: c.CustomerCompanyName,
			FormerNames:         c.FormerNames,
		}
	}
	return customerMap
}

func GetCustomerMap(db *sqlx.DB) (map[string]myJson.Customer, error) {
	customers, err := entity.GetCustomers(db)
	if err != nil {
		return nil, err
	}

	return Customer2CustomerMap(customers), nil
}

func SortCustomerPerSubTrend(t myJson.CustomerMonthlyPerSubTrend) myJson.CustomerMonthlyPerSubTrend {
	csa := t.Trend[0].CostPerSubs
	csa0 := make([]myJson.CostSub, len(csa))
	for k, v := range csa {
		csa0[k] = myJson.CostSub{
			Subscription: v.Subscription,
			Cost:         0.0,
		}
	}

	ma := make(map[myJson.YearMonth][]myJson.CostSub)
	for _, v := range t.Trend {
		ma[v.Date] = v.CostPerSubs
	}

	st := make([]myJson.DateCostSubItem, 0, len(t.Trend))

	yma := ExpandYearMonthRange2YearMonths(t.DateRange)
	for _, ym := range yma {
		dc := myJson.DateCostSubItem{
			Date: ym,
		}
		if k, exist := ma[ym]; exist {
			dc.CostPerSubs = k
		} else {
			dc.CostPerSubs = csa0
		}
		st = append(st, dc)
	}

	return myJson.CustomerMonthlyPerSubTrend{
		Owner:     t.Owner,
		DateRange: t.DateRange,
		Trend:     st,
	}

}

func SortCustomerTrend(t myJson.CustomerMonthlyTrend) myJson.CustomerMonthlyTrend {

	ma := make(map[myJson.YearMonth]float32)
	for _, v := range t.Trend {
		ma[v.Date] = v.Cost
	}

	st := make([]myJson.DateCostItem, 0, len(t.Trend))

	yma := ExpandYearMonthRange2YearMonths(t.DateRange)
	for _, ym := range yma {
		dc := myJson.DateCostItem{
			Date: ym,
		}
		if k, exist := ma[ym]; exist {
			dc.Cost = k
		} else {
			dc.Cost = 0.0
		}
		st = append(st, dc)
	}

	return myJson.CustomerMonthlyTrend{
		Owner:     t.Owner,
		DateRange: t.DateRange,
		Trend:     st,
	}

}

func ExpandYearMonthRange2YearMonths(dr myJson.YearMonthRange) []myJson.YearMonth {
	sy := dr.StartDate.Year
	sm := dr.StartDate.Month

	ey := dr.EndDate.Year
	em := dr.EndDate.Month

	y := sy
	m := sm
	yma := make([]myJson.YearMonth, 0)
	for {

		ym := myJson.YearMonth{
			Year:  y,
			Month: m,
		}

		yma = append(yma, ym)
		if y == ey && m == em {
			break
		}

		m++

		if m == 13 {
			m = 1
			y++
		}
	}
	return yma

}

func GetRange(sy, sm, ey, em int) []monthRangePerYr {

	d := make(map[int]monthRangePerYr)
	m := sm
	for y := sy; y <= ey; y++ {
		e := monthRangePerYr{y, m, m}
		d[y] = e
		m++
		for {
			if m == 13 {
				m = 1
				break
			}
			if y == ey && m > em {
				break
			}
			e.em = m
			d[y] = e
			m++
		}

	}
	res := make([]monthRangePerYr, len(d))

	for y := sy; y <= ey; y++ {
		res[y-sy] = d[y]

	}
	return res
}

func AddLineItem2Map(subMap map[string]map[string]float32, l myJson.SubscriptionServiceCostItem) {
	sub := l.Suscription
	svc := l.ServiceNameAndType
	cost := l.Cost

	svcMap, ok := subMap[sub]
	if !ok {
		subMap[sub] = map[string]float32{svc: cost}
		return
	}
	svcMap[svc] += cost
}

func map2LineItems(subMap map[string]map[string]float32) []myJson.SubscriptionServiceCostItem {
	lineItems := []myJson.SubscriptionServiceCostItem{}
	for sub, svcmap := range subMap {
		for svc, cost := range svcmap {
			lineItem := myJson.SubscriptionServiceCostItem{
				Suscription:        sub,
				ServiceNameAndType: svc,
				Cost:               cost,
			}
			lineItems = append(lineItems, lineItem)
		}
	}
	return lineItems
}

func ParseIPNet(ipnetString string) (net.IPNet, error) {
	ipnet := net.IPNet{}
	ipmask := strings.Split(ipnetString, "/")
	if len(ipmask) != 2 {
		return ipnet, fmt.Errorf("%s is not valid cdir  string format", ipnetString)
	}
	ip := net.ParseIP(ipmask[0])
	if ip == nil {
		return ipnet, fmt.Errorf("IP part of cdir string (%s) is not valid ip string format", ipmask[0])
	}
	ipnet.IP = ip
	subnet := strings.Split(ipmask[1], ".")
	if len(subnet) == 4 {
		sbn := make([]byte, 4)
		for i, v := range subnet {
			in, err := strconv.Atoi(v)
			if err != nil {
				return ipnet, fmt.Errorf("invalid subnet mask %s", ipmask[1])
			}
			sbn[i] = byte(in)
		}
		mask := net.IPv4Mask(sbn[0], sbn[1], sbn[2], sbn[3])
		ipnet.Mask = mask
		return ipnet, nil

	}
	if len(subnet) == 1 {
		netMask, err := strconv.Atoi(ipmask[1])
		if err != nil {
			return ipnet, err
		}

		if netMask < 0 || netMask > 32 {
			return ipnet, fmt.Errorf("invalid subnet number %d", netMask)
		}

		mask := net.CIDRMask(netMask, 32)

		ipnet.Mask = mask

		return ipnet, nil

	}

	return ipnet, fmt.Errorf("Invalid sumbet mask %s", ipmask[1])

}
