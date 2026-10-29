package rest

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ShantanuBadmanji/saga-simulator/backend/internal/platform/httprequest"
	"github.com/ShantanuBadmanji/saga-simulator/backend/internal/platform/httpresponse"
)

const (
	paramCategory = "category"
	paramParent   = "parent"
	pathID        = "id"
)

type CatalogHandler struct {
	catalogService *CatalogService
	log            *slog.Logger
}

func NewCatalogHandler(catalogService *CatalogService, logger *slog.Logger) *CatalogHandler {
	return &CatalogHandler{
		catalogService: catalogService,
		log:            logger,
	}
}

func (h *CatalogHandler) RegisterRoutes(mux *http.ServeMux) {

	mux.HandleFunc("GET /categories", h.listCategories)

	mux.HandleFunc("GET /products", h.listProducts)
	mux.HandleFunc("GET /products/{id}", h.getProductByID)

	mux.HandleFunc("GET /skus/{id}", h.getSkuByID)

}

func (h *CatalogHandler) listCategories(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	categories, err := h.catalogService.getCategories(r.Context(), query.Get(paramParent))
	if err != nil {
		h.log.ErrorContext(r.Context(), "failed to get categories", "error", err)
		httpresponse.Error(w, http.StatusInternalServerError, "")
		return
	}

	httpresponse.JSON(w,
		http.StatusOK,
		httpresponse.DataResponse{Data: categories},
	)
}

func (h *CatalogHandler) listProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	categoryID := query.Get(paramParent)
	pagination, err := httprequest.ParsePagination(r)
	if err != nil {
		httpresponse.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	productOverviews, totalItems, err := h.catalogService.getProducts(r.Context(), categoryID, pagination)
	if err != nil {
		h.log.ErrorContext(r.Context(), "failed to get products", "error", err)
		httpresponse.Error(w, http.StatusInternalServerError, "")
		return
	}

	httpresponse.JSON(w,
		http.StatusOK,
		httpresponse.DataResponse{
			Data: productOverviews,
			Meta: httpresponse.NewPageMetaFromPagination(pagination, totalItems),
		},
	)
}

func (h *CatalogHandler) getProductByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue(pathID)

	productDto, err := h.catalogService.getProductByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpresponse.Error(w, http.StatusNotFound, "")
			return
		}

		h.log.ErrorContext(r.Context(), "failed to get product", "id", id, "error", err)
		httpresponse.Error(w, http.StatusInternalServerError, "")
		return
	}

	httpresponse.JSON(
		w,
		http.StatusOK,
		httpresponse.DataResponse{
			Data: productDto,
		},
	)
}

func (h *CatalogHandler) getSkuByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue(pathID)

	sku, err := h.catalogService.getSkuByID(r.Context(), id)
	if err != nil {
		h.log.ErrorContext(r.Context(), "failed to get sku", "id", id, "error", err)
		httpresponse.Error(w, http.StatusInternalServerError, "")
		return
	}

	httpresponse.JSON(
		w,
		http.StatusOK,
		httpresponse.DataResponse{
			Data: sku,
		},
	)
}
