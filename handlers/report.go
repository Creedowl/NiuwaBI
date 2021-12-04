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

func GetAllReports(_ *gin.Context, pagination models.Pagination) (*models.PaginationResp, error) {
	return models.GetAllReports(&pagination)
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
