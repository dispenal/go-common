package common_utils

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Pagination[T any] struct {
	Limit      int    `json:"limit,omitempty"`
	Page       int    `json:"page,omitempty"`
	Sort       string `json:"sort,omitempty"`
	TotalRows  int    `json:"totalRows,omitempty"`
	TotalPages int    `json:"totalPages,omitempty"`
	Rows       T      `json:"rows,omitempty"`
}

func ValidatePagination(r *http.Request) *Pagination[any] {
	pageStr := chi.URLParam(r, "page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	pageSizeStr := chi.URLParam(r, "page")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10
	}
	sortBy := chi.URLParam(r, "sortBy")
	var sortWithDirection string
	if sortBy != "" {
		direction := chi.URLParam(r, "sortDesc")
		if direction != "" {
			if direction == "true" || direction == "1" {
				sortWithDirection = sortBy + " desc"
			} else if direction == "false" || direction == "0" {
				sortWithDirection = sortBy + " asc"
			}
		}

	}
	return &Pagination[any]{
		Page:  page,
		Limit: pageSize,
		Sort:  sortWithDirection,
	}
}

func Paginate[T []any](pagination Pagination[T], rows T) *Pagination[T] {
	totalRows := len(rows)
	limit := pagination.Limit
	page := pagination.Page

	totalPages := totalRows / limit
	if totalRows%limit != 0 {
		totalPages++
	}
	return &Pagination[T]{
		Limit:      limit,
		Page:       page,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Rows:       rows,
	}
}
