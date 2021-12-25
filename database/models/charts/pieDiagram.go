package charts

import "gorm.io/gorm"

type PieDiagram struct {
	ChartBase
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
