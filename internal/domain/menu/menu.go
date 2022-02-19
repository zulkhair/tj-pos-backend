package menudomain

import "sync"

type Menu struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	MenuOrder int64  `json:"menuOrder"`
	MenuPath  string `json:"menuPath"`
	Icon      string `json:"icon"`
}

type MenuCache struct {
	sync.RWMutex
	DataMap map[string]*Menu
}
