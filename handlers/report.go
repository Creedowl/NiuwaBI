package handlers

import (
	"github.com/Creedowl/NiuwaBI/database/models"
	"github.com/gin-gonic/gin"
)

type ReportIDParam struct {
	ID uint `json:"id"`
}

func CreateReport(_ *gin.Context, report models.Report) (*models.Report, error) {
	return report.Save()
}

type GetAllReportsParam struct {
	models.Pagination
	WorkspaceID uint `json:"workspace_id"`
}

func GetAllReports(_ *gin.Context, param GetAllReportsParam) (*models.PaginationResp, error) {
	return models.GetAllReports(&param.Pagination, param.WorkspaceID)
}

func GetReport(_ *gin.Context, param ReportIDParam) (*models.Report, error) {
	return models.GetReportByID(param.ID)
}

func UpdateReport(_ *gin.Context, workspace models.Report) (*models.Report, error) {
	return workspace.Update()
}

func ExecuteReport(_ *gin.Context, param ReportIDParam) (interface{}, error) {
	report, err := models.GetReportByID(param.ID)
	if err != nil {
		return nil, err
	}
	res, err := report.Execute()
	if err != nil {
		return nil, err
	}
	return res, nil
}
