package utils

import (
	"math"

	"gorm.io/gorm"
)

type PaginationPayload struct {
	Limit  int    `json:"limit" form:"limit"`
	Page   int    `json:"page" form:"page"`
	SortBy string `json:"sort_by" form:"sort_by"`
}

type Pagination[T any] struct {
	Limit      int    `json:"limit"`
	Page       int    `json:"page"`
	SortBy     string `json:"sort_by"`
	TotalPages int    `json:"total_pages"`
	TotalRows  int64  `json:"total_rows"`
	Rows       T      `json:"rows"`
}

func (p *Pagination[T]) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination[T]) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination[T]) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination[T]) GetSort() string {
	if p.SortBy == "" {
		p.SortBy = "Id"
	}

	return p.SortBy
}

func Paginate[T any](value T, pagination *Pagination[T], db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64

	db.Model(value).Count(&totalRows)
	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))
	pagination.TotalPages = totalPages

	return func(w *gorm.DB) *gorm.DB {
		return w.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}
