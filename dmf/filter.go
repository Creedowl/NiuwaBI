package dmf

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

var OpMap map[string]string

func init() {
	OpMap = map[string]string{
		">":           ">",
		">=":          ">=",
		"=":           "=",
		"<":           "<",
		"<=":          "<=",
		"in":          "in",
		"not in":      "not in",
		"is null":     "is null",
		"is not null": "is not null",
	}
}

type Filter struct {
	Field string      `json:"field"`
	Dim   DimensionI  `json:"-"`
	Met   MetricI     `json:"-"`
	Op    string      `json:"op"`
	Value interface{} `json:"value"`
}

func (f *Filter) Apply(db *gorm.DB) (*gorm.DB, error) {
	field := ""
	stmt := ""
	if f.Dim != nil {
		field = f.Dim.GetName()
	} else if f.Met != nil {
		field = f.Met.GetName()
	} else {
		return db, fmt.Errorf("filter must contain one of dimensions or metrics")
	}
	op, ok := OpMap[f.Op]
	if !ok {
		return db, fmt.Errorf("unsupport op %s", f.Op)
	}
	switch op {
	case ">", ">=", "=", "<", "<=":
		stmt = fmt.Sprintf("%s %s ?", field, op)
	case "in", "not in":
		if reflect.TypeOf(f.Value).Kind() != reflect.Slice {
			return db, fmt.Errorf("'in' or 'not in' must have array as value")
		}
		stmt = fmt.Sprintf("%s %s ?", field, op)
	case "is null", "is not null":
		stmt = fmt.Sprintf("%s %s", field, op)
	}
	if f.Dim != nil {
		db.Where(stmt, f.Value)
	} else if f.Met != nil {
		db.Having(stmt, f.Value)
	}
	return db, nil
}
