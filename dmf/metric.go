package dmf

import "fmt"

type MetricI interface {
	GetName() string
	GetLabel() string
	GenStmt(bool) string
}

type Metric struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Field string `json:"field"`
	Table string `json:"table"`
	Aggr  string `json:"aggr"`
	Label string `json:"label"`
}

func (m *Metric) GetName() string {
	return m.Name
}

func (m *Metric) GetLabel() string {
	return m.Label
}

func (m *Metric) GenStmt(asElement bool) string {
	if !asElement {
		return fmt.Sprintf("%s(`%s`.`%s`) as `%s`", m.Aggr, m.Table, m.Field, m.Name)
	} else {
		return fmt.Sprintf("%s(`%s`.`%s`)", m.Aggr, m.Table, m.Field)
	}
}

type EquationMetric struct {
	Type          string       `json:"type"`
	Name          string       `json:"name"`
	Label         string       `json:"label"`
	Formula       string       `json:"formula"`
	Elements      []DimensionI `json:"-"`
	ElementFields []string     `json:"elements"`
}

func (e *EquationMetric) GetName() string {
	return e.Name
}

func (e *EquationMetric) GetLabel() string {
	return e.Label
}

func (e *EquationMetric) GenStmt(asElement bool) string {
	elements := make([]interface{}, len(e.Elements))
	for i, element := range e.Elements {
		elements[i] = element.GenStmt(true)
	}
	if !asElement {
		elements = append(elements, e.Name)
		return fmt.Sprintf(e.Formula+" as `%s`", elements...)
	} else {
		return fmt.Sprintf(e.Formula, elements...)
	}
}
