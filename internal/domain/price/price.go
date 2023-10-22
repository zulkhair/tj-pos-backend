package pricedomain

type PriceTemplate struct {
	ID                  string                 `json:"id"`
	Name                string                 `json:"name"`
	AppliedTo           string                 `json:"appliedTo"`
	PriceTemplateDetail []*PriceTemplateDetail `json:"priceTemplateDetail"`
}

type PriceTemplateDetail struct {
	ID        string  `json:"id"`
	ProductID string  `json:"productId"`
	Price     float64 `json:"price"`
	Checked   bool    `json:"checked"`
}

type ApplyToCustomerReq struct {
	TemplateID  string   `json:"templateId"`
	CustomerIDs []string `json:"customerIds"`
}

type ApplyToTrxReq struct {
	TemplateID string `json:"templateId"`
	Date       string `json:"date"`
}

type Download struct {
	TemplateID        string   `json:"templateId"`
	TemplateDetailIDs []string `json:"templateDetailIds"`
}

type DeleteTemplateReq struct {
	TemplateID string `json:"templateId"`
}
