package charts

import "gorm.io/gorm"

const (
	DataTable_           = "table"
	DatatableLineDiagram = "line"
	DatatablePieDiagram  = "pie"
)

type Pos struct {
	X int    `json:"x"`
	Y int    `json:"y"`
	W int    `json:"w"`
	H int    `json:"h"`
	I string `json:"i"`
}

type Kv struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

type ChartBase struct {
	Type string `json:"type"`
	Sql  string `json:"sql"`
	Name string `json:"name"`
	Pos  Pos    `json:"pos"`
	Kv   []Kv   `json:"kv"`
}

type Chart interface {
	Execute(db *gorm.DB) (interface{}, error)
	GetType() string
}
