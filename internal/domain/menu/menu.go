package menudomain

type Menu struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	MenuOrder int    `json:"menuOrder"`
	MenuPath  string `json:"menuPath"`
	Icon      string `json:"icon"`
}
