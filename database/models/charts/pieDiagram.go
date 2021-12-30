package charts

import (
	"errors"
	"fmt"
	"github.com/Creedowl/NiuwaBI/dmf"
	"gorm.io/gorm"
)

type PieDiagram struct {
	ChartBase
	Data       []string `json:"data"`   // works for oneRow
	OneRow     bool     `json:"oneRow"` //if oneRow is enabled, only works for one line data
	NameField  string   `json:"nameField"`
	ValueField string   `json:"valueField"`
	RoseType   bool     `json:"roseType"`
}

func (t *PieDiagram) ExecuteDmf(db *gorm.DB, dmf *dmf.DMF) (interface{}, error) {
	results, err := dmf.Execute(db, t.Fields, t.Filters)
	if err != nil {
		return nil, err
	}
	return t.genEchartsJsonData(results.([]map[string]interface{}))
}

func (t *PieDiagram) UpdateKv(dmf *dmf.DMF) error {
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

func (t *PieDiagram) GetChartBase() *ChartBase {
	return &t.ChartBase
}

func (t *PieDiagram) Execute(db *gorm.DB) (interface{}, error) {
	var results []map[string]interface{}
	err := db.Raw(t.Sql).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return t.genEchartsJsonData(results)
}

func (t *PieDiagram) genEchartsJsonData(results []map[string]interface{}) (interface{}, error) {
	kvCache := map[string]string{}
	for _, kv := range t.Kv {
		kvCache[kv.Key] = kv.Label
	}
	if len(results) != 1 && (len(results) != 1 && len(results[0]) != 2) {
		return nil, errors.New("data can not convert to pie diagram")
	}
	var echartsJsonData = map[string]interface{}{}
	if t.SubName != "" {
		echartsJsonData["title"] = CompileTitle(t.Name, &t.SubName, "center")
	} else {
		echartsJsonData["title"] = CompileTitle(t.Name, nil, "center")
	}
	echartsJsonData["tooltip"] = CompileTooltips("item")
	series := make(map[string]interface{})
	itemStyle := make(map[string]interface{})
	if t.OneRow {
		//TODO test this
		series["data"] = CompileOneRowPieData(results[0], kvCache)
	} else {
		series["data"] = CompileTwoColumnsPieData(results, t.NameField, t.ValueField)
	}
	//echartsJsonData["legend"] = CompileLegends(t.Kv, t.Data)
	//series["roseType"] = "rose"
	if t.RoseType {
		series["roseType"] = "rose"
		itemStyle["borderRadius"] = 8
	}
	echartsJsonData["legend"] = CompileLeftLegends()
	series["emphasis"] = CompilePieEmphasis()
	series["type"] = "pie"
	series["itemStyle"] = itemStyle
	echartsJsonData["series"] = series
	echartsJsonData["toolbox"] = CompileFeatures(true, true)
	return echartsJsonData, nil
}

func (t *PieDiagram) GetType() string {
	return DatatablePieDiagram
}

func (t *PieDiagram) GetChartType() string {
	return t.ChartType
}
