package charts

import (
	"fmt"
	"github.com/Creedowl/NiuwaBI/dmf"
	"gorm.io/gorm"
)

type DataTable struct {
	ChartBase
}

func (t *DataTable) Execute(db *gorm.DB) (interface{}, error) {
	var results []map[string]interface{}
	err := db.Raw(t.Sql).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (t *DataTable) ExecuteDmf(db *gorm.DB, dmf *dmf.DMF) (interface{}, error) {
	return dmf.Execute(db, t.Fields, t.Filters)
}

func (t *DataTable) GetType() string {
	return DataTable_
}

func (t *DataTable) GetChartType() string {
	return t.ChartType
}

func (t *DataTable) UpdateKv(dmf *dmf.DMF) error {
	t.Kv = nil
	for _, field := range t.Fields {
		dimension := dmf.GetDimensionByName(field)
		if dimension != nil {
			t.Kv = append(t.Kv, Kv{
				Key:   field,
				Label: dimension.GetLabel(),
			})
			continue
		}

		metric := dmf.GetMetricByName(field)
		if metric != nil {
			t.Kv = append(t.Kv, Kv{
				Key:   field,
				Label: metric.GetLabel(),
			})
			continue
		}
		return fmt.Errorf("field %s not found in dimensions and metrics", field)
	}
	return nil
}

func (t *DataTable) GetChartBase() *ChartBase {
	return &t.ChartBase
}
