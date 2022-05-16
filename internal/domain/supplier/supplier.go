package supplierdomain

import (
	productdomain "dromatech/pos-backend/internal/domain/product"
	unitdomain "dromatech/pos-backend/internal/domain/unit"
	"time"
)

type Supplier struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type BuyPrice struct {
	Date     time.Time             `json:"date"`
	Supplier Supplier              `json:"supplier"`
	Product  productdomain.Product `json:"product"`
	Unit     unitdomain.Unit       `json:"unit"`
	Price    float64               `json:"price"`
}

type BuyPriceRequest struct {
	Date       string                  `json:"date"`
	SupplierId string                  `json:"supplierId"`
	UnitId     string                  `json:"unitId"`
	Prices     []BuyPriceDetailRequest `json:"prices"`
}

type BuyPriceDetailRequest struct {
	ProductID string  `json:"productId"`
	Price     float64 `json:"price"`
}

type BuyPriceResponse struct {
	ProductID   string  `json:"productId"`
	ProductCode string  `json:"productCode"`
	ProductName string  `json:"productName"`
	ProductDesc string  `json:"productDesc"`
	Price       float64 `json:"price"`
}

type AddPriceRequest struct {
	ID         string  `json:"id"`
	Date       string  `json:"date"`
	SupplierId string  `json:"supplierId"`
	UnitId     string  `json:"unitId"`
	ProductID  string  `json:"productId"`
	Price      float64 `json:"price"`
	WebUserId  string  `json:"webUserId"`
	Latest     bool    `json:"latest"`
}

type PriceResponse struct {
	ID          string  `json:"id"`
	Date        string  `json:"date"`
	SupplierId  string  `json:"supplierId"`
	UnitId      string  `json:"unitId"`
	ProductID   string  `json:"productId"`
	Price       float64 `json:"price"`
	WebUsername string  `json:"webUsername"`
	WebUserName string  `json:"webUserName"`
}
