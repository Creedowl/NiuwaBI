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
	Text    string  `json:"text"`
	Subtext *string `json:"subtext"`
	Left    string  `json:"left"`
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

func CompileOneRowPieData(m map[string]interface{}, cache map[string]string) interface{} {
	var data []map[string]interface{}
	for k, v := range m {
		temp := map[string]interface{}{}
		temp["value"] = v
		if name, exist := cache[k]; exist {
			temp["name"] = name
		} else {
			temp["name"] = k
		}
	}
	return data
}

func CompileTwoColumnsPieData(m []map[string]interface{}, name string, value string) interface{} {
	var data []map[string]interface{}
	for _, s := range m {
		if len(s) != 2 {
			return nil
		}
		//s[0] -> name , s[1] -> value
		temp := map[string]interface{}{}
		temp["name"] = s[name]
		temp["value"] = s[value]
		data = append(data, temp)
	}
	return data
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

func CompileTitle(title string, subtext *string, left string) EchartsTitle {
	return EchartsTitle{Text: title, Subtext: subtext, Left: left}
}

func CompileFeatures(saveImage bool, restore bool) (m map[string]interface{}) {
	m = map[string]interface{}{}
	featureMap := make(map[string]interface{})
	if saveImage {
		featureMap["saveAsImage"] = void{}
	}
	if restore {
		featureMap["restore"] = void{}
	}
	m["feature"] = featureMap
	return
}

func CompileTooltips(trigger string) (m map[string]interface{}) {
	m = map[string]interface{}{}
	m["trigger"] = trigger
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

func CompilePieEmphasis() (m map[string]interface{}) {
	m = map[string]interface{}{}
	style := map[string]interface{}{}
	style["shadowBlur"] = 10
	style["shadowOffsetX"] = 0
	style["shadowColor"] = "rgba(0, 0, 0, 0.5)"
	m["itemStyle"] = style
	return
}

func CompileLeftLegends() (m map[string]interface{}) {
	m = map[string]interface{}{}
	m["orient"] = "vertical"
	m["left"] = "left"
	return
}
