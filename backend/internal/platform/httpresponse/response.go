package httpresponse

import (
	"encoding/json"
	"net/http"

	"github.com/ShantanuBadmanji/saga-simulator/backend/internal/platform/httprequest"
)

type DataResponse struct {
	Data any `json:"data"`
	Meta any `json:"meta,omitempty"`
}

type PageMeta struct {
	Page       int64 `json:"page"`
	PageSize   int64 `json:"pageSize"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int64 `json:"totalPages"`
}

func NewPageMetaFromPagination(pagination httprequest.Pagination, totalItems int64) PageMeta {
	return PageMeta{
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalItems: totalItems,
		TotalPages: TotalPages(totalItems, pagination.PageSize),
	}
}

func TotalPages(totalItems int64, pageSize int64) int64 {
	if pageSize < 1 {
		return 0
	}

	return (totalItems + pageSize - 1) / pageSize
}

type CursorMeta struct {
	NextCursor string `json:"nextCursor,omitempty"`
	HasMore    bool   `json:"hasMore"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

func JSON(w http.ResponseWriter, status int, body any) {
	data, err := json.Marshal(body)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, _ = w.Write(data)
}

func Error(
	w http.ResponseWriter,
	status int,
	message string,
) {
	if message == "" {
		message = http.StatusText(status)
	}

	JSON(w, status, ErrorResponse{
		Error: ErrorDetail{
			Message: message,
		},
	})
}
