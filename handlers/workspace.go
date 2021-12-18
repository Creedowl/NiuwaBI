package handlers

import (
	"github.com/Creedowl/NiuwaBI/database"
	"github.com/Creedowl/NiuwaBI/database/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type WorkspaceIDParam struct {
	ID uint `json:"id"`
}

func CreateWorkspace(_ *gin.Context, workspace models.Workspace) (*models.Workspace, error) {
	return workspace.Save()
}

func GetAllWorkspaces(_ *gin.Context, pagination models.Pagination) (*models.PaginationResp, error) {
	return models.GetAllWorkspaces(&pagination)
}

func GetWorkspace(_ *gin.Context, param WorkspaceIDParam) (*models.Workspace, error) {
	return models.GetWorkspaceByID(param.ID)
}

func UpdateWorkspace(_ *gin.Context, workspace models.Workspace) (*models.Workspace, error) {
	return workspace.Update()
}

func RemoveWorkspace(_ *gin.Context, param WorkspaceIDParam) (*DumbResp, error) {
	return &DumbResp{OK: true}, models.RemoveWorkspace(param.ID)
}

type TestConnParam struct {
	WorkspaceIDParam
	database.DBConfig
}

type DumbResp struct {
	OK bool `json:"ok"`
}

func TestConn(_ *gin.Context, param TestConnParam) (*DumbResp, error) {
	var dbConfig database.DBConfig
	if param.WorkspaceIDParam.ID != 0 {
		workspace, err := models.GetWorkspaceByID(param.ID)
		if err != nil {
			return nil, err
		}
		logrus.Infof("db config: %+v", workspace.Config.DB)
		dbConfig = workspace.Config.DB
	} else {

		dbConfig = param.DBConfig
	}
	_, err := dbConfig.TestConn()
	if err != nil {
		return nil, err
	}
	return &DumbResp{OK: true}, nil
}
