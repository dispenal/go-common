package common_utils

import (
	"net/http"
	"strconv"
)

type Pagination[T any] struct {
	SearchField string `json:"searchField,omitempty"`
	SearchValue string `json:"searchValue,omitempty"`
	Limit       int    `json:"limit,omitempty"`
	Page        int    `json:"page,omitempty"`
	Sort        string `json:"sort,omitempty"`
	TotalRows   int    `json:"totalRows,omitempty"`
	TotalPages  int    `json:"totalPages,omitempty"`
	Rows        []T    `json:"rows,omitempty"`
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
	if sortBy == "" {
		sortBy = "createdAt"
	}
	searchField := r.URL.Query().Get("searchField")
	searchValue := r.URL.Query().Get("searchValue")
	return &Pagination[T]{
		SearchField: searchField,
		SearchValue: searchValue,
		Page:        page,
		Limit:       pageSize,
		Sort:        sortBy,
	}
}

func Paginate[T any](pagination *Pagination[T], rows []T) *Pagination[T] {
	totalRows := pagination.TotalRows
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
