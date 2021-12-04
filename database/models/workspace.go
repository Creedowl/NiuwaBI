package models

import (
	"github.com/Creedowl/NiuwaBI/database"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm/clause"
)

type Workspace struct {
	BaseModel
	Name      string          `json:"name"`
	Owner     uint            `json:"owner"`
	ConfigStr string          `json:"-" gorm:"column:config"`
	Config    WorkspaceConfig `json:"config" gorm:"-"`
}

type WorkspaceConfig struct {
	DB        database.DBConfig `json:"db"`
	Operators []uint            `json:"operators"`
	Tables    []Table           `json:"tables"`
}

type Table struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

func (w *Workspace) Save() (*Workspace, error) {
	config, err := jsoniter.Marshal(&w.Config)
	if err != nil {
		return nil, err
	}
	w.ConfigStr = string(config)
	err = database.GetDB().Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(w).Error
	return w, err
}

func (w *Workspace) Update() (*Workspace, error) {
	_, err := GetWorkspaceByID(w.ID)
	if err != nil {
		return nil, err
	}
	return w.Save()
}

func GetAllWorkspaces(pagination *Pagination) (*PaginationResp, error) {
	var workspaces []Workspace
	var count int64
	db := database.GetDB().Find(&workspaces).
		Limit(pagination.PageSize).Offset((pagination.PageNum - 1) * pagination.PageSize)
	err := db.Error
	if err != nil {
		return nil, err
	}
	for i := range workspaces {
		err = jsoniter.Unmarshal([]byte(workspaces[i].ConfigStr), &workspaces[i].Config)
		if err != nil {
			return nil, err
		}
	}
	db.Count(&count)
	return &PaginationResp{
		Total: count,
		Data:  workspaces,
	}, nil
}

func GetWorkspaceByID(id uint) (*Workspace, error) {
	var workspace Workspace
	err := database.GetDB().First(&workspace, id).Error
	if err != nil {
		return nil, err
	}
	err = jsoniter.Unmarshal([]byte(workspace.ConfigStr), &workspace.Config)
	if err != nil {
		return nil, err
	}
	return &workspace, nil
}
