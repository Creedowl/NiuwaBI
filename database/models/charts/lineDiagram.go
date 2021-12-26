package charts

import (
	"errors"
	"github.com/Creedowl/NiuwaBI/dmf"
	"gorm.io/gorm"
)

type LineDiagram struct {
	ChartBase
	X            string    `json:"x"`
	Y            []string  `json:"y"`
	XAxisType    string    `json:"xAxisType"`
	YDataType    string    `json:"yDataType"`
	Datatype     []string  `json:"datatype"`
	YExtraOption []yOption `json:"yExtraOption"`
}

func (t *LineDiagram) ExecuteDmf(db *gorm.DB, dmf *dmf.DMF) (interface{}, error) {
	return dmf.Execute(db, t.Fields, t.Filters)
}

func (t *LineDiagram) UpdateKv(dmf *dmf.DMF) error {
	panic("implement me")
}

func (t *LineDiagram) GetChartBase() *ChartBase {
	return &t.ChartBase
}

type yOption struct {
	Smooth bool `json:"smooth"`
}

func (t *LineDiagram) Execute(db *gorm.DB) (interface{}, error) {
	var results []map[string]interface{}
	err := db.Raw(t.Sql).Find(&results).Error
	if err != nil {
		return nil, err
	}
	//Generate echarts json option
	keys := make(map[string]bool)
	for _, key := range t.Y {
		keys[key] = true
	}
	xData := make([]interface{}, 0)
	yData := make([][]interface{}, len(t.Y)) // lines
	kvCache := map[string]string{}
	for _, kv := range t.Kv {
		kvCache[kv.Key] = kv.Label
	}
	keyIndex := make(map[string]int)
	xName := t.X
	first := true
	var yName []string
	if v, ok := kvCache[xName]; ok {
		xName = v
	}

	for _, result := range results {
		index := 0
		for k, v := range result {
			if k == t.X {
				xData = append(xData, v)
			}
			if _, exist := keys[k]; exist {
				if first {
					keyIndex[k] = index
					if vv, ok := kvCache[k]; ok {
						yName = append(yName, vv)
					} else {
						yName = append(yName, k)
					}
				}
				index++
				yData[keyIndex[k]] = append(yData[keyIndex[k]], v)
			}
		}
		if index != len(t.Y) {
			return nil, errors.New("data unmatched with keys")
		}
		first = false
	}

	var echartsJsonData = map[string]interface{}{}
	echartsJsonData["title"] = CompileTitle(t.Name)
	echartsJsonData["tooltip"] = CompileTooltips(true)
	echartsJsonData["toolbox"] = CompileFeatures(true)
	echartsJsonData["legend"] = CompileLegends(t.Kv, t.Y)
	//X-axis
	echartsJsonData["xAxis"] = CompileDataX(xName, t.XAxisType, xData)
	echartsJsonData["yAxis"] = CompileYAxis(t.YDataType)
	//Data
	echartsJsonData["series"] = CompileSequentialData(yName, t.Datatype, yData, t.YExtraOption)
	//default grid style
	echartsJsonData["grid"] = CompileGridStyle(3, 4, 3, true)
	return echartsJsonData, nil
}

func (t *LineDiagram) GetType() string {
	return DatatableLineDiagram
}
