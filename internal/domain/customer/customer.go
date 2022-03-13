package customerdomain

import (
	productdomain "dromatech/pos-backend/internal/domain/product"
	unitdomain "dromatech/pos-backend/internal/domain/unit"
	"time"
)

type Customer struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type SellPrice struct {
	Date     time.Time             `json:"date"`
	Customer Customer              `json:"customer"`
	Product  productdomain.Product `json:"product"`
	Unit     unitdomain.Unit       `json:"unit"`
	Price    float64               `json:"price"`
}

type SellPriceRequest struct {
	Date       string                   `json:"date"`
	CustomerId string                   `json:"customerId"`
	UnitId     string                   `json:"unitId"`
	Prices     []SellPriceDetailRequest `json:"prices"`
}

type SellPriceDetailRequest struct {
	ProductID string  `json:"productId"`
	Price     float64 `json:"price"`
}

type SellPriceResponse struct {
	ProductID   string  `json:"productId"`
	ProductCode string  `json:"productCode"`
	ProductName string  `json:"productName"`
	ProductDesc string  `json:"productDesc"`
	Price       float64 `json:"price"`
}
