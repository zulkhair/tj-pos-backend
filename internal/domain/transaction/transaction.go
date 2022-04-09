package transactiondomain

const TRANSACTION_TYPE_BUY = "BUY"
const TRANSACTION_TYPE_SELL = "SELL"
const TRANSACTION_STATUS_PEMBUATAN = "PEMBUATAN"
const TRANSACTION_STATUS_BELUM_LUNAS = "BELUM_LUNAS"
const TRANSACTION_STATUS_SUDAH_DITAGIH = "SUDAH_DITAGIH"
const TRANSACTION_STATUS_LUNAS = "LUNAS"
const TRANSACTION_STATUS_BATAL = "BATAL"

type Transaction struct {
	ID                string               `json:"id"`
	Code              string               `json:"code"`
	Date              string               `json:"date"`
	StakeholderID     string               `json:"stakeholderId"`
	TransactionType   string               `json:"transactionType"`
	Status            string               `json:"status"`
	ReferenceCode     string               `json:"referenceCode"`
	TransactionDetail []*TransactionDetail `json:"transactionDetail"`
}

type TransactionDetail struct {
	TransactionID string  `json:"transactionId"`
	UnitID        string  `json:"unitId"`
	ProductID     string  `json:"productId"`
	Price         float64 `json:"price"`
	Quantity      int64   `json:"quantity"`
}
