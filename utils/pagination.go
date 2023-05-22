package utils

import (
	"gorm.io/gorm"
	"math"
)

type Pagination struct {
	Paginate  int   `json:"paginate,omitempty"`
	Page      int   `json:"page,omitempty"`
	PerPage   int   `json:"per_page,omitempty"`
	PageCount int   `json:"page_count"`
	Total     int64 `json:"total"`
	Next      int   `json:"next,omitempty"`
	Previous  int   `json:"previous,omitempty"`
}

func Paginate(value interface{}, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalData int64
	db.Model(value).Count(&totalData)

	pagination.Total = totalData
	pagination.PageCount = int(math.Ceil(float64(totalData) / float64(pagination.PerPage)))

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	return p.PerPage
}
