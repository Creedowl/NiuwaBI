package charts

import (
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

func (t *DataTable) GetType() string {
	return DataTable_
}
