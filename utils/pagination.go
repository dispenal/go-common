package common_utils

import (
	"net/http"
	"strconv"
)

type Pagination[T any] struct {
	Limit      int    `json:"limit,omitempty"`
	Page       int    `json:"page,omitempty"`
	Sort       string `json:"sort,omitempty"`
	TotalRows  int    `json:"totalRows,omitempty"`
	TotalPages int    `json:"totalPages,omitempty"`
	Rows       []T    `json:"rows,omitempty"`
}

func ValidatePagination[T any](r *http.Request) *Pagination[T] {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10
	}
	sortBy := r.URL.Query().Get("sortBy")
	var sortWithDirection string
	if sortBy != "" {
		direction := r.URL.Query().Get("sortDesc")
		isDesc, err := strconv.ParseBool(direction)
		if err != nil {
			isDesc = false
		}
		if isDesc {
			sortWithDirection = sortBy + " desc"
		} else {
			sortWithDirection = sortBy + " asc"
		}

	}
	return &Pagination[T]{
		Page:  page,
		Limit: pageSize,
		Sort:  sortWithDirection,
	}
}

func Paginate[T any](pagination *Pagination[T], rows []T) *Pagination[T] {
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
