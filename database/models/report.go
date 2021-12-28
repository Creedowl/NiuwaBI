package models

import (
	"fmt"
	"github.com/Creedowl/NiuwaBI/database"
	"github.com/Creedowl/NiuwaBI/database/models/charts"
	"github.com/Creedowl/NiuwaBI/dmf"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm/clause"
)

type Report struct {
	BaseModel
	WorkspaceID uint         `json:"workspace_id"`
	Name        string       `json:"name"`
	Owner       uint         `json:"owner"`
	ConfigStr   string       `json:"-" gorm:"column:config"`
	Config      ReportConfig `json:"config" gorm:"-"`
}

type ReportConfig struct {
	Charts []charts.Chart `json:"charts"`
	Dmf    dmf.DMF        `json:"dmf"`
}

type ChartData struct {
	Chart charts.Chart `json:"chart"`
	Data  interface{}  `json:"data"`
}

func (r *ReportConfig) UnmarshalJSON(data []byte) error {
	c := struct {
		Charts []jsoniter.RawMessage `json:"charts"`
		Dmf    dmf.DMF               `json:"dmf"`
	}{}
	t := struct {
		Type string `json:"type"`
	}{}
	err := jsoniter.Unmarshal(data, &c)
	if err != nil {
		return err
	}
	r.Dmf = c.Dmf

	for _, ct := range c.Charts {
		err = jsoniter.Unmarshal(ct, &t)
		if err != nil {
			return err
		}
		switch t.Type {
		case charts.DataTable_:
			var dt charts.DataTable
			err = jsoniter.Unmarshal(ct, &dt)
			if err != nil {
				return err
			}
			r.Charts = append(r.Charts, &dt)
		case charts.DatatableLineDiagram:
			var dt charts.LineDiagram
			err = jsoniter.Unmarshal(ct, &dt)
			if err != nil {
				return err
			}
			r.Charts = append(r.Charts, &dt)
		case charts.DatatablePieDiagram:
			var dt charts.PieDiagram
			err = jsoniter.Unmarshal(ct, &dt)
			if err != nil {
				return err
			}
			r.Charts = append(r.Charts, &dt)
		default:
			return fmt.Errorf("unkown chart type %s", t.Type)
		}
	}
	return nil
}

func (r *Report) Check() error {
	// check dimensions and metrics
	err := r.Config.Dmf.Check()
	if err != nil {
		return err
	}

	// check filters
	for _, chart := range r.Config.Charts {
		err = chart.GetChartBase().Check(&r.Config.Dmf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Report) Save() (*Report, error) {
	config, err := jsoniter.Marshal(&r.Config)
	if err != nil {
		return nil, err
	}
	r.ConfigStr = string(config)

	_, err = GetWorkspaceByID(r.WorkspaceID)
	if err != nil {
		return nil, fmt.Errorf("workspace %d not found", r.WorkspaceID)
	}

	err = r.Check()
	if err != nil {
		return nil, err
	}

	// update kv
	for _, chart := range r.Config.Charts {
		if chart.GetChartType() == charts.DMFChart {
			err = chart.UpdateKv(&r.Config.Dmf)
			if err != nil {
				return nil, err
			}
		}
	}
	err = database.GetDB().Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(r).Error
	return r, err
}

func (r *Report) Update() (*Report, error) {
	_, err := GetReportByID(r.ID)
	if err != nil {
		return nil, err
	}
	err = r.Check()
	if err != nil {
		return nil, err
	}

	// update kv
	for _, chart := range r.Config.Charts {
		if chart.GetChartType() == charts.DMFChart {
			err = chart.UpdateKv(&r.Config.Dmf)
			if err != nil {
				return nil, err
			}
		}
	}
	return r.Save()
}

func (r *Report) Execute() ([]ChartData, error) {
	var results []ChartData

	workspace, err := GetWorkspaceByID(r.WorkspaceID)
	if err != nil {
		return nil, err
	}
	db, err := database.GetCachedDB(r.WorkspaceID, &workspace.Config.DB)
	if err != nil {
		return nil, err
	}
	err = r.Check()
	if err != nil {
		return nil, err
	}

	for _, chart := range r.Config.Charts {
		var res interface{}
		var err error

		switch chart.GetChartType() {
		case charts.SqlChart:
			res, err = chart.Execute(db)
		case charts.DMFChart:
			res, err = chart.ExecuteDmf(db, &r.Config.Dmf)
		default:
			return nil, fmt.Errorf("unknow chart type %s", chart.GetChartType())
		}
		if err != nil {
			return nil, err
		}
		results = append(results, ChartData{
			Chart: chart,
			Data:  res,
		})
	}
	return results, nil
}

func GetAllReports(pagination *Pagination, workspaceID uint) (*PaginationResp, error) {
	var reports []Report
	var count int64
	db := database.GetDB().Model(&Report{})
	if workspaceID != 0 {
		db.Where("workspace_id = ?", workspaceID)
	}
	db.Find(&reports).
		Limit(pagination.PageSize).Offset((pagination.PageNum - 1) * pagination.PageSize)
	err := db.Error
	if err != nil {
		return nil, err
	}
	for i := range reports {
		err = jsoniter.Unmarshal([]byte(reports[i].ConfigStr), &reports[i].Config)
		if err != nil {
			return nil, err
		}
	}
	db.Count(&count)
	return &PaginationResp{
		Total: count,
		Data:  reports,
	}, nil
}

func GetReportByID(id uint) (*Report, error) {
	var report Report
	err := database.GetDB().First(&report, id).Error
	if err != nil {
		return nil, err
	}
	err = jsoniter.Unmarshal([]byte(report.ConfigStr), &report.Config)
	if err != nil {
		return nil, err
	}
	return &report, nil
}
