package dmf

import (
	"fmt"
)

type DimensionI interface {
	GetName() string
	GetLabel() string
	GenStmt(bool) string
}

type Dimension struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Field string `json:"field"`
	Table string `json:"table"`
	Label string `json:"label"`
}

func (d *Dimension) GetName() string {
	return d.Name
}

func (d *Dimension) GetLabel() string {
	return d.Label
}

func (d *Dimension) GenStmt(asElement bool) string {
	if !asElement {
		return fmt.Sprintf("`%s`.`%s` as `%s`", d.Table, d.Field, d.Name)
	} else {
		return fmt.Sprintf("`%s`.`%s`", d.Table, d.Field)
	}
}

type EquationDimension struct {
	Type          string       `json:"type"`
	Name          string       `json:"name"`
	Label         string       `json:"label"`
	Formula       string       `json:"formula"`
	Elements      []DimensionI `json:"-"`
	ElementFields []string     `json:"elements"`
}

func (e *EquationDimension) GetName() string {
	return e.Name
}

func (e *EquationDimension) GetLabel() string {
	return e.Label
}

func (e *EquationDimension) GenStmt(asElement bool) string {
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
