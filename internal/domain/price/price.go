package pricedomain

type PriceTemplate struct {
	ID                  string                 `json:"id"`
	Name                string                 `json:"name"`
	PriceTemplateDetail []*PriceTemplateDetail `json:"priceTemplateDetail"`
}

type PriceTemplateDetail struct {
	ProductID   string  `json:"productId"`
	Price       float64 `json:"price"`
}

type ApplyToCustomerReq struct {
	TemplateID  string   `json:"templateId"`
	CustomerIDs []string `json:"customerIds"`
}
