package transactiondomain

import (
	"time"
)

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
	CreatedTime       time.Time            `json:"createdTime"`
	Total             float64              `json:"total"`
	TransactionDetail []*TransactionDetail `json:"transactionDetail"`
}

type TransactionDetail struct {
	ID            string    `json:"id"`
	TransactionID string    `json:"transactionId"`
	UnitID        string    `json:"unitId"`
	ProductID     string    `json:"productId"`
	BuyPrice      float64   `json:"buyPrice"`
	SellPrice     float64   `json:"sellPrice"`
	Quantity      float64   `json:"quantity"`
	BuyQuantity   float64   `json:"buyQuantity"`
	CreatedTime   time.Time `json:"-"`
	WebUserID     string    `json:"-"`
	Latest        bool      `json:"-"`
	SortingVal    int64     `json:"-"`
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
	Total             float64                    `json:"total"`
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
	Quantity      float64 `json:"quantity"`
	BuyQuantity   float64 `json:"buyQuantity"`
}

type ReportDate struct {
	Date    string    `json:"date"`
	Reports []*Report `json:"reports"`
}

type Report struct {
	Code          string          `json:"code"`
	ReferenceCode string          `json:"referenceCode"`
	Status        string          `json:"status"`
	ReportDetails []*ReportDetail `json:"reportDetails"`
}

type ReportDetail struct {
	ID          string  `json:"id"`
	ProductCode string  `json:"productCode"`
	ProductName string  `json:"productName"`
	BuyPrice    float64 `json:"buyPrice"`
	SellPrice   float64 `json:"sellPrice"`
	Quantity    float64 `json:"quantity"`
	BuyQuantity float64 `json:"buyQuantity"`
}

type UpdateHargaBeliRequest struct {
	TransactionDetailID string `json:"transactionDetailId"`
	BuyPrice            int64  `json:"buyPrice"`
	WebUserID           string `json:"-"`
}

type TransactionBuy struct {
	ID            string    `json:"id"`
	TransactionID string    `json:"transactionId"`
	ProductID     string    `json:"productId"`
	Quantity      int64     `json:"quantity"`
	Price         int64     `json:"price"`
	PaymentMethod string    `json:"paymentMethod"`
	CreatedTime   time.Time `json:"createdTime"`
	WebUserID     string    `json:"-"`
}

type InsertTransactionBuyRequestBulk struct {
	TransactionID string                        `json:"transactionId"`
	Details       []InsertTransactionBuyRequest `json:"details"`
	WebUserID     string                        `json:"-"`
}

type InsertTransactionBuyRequest struct {
	ProductID     string `json:"productId"`
	Quantity      int64  `json:"quantity"`
	Price         int64  `json:"price"`
	PaymentMethod string `json:"paymentMethod"`
}

type TransactionBuyStatus struct {
	ID           string `json:"id"`
	Code         string `json:"code"`
	CustomerCode string `json:"customerCode"`
	CustomerName string `json:"customerName"`
	TotalBuy     int64  `json:"totalBuy"`
	TotalSell    int64  `json:"totalSell"`
}

type TransactionCredit struct {
	PreviousMonth string                  `json:"previousMonth"`
	Days          int                     `json:"days"`
	Transactions  []TransactionCreditDate `json:"transactions"`
}

type TransactionCreditDate struct {
	CustomerCode string        `json:"customerCode"`
	CustomerName string        `json:"customerName"`
	LastCredit   int64         `json:"lastCredit"`
	Credits      map[int]int64 `json:"credits"`
}
