package kontrabondomain

const STATUS_CREATED = "CREATED"
const STATUS_LUNAS = "LUNAS"

type Kontrabon struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	CreatedTime string `json:"createdTime"`
	Status      string `json:"status"`
	Total       int64  `json:"total"`
	CustomerID  string `json:"customerId"`
}

type CreateRequest struct {
	CustomerID     string   `json:"customerId"`
	TransactionIDs []string `json:"transactionIds"`
}

type UpdateRequest struct {
	KontrabonID    string   `json:"kontrabonId"`
	TransactionIDs []string `json:"transactionIds"`
}