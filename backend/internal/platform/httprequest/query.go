package httprequest

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
)

const (
	ParamPage     = "page"
	ParamPageSize = "pageSize"

	DefaultPage     int64 = 1
	DefaultPageSize int64 = 20
	MaxPageSize     int64 = 100
)

type Pagination struct {
	Page     int64
	PageSize int64
}

// ParsePagination parses and validates page-based pagination parameters.
func ParsePagination(r *http.Request) (Pagination, error) {
	query := r.URL.Query()

	p := Pagination{
		Page:     DefaultPage,
		PageSize: DefaultPageSize,
	}

	if raw := query.Get(ParamPage); raw != "" {
		page, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || page < 1 {
			return Pagination{}, fmt.Errorf("page must be a positive integer")
		}
		p.Page = page
	}

	if raw := query.Get(ParamPageSize); raw != "" {
		pageSize, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || pageSize < 1 || pageSize > MaxPageSize {
			return Pagination{}, fmt.Errorf(
				"pageSize must be between 1 and %d",
				MaxPageSize,
			)
		}
		p.PageSize = pageSize
	}

	// Prevent overflow when converting the page to an offset.
	if p.Page-1 > math.MaxInt64/p.PageSize {
		return Pagination{}, fmt.Errorf("page is too large")
	}

	return p, nil
}

// LimitOffset converts page-based pagination into database pagination.
func (p Pagination) LimitOffset() (limit, offset int64) {
	return p.PageSize, (p.Page - 1) * p.PageSize
}
