package models

import (
	"fmt"
	"github.com/Creedowl/NiuwaBI/database"
	"github.com/Creedowl/NiuwaBI/database/models/charts"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm/clause"
)

const (
	SqlReport = "sql"
	DMFReport = "dmf"
)

type Report struct {
	BaseModel
	WorkspaceID uint         `json:"workspace_id"`
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Owner       uint         `json:"owner"`
	ConfigStr   string       `json:"-" gorm:"column:config"`
	Config      ReportConfig `json:"config" gorm:"-"`
}

type ReportConfig struct {
	Charts []charts.Chart `json:"charts"`
}

type ChartData struct {
	Chart charts.Chart `json:"chart"`
	Data  interface{}  `json:"data"`
}

func (r *ReportConfig) UnmarshalJSON(data []byte) error {
	c := struct {
		Charts []jsoniter.RawMessage `json:"charts"`
	}{}
	t := struct {
		Type string `json:"type"`
	}{}
	err := jsoniter.Unmarshal(data, &c)
	if err != nil {
		return err
	}

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
	for _, c := range r.Config.Charts {
		res, err := c.Execute(db) //Generate chart data
		if err != nil {
			return nil, err
		}
		results = append(results, ChartData{
			Chart: c,
			Data:  res, //Generate by Execute
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
