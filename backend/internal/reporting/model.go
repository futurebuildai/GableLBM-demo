package reporting

type DailyTillReport struct {
	Date             string             `json:"date"`
	TotalCollected   float64            `json:"total_collected"`
	ByMethod         map[string]float64 `json:"by_method"`
	TransactionCount int                `json:"transaction_count"`
}

type SalesSummaryReport struct {
	StartDate      string  `json:"start_date"`
	EndDate        string  `json:"end_date"`
	TotalInvoiced  float64 `json:"total_invoiced"`
	TotalCollected float64 `json:"total_collected"`
	OutstandingAR  float64 `json:"outstanding_ar"`
	InvoiceCount   int     `json:"invoice_count"`
}

// AR Aging Report
type ARAgingBucket struct {
	CustomerID   string  `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	Current      float64 `json:"current"`   // 0-30 days
	Days31to60   float64 `json:"days_31_60"`
	Days61to90   float64 `json:"days_61_90"`
	Over90       float64 `json:"over_90"`
	Total        float64 `json:"total"`
}

type ARAgingReport struct {
	AsOfDate      string          `json:"as_of_date"`
	Buckets       []ARAgingBucket `json:"buckets"`
	TotalCurrent  float64         `json:"total_current"`
	Total31to60   float64         `json:"total_31_60"`
	Total61to90   float64         `json:"total_61_90"`
	TotalOver90   float64         `json:"total_over_90"`
	GrandTotal    float64         `json:"grand_total"`
}

// Customer Statement
type StatementLine struct {
	Date        string  `json:"date"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
	Balance     float64 `json:"balance"`
}

type CustomerStatement struct {
	CustomerID   string          `json:"customer_id"`
	CustomerName string          `json:"customer_name"`
	StartDate    string          `json:"start_date"`
	EndDate      string          `json:"end_date"`
	OpenBalance  float64         `json:"open_balance"`
	CloseBalance float64         `json:"close_balance"`
	Lines        []StatementLine `json:"lines"`
}
