package dmf

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"strings"
)

type DMF struct {
	Table      string       `json:"table"` // 目前只支持单表查询
	Dimensions []DimensionI `json:"dimensions"`
	Metrics    []MetricI    `json:"metrics"`
}

func (d *DMF) GetDimensionByName(name string) DimensionI {
	for _, dimension := range d.Dimensions {
		if dimension.GetName() == name {
			return dimension
		}
	}
	return nil
}

func (d *DMF) GetMetricByName(name string) MetricI {
	for _, metric := range d.Metrics {
		if metric.GetName() == name {
			return metric
		}
	}
	return nil
}

func (d *DMF) Check() error {
	// dimensions
	for _, dimension := range d.Dimensions {
		if ed, ok := dimension.(*EquationDimension); ok {
			for _, field := range ed.ElementFields {
				dim := d.GetDimensionByName(field)
				if dim != nil {
					ed.Elements = append(ed.Elements, dim)
				} else {
					return fmt.Errorf("dimension %s not found", field)
				}
			}
		}
	}

	// metrics
	for _, metric := range d.Metrics {
		if em, ok := metric.(*EquationMetric); ok {
			for _, field := range em.ElementFields {
				met := d.GetMetricByName(field)
				if met != nil {
					em.Elements = append(em.Elements, met)
				} else {
					return fmt.Errorf("metric %s not found", field)
				}
			}
		}
	}

	return nil
}

func (d *DMF) UnmarshalJSON(data []byte) error {
	c := struct {
		Table      string                `json:"table"`
		Dimensions []jsoniter.RawMessage `json:"dimensions"`
		Metrics    []jsoniter.RawMessage `json:"metrics"`
	}{}
	t := struct {
		Type string `json:"type"`
	}{}
	err := jsoniter.Unmarshal(data, &c)
	if err != nil {
		return err
	}
	d.Table = c.Table

	for _, ct := range c.Dimensions {
		err = jsoniter.Unmarshal(ct, &t)
		if err != nil {
			return err
		}
		switch t.Type {
		case "dimension":
			var dimension Dimension
			err = jsoniter.Unmarshal(ct, &dimension)
			if err != nil {
				return err
			}
			d.Dimensions = append(d.Dimensions, &dimension)
		case "equation_dimension":
			var ed EquationDimension
			err = jsoniter.Unmarshal(ct, &ed)
			if err != nil {
				return err
			}
			d.Dimensions = append(d.Dimensions, &ed)
		default:
			return fmt.Errorf("unkown dmf type %s", t.Type)
		}
	}

	for _, ct := range c.Metrics {
		err = jsoniter.Unmarshal(ct, &t)
		if err != nil {
			return err
		}
		switch t.Type {
		case "metric":
			var m Metric
			err = jsoniter.Unmarshal(ct, &m)
			if err != nil {
				return err
			}
			d.Metrics = append(d.Metrics, &m)
		case "equation_metric":
			var em EquationMetric
			err = jsoniter.Unmarshal(ct, &em)
			if err != nil {
				return err
			}
			d.Metrics = append(d.Metrics, &em)
		default:
			return fmt.Errorf("unkown dmf type %s", t.Type)
		}
	}
	return nil
}

func (d *DMF) Execute(db *gorm.DB, fields []string, filters []Filter) (interface{}, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("fields are empty")
	}

	var dimensions []DimensionI
	var metrics []MetricI
	var selects []string
	var err error

	for _, field := range fields {
		dimension := d.GetDimensionByName(field)
		if dimension != nil {
			dimensions = append(dimensions, dimension)
			selects = append(selects, dimension.GenStmt(false))
			continue
		}

		metric := d.GetMetricByName(field)
		if metric != nil {
			metrics = append(metrics, metric)
			selects = append(selects, metric.GenStmt(false))
			continue
		}

		return nil, fmt.Errorf("field %s not found in dimensions and metrics", field)
	}

	db = db.Table(d.Table).Select(strings.Join(selects, ", "))

	for _, dimension := range dimensions {
		db = db.Group(dimension.GetName())
	}

	for _, filter := range filters {
		db, err = filter.Apply(db)
		if err != nil {
			return nil, err
		}
	}

	var results []map[string]interface{}
	err = db.Find(&results).Error

	return results, err
}
