package transactiondomain

const TRANSACTION_TYPE_BUY = "BUY"
const TRANSACTION_TYPE_SELL = "SELL"
const TRANSACTION_PEMBUATAN = "PEMBUATAN"
const TRANSACTION_DICETAK = "DICETAK"
const TRANSACTION_KONTRABON = "KONTRABON"
const TRANSACTION_DIBAYAR = "DIBAYAR"
const TRANSACTION_BATAL = "BATAL"

type Transaction struct {
	ID                string               `json:"id"`
	Code              string               `json:"code"`
	Date              string               `json:"date"`
	StakeholderID     string               `json:"stakeholderId"`
	TransactionType   string               `json:"transactionType"`
	Status            string               `json:"status"`
	ReferenceCode     string               `json:"referenceCode"`
	UserId            string               `json:"userId"`
	CreatedTime       string               `json:"createdTime"`
	Total             float64                `json:"total"`
	TransactionDetail []*TransactionDetail `json:"transactionDetail"`
}

type TransactionDetail struct {
	ID            string  `json:"id"`
	TransactionID string  `json:"transactionId"`
	UnitID        string  `json:"unitId"`
	ProductID     string  `json:"productId"`
	BuyPrice      float64 `json:"buyPrice"`
	SellPrice     float64 `json:"sellPrice"`
	Quantity      float64   `json:"quantity"`
	BuyQuantity   float64   `json:"buyQuantity"`
}

type TransactionStatus struct {
	ID                string                     `json:"id"`
	Code              string                     `json:"code"`
	Date              string                     `json:"date"`
	StakeholderID     string                     `json:"stakeholderId"`
	StakeholderCode   string                     `json:"stakeholderCode"`
	StakeholderName   string                     `json:"stakeholderName"`
	TransactionType   string                     `json:"transactionType"`
	Status            string                     `json:"status"`
	ReferenceCode     string                     `json:"referenceCode"`
	UserId            string                     `json:"userId"`
	UserName          string                     `json:"userName"`
	CreatedTime       string                     `json:"createdTime"`
	Total             float64                      `json:"total"`
	TransactionDetail []*TransactionStatusDetail `json:"transactionDetail"`
}

type TransactionStatusDetail struct {
	TransactionID string  `json:"transactionId"`
	UnitID        string  `json:"unitId"`
	UnitCode      string  `json:"unitCode"`
	ProductID     string  `json:"productId"`
	ProductCode   string  `json:"productCode"`
	ProductName   string  `json:"productName"`
	BuyPrice      float64 `json:"buyPrice"`
	SellPrice     float64 `json:"sellPrice"`
	Quantity      float64   `json:"quantity"`
	BuyQuantity   float64   `json:"buyQuantity"`
}
