package transactiondomain

import "time"

type DanaStatus string

const (
	DanaStatusPending  DanaStatus = "pending"
	DanaStatusApproved DanaStatus = "approved"
	DanaStatusRejected DanaStatus = "rejected"
	DanaStatusCanceled DanaStatus = "canceled"
)

type Belanja struct {
	ID          string    `json:"id"`
	Date        time.Time `json:"date"`
	WebUserID   string    `json:"webUserId"`
	ProductID   string    `json:"productId"`
	ProductCode string    `json:"productCode"`
	ProductName string    `json:"productName"`
	Quantity    int16     `json:"quantity"`
	Price       float64   `json:"price"`
	CreatedTime time.Time `json:"createdTime"`
}

type Dana struct {
	ID           string    `json:"id"`
	Date         time.Time `json:"date"`
	WebUserID    string    `json:"webUserId"`
	SaldoAwal    float64   `json:"saldoAwal"`
	DanaTambahan float64   `json:"danaTambahan"`
	CreatedTime  time.Time `json:"createdTime"`
}

type DanaTransaction struct {
	ID          string     `json:"id"`
	Date        time.Time  `json:"date"`
	Sender      string     `json:"sender"`
	Receiver    string     `json:"receiver"`
	Amount      float64    `json:"amount"`
	Status      DanaStatus `json:"status"`
	CreatedTime time.Time  `json:"createdTime"`
}

type PenjualanTunai struct {
	ID          string    `json:"id"`
	Date        time.Time `json:"date"`
	WebUserID   string    `json:"webUserId"`
	ProductID   string    `json:"productId"`
	ProductCode string    `json:"productCode"`
	ProductName string    `json:"productName"`
	Quantity    int16     `json:"quantity"`
	Price       float64   `json:"price"`
	CreatedTime time.Time `json:"createdTime"`
}

type Operasional struct {
	ID          string    `json:"id"`
	Date        time.Time `json:"date"`
	WebUserID   string    `json:"webUserId"`
	Description string    `json:"description"`
	Quantity    int16     `json:"quantity"`
	Price       float64   `json:"price"`
	CreatedTime time.Time `json:"createdTime"`
}

type DanaRequest struct {
	ID           string  `json:"id"`
	Date         string  `json:"date"`
	SaldoAwal    float64 `json:"saldoAwal"`
	DanaTambahan float64 `json:"danaTambahan"`
}

type DanaTransactionRequest struct {
	Date     string  `json:"date"`
	Receiver string  `json:"receiver"`
	Amount   float64 `json:"amount"`
}

type DanaInquiryResponse struct {
	ID           string                           `json:"id"`
	SaldoAwal    float64                          `json:"saldoAwal"`
	DanaTambahan float64                          `json:"danaTambahan"`
	DanaMasuk    []DanaTransactionInquiryResponse `json:"danaMasuk"`
	DanaKeluar   []DanaTransactionInquiryResponse `json:"danaKeluar"`
}

type DanaTransactionInquiryResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Amount      float64    `json:"amount"`
	Status      DanaStatus `json:"status"`
	CreatedTime time.Time  `json:"createdTime"`
}

type WebUserMobile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TrxCreateRequest struct {
	Date      string  `json:"date"`
	ProductID string  `json:"productId"`
	Quantity  int16   `json:"quantity"`
	Price     float64 `json:"price"`
}

type TrxInquiryResponse struct {
	ID          string  `json:"id"`
	ProductName string  `json:"productName"`
	Quantity    int16   `json:"quantity"`
	Price       float64 `json:"price"`
}

type TrxCreateOperasionalRequest struct {
	Date        string  `json:"date"`
	Description string  `json:"description"`
	Quantity    int16   `json:"quantity"`
	Price       float64 `json:"price"`
}

type TrxInquiryOperasionalResponse struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Quantity    int16   `json:"quantity"`
	Price       float64 `json:"price"`
}

type SaldoResponse struct {
	SaldoAwal    float64 `json:"saldoAwal"`
	DanaTambahan float64 `json:"danaTambahan"`
	DanaMasuk    float64 `json:"danaMasuk"`
	Belanja      float64 `json:"belanja"`
	Operasional  float64 `json:"operasional"`
	DanaKeluar   float64 `json:"danaKeluar"`
	SaldoAkhir   float64 `json:"saldoAkhir"`
}

type RekapitulasiResponse struct {
	Rekapitulasi     []RekapitulasiDetail `json:"rekapitulasi"`
	TotalBelanja     float64              `json:"totalBelanja"`
	TotalOperasional float64              `json:"totalOperasional"`
}

type RekapitulasiDetail struct {
	Name        string  `json:"name"`
	Belanja     float64 `json:"belanja"`
	Operasional float64 `json:"operasional"`
}
