package models

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Pagination struct {
	PageNum  int `json:"page_num"`
	PageSize int `json:"page_size"`
}

func (p *Pagination) Apply(db *gorm.DB) *gorm.DB {
	if p.PageNum != 0 && p.PageSize != 0 {
		db.Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize)
	}
	return db
}

type PaginationResp struct {
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
}
