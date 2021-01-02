package json

//Customer ...
type Customer struct {
	CustomerId          string `json:"customer_id"`
	CustomerCompanyName string `json:"customer_company_name"`
	FormerNames         string `json:"other_names"`
}

//SubscriptionServiceCostItem ...
type SubscriptionServiceCostItem struct {
	Suscription        string  `json:"subscription"`
	ServiceNameAndType string  `json:"service_name_type"`
	Cost               float32 `json:"cost"`
}

type SubscriptionCostItem struct {
	Suscription string  `json:"subscription"`
	Cost        float32 `json:"cost"`
}

//YearMonth ...
type YearMonth struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}

type YearMonthRange struct {
	StartDate YearMonth `json:"start_Date"`
	EndDate   YearMonth `json:"end_Date"`
}

//CustomerMonthlyBill ...
type CustomerMonthlyBill struct {
	Date      YearMonth                     `json:"date"`
	Owner     Customer                      `json:"customer"`
	LineItems []SubscriptionServiceCostItem `json:"line_items"`
}

type CustomerMonthlyCostPerSub struct {
	Date      YearMonth              `json:"date"`
	Owner     Customer               `json:"customer"`
	LineItems []SubscriptionCostItem `json:"line_items"`
}

type CustomerRangeCostPerSub struct {
	DateRange YearMonthRange         `json:"date_range"`
	Owner     Customer               `json:"customer"`
	LineItems []SubscriptionCostItem `json:"line_items"`
}

type CustomerRangeBill struct {
	DateRange YearMonthRange                `json:"date_range"`
	Owner     Customer                      `json:"customer"`
	LineItems []SubscriptionServiceCostItem `json:"line_items"`
}

//CustomerBillNoDate ...
type CustomerBillNoDate struct {
	Owner     Customer                      `json:"customer"`
	LineItems []SubscriptionServiceCostItem `json:"line_items"`
}

//MonthlyBill ...
type MonthlyBill struct {
	Date                 YearMonth            `json:"date"`
	Summary              MonthlyCostSummary   `json:"summary"`
	CustomerMonthlyBills []CustomerBillNoDate `json:"customer_monthly_bills"`
}

type RangeBill struct {
	DateRange          YearMonthRange       `json:"date_range"`
	Summary            RangeCostSummary     `json:"summary"`
	CustomerRangeBills []CustomerBillNoDate `json:"customer_monthly_bills"`
}

//CustomerIDAndDate ...
type CustomerIDAndDate struct {
	CustomerId string    `json:"customer_id"`
	Date       YearMonth `json:"date"`
}

//CustomerIDAndDate ...
type CustomerIDAndDateRange struct {
	CustomerId string         `json:"customer_id"`
	DateRange  YearMonthRange `json:"date_range"`
}

type DateCostItem struct {
	Date YearMonth `json:"date"`
	Cost float32   `json:"cost"`
}

type CustomerMonthlyTrend struct {
	Owner     Customer       `json:"customer"`
	DateRange YearMonthRange `json:"date_range"`
	Trend     []DateCostItem `json:"trend"`
}

type CostSub struct {
	Subscription string  `json:"subscription"`
	Cost         float32 `json:"cost"`
}

type DateCostSubItem struct {
	Date        YearMonth `json:"date"`
	CostPerSubs []CostSub `json:"cost_subs"`
}
type CustomerMonthlyPerSubTrend struct {
	Owner     Customer          `json:"customer"`
	DateRange YearMonthRange    `json:"date_range"`
	Trend     []DateCostSubItem `json:"trend"`
}

type CustomerMonthlyTrendNoDateRange struct {
	Owner Customer       `json:"customer"`
	Trend []DateCostItem `json:"trend"`
}

type MonthlyTrend struct {
	DateRange YearMonthRange                    `json:"date_range"`
	Summary   []DateCostItem                    `json:"summary"`
	Trend     []CustomerMonthlyTrendNoDateRange `json:"trend"`
}

type CustomerCostItem struct {
	Owner Customer `json:"customer"`
	Cost  float32  `json:"cost"`
}

type MonthlyCostSummary struct {
	Date      YearMonth          `json:"date"`
	Customers []CustomerCostItem `json:"customers"`
}

type RangeCostSummary struct {
	DateRange YearMonthRange     `json:"date_range"`
	Customers []CustomerCostItem `json:"customers"`
}

type Login struct {
	UserName string `json:"user_name"`
	Password string `json:"user_password"`
	LoggedIn int    `json:"logged_in"`
}
