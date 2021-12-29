package charts

import (
	"fmt"
	"github.com/Creedowl/NiuwaBI/dmf"
	"gorm.io/gorm"
)

const (
	DataTable_           = "table"
	DatatableLineDiagram = "line"
	DatatablePieDiagram  = "pie"
)

const (
	SqlChart = "sql"
	DMFChart = "dmf"
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
	Type      string       `json:"type"`
	ChartType string       `json:"chart_type"`
	Sql       string       `json:"sql"`
	Name      string       `json:"name"`
	SubName   string       `json:"subName"`
	Pos       Pos          `json:"pos"`
	Kv        []Kv         `json:"kv"`
	Fields    []string     `json:"fields"`
	Filters   []dmf.Filter `json:"filters"`
}

type Chart interface {
	Execute(*gorm.DB) (interface{}, error)
	ExecuteDmf(*gorm.DB, *dmf.DMF) (interface{}, error)
	GetType() string
	GetChartType() string
	UpdateKv(dmf *dmf.DMF) error
	GetChartBase() *ChartBase
}

func (c *ChartBase) Check(dmf *dmf.DMF) error {
	for i, filter := range c.Filters {
		dimension := dmf.GetDimensionByName(filter.Field)
		if dimension != nil {
			c.Filters[i].Dim = dimension
			continue
		}

		metric := dmf.GetMetricByName(filter.Field)
		if metric != nil {
			c.Filters[i].Met = metric
			continue
		}

		return fmt.Errorf("unknow dimension or metric %s or field is not set", filter.Field)
	}
	return nil
}
