package charts

import "gorm.io/gorm"

const (
	DataTable_ = "table"
)

type Pos struct {
	X int    `json:"x"`
	Y int    `json:"y"`
	W int    `json:"w"`
	H int    `json:"h"`
	I string `json:"i"`
}

type ChartBase struct {
	Type string `json:"type"`
	Sql  string `json:"sql"`
	Name string `json:"name"`
	Pos  Pos    `json:"pos"`
}

type Chart interface {
	Execute(db *gorm.DB) (interface{}, error)
	GetType() string
}
