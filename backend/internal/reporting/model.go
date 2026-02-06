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
