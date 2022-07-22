package productdomain

type Product struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	UnitCode    string `json:"unitCode"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
	UnitID      string `json:"unitId"`
}
