package internal

import (
	"net/http"
	"strconv"
)

const (
	PAGINATIN_DEFAULT_PAGE         = 1
	PAGINATIN_DEFAULT_PER_PAGE     = 20
	PAGINATIN_DEFAULT_MAX_PER_PAGE = 100
)

type Pagination struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func ParsePaginationRequest(r *http.Request) *Pagination {

	pageQuery := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageQuery)
	if err != nil || page < 1 {
		page = PAGINATIN_DEFAULT_PAGE
	}

	perPageQuery := r.URL.Query().Get("per_page")
	perPage, err := strconv.Atoi(perPageQuery)
	if err != nil || perPage < 1 {
		page = PAGINATIN_DEFAULT_PER_PAGE
	}

	if perPage > PAGINATIN_DEFAULT_MAX_PER_PAGE {
		perPage = PAGINATIN_DEFAULT_MAX_PER_PAGE
	}

	return &Pagination{
		Page:    page,
		PerPage: perPage,
	}
}

func (p *Pagination) GetOffset() int32 {
	return int32((p.Page - 1) * p.PerPage)
}

func (p *Pagination) GetLimit() int32 {
	return int32(p.PerPage)
}
