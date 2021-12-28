package charts

import (
	"github.com/Creedowl/NiuwaBI/dmf"
	"gorm.io/gorm"
)

type PieDiagram struct {
	ChartBase
}

func (t *PieDiagram) ExecuteDmf(db *gorm.DB, dmf *dmf.DMF) (interface{}, error) {
	return dmf.Execute(db, t.Fields, t.Filters)
}

func (t *PieDiagram) UpdateKv(dmf *dmf.DMF) error {
	panic("implement me")
}

func (t *PieDiagram) GetChartBase() *ChartBase {
	return &t.ChartBase
}

func (t *PieDiagram) Execute(db *gorm.DB) (interface{}, error) {
	var results []map[string]interface{}
	err := db.Raw(t.Sql).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *PieDiagram) GetType() string {
	return DatatablePieDiagram
}

func (t *PieDiagram) GetChartType() string {
	return t.ChartType
}
