package charts

type ChartDataY struct {
	Smooth bool          `json:"smooth"`
	Name   string        `json:"name"`
	Type   string        `json:"type"`
	Data   []interface{} `json:"data"`
}

type ChartDataX struct {
	Name string        `json:"name"`
	Type string        `json:"type"`
	Data []interface{} `json:"data"`
}

type Grid struct {
	Left         string `json:"left,omitempty"`
	Right        string `json:"right,omitempty"`
	Bottom       string `json:"bottom,omitempty"`
	ContainLabel bool   `json:"containLabel,omitempty"`
}
type void struct{}

type EchartsTitle struct {
	Text string `json:"text"`
}

func CompileGridStyle(left int, right int, bottom int, containLabel bool) Grid {
	return Grid{
		Left:         string(rune(left)) + "%",
		Right:        string(rune(right)) + "%",
		Bottom:       string(rune(bottom)) + "%",
		ContainLabel: containLabel,
	}
}
func CompileSequentialData(name []string, datatype []string, data [][]interface{}, options []yOption) []ChartDataY {
	results := make([]ChartDataY, 0)
	l := len(name)
	for i := 0; i < l; i++ {
		results = append(results, CompileData(name[i], datatype[i], data[i], options[i].Smooth))
	}
	return results
}

func CompileData(name string, datatype string, data []interface{}, smooth bool) ChartDataY {
	c := ChartDataY{
		Name:   name,
		Type:   datatype,
		Data:   data,
		Smooth: smooth,
	}
	return c
}

func CompileDataX(name string, datatype string, data []interface{}) ChartDataX {
	c := ChartDataX{
		Name: name,
		Type: datatype,
		Data: data,
	}
	return c
}

func CompileYAxis(yDataType string) (m map[string]interface{}) {
	m = map[string]interface{}{}
	switch yDataType {
	case "value":
		m["type"] = "value"
	case "category":
		m["type"] = "category"
	default:
		m["type"] = "value"
	}
	return
}

func CompileTitle(title string) EchartsTitle {
	return EchartsTitle{Text: title}
}

func CompileFeatures(saveImage bool) (m map[string]interface{}) {
	m = map[string]interface{}{}
	featureMap := make(map[string]interface{})
	if saveImage {
		featureMap["saveAsImage"] = void{}
	}
	m["feature"] = featureMap
	return
}

func CompileTooltips(isAxis bool) (m map[string]interface{}) {
	m = map[string]interface{}{}
	if isAxis {
		m["trigger"] = "axis"
	}
	return
}

func CompileLegends(kv []Kv, fields []string) (m map[string]interface{}) {
	m = map[string]interface{}{}
	data := make([]string, 0)
	for _, field := range fields {
		exist := false
		for _, item := range kv {
			if item.Key == field {
				data = append(data, item.Label)
				exist = false
				break
			}
		}
		if exist == false {
			data = append(data, field)
		}
	}
	m["data"] = data
	return
}
