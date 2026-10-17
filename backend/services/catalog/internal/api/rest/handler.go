package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type CatalogHandler struct {
	catalogService *CatalogService
}

// NewCatalogHandler creates a new instance of CatalogController.
func NewCatalogHandler(catalogService *CatalogService) *CatalogHandler {
	return &CatalogHandler{
		catalogService: catalogService,
	}
}

func (h *CatalogHandler) RegisterRoutes(mux *http.ServeMux) {

	mux.HandleFunc("GET /categories", h.listCategories)

	mux.HandleFunc("GET /products", h.listProducts)
	mux.HandleFunc("GET /products/{id}", h.getProductByID)

	mux.HandleFunc("GET /skus/{id}", h.getSkuByID)

}

// GetCatalog handles the request to retrieve the catalog.
func (h *CatalogHandler) listCategories(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	catagories, err := h.catalogService.GetCategories(r.Context(), query.Get("parent"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(catagories)
	}

	w.Header().Set("Content-Type", "application/json")
}

func (h *CatalogHandler) listProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	categoryID := query.Get("category")
	limit, err := strconv.ParseInt(query.Get("limit"), 10, 64)
	if err != nil {
		limit = 20
	}
	offset, err := strconv.ParseInt(query.Get("offset"), 10, 64)
	if err != nil {
		offset = 0
	}

	productsWithSkus, err := h.catalogService.GetProducts(r.Context(), categoryID, limit, offset)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(productsWithSkus)
	}

	w.Header().Set("Content-Type", "application/json")
}

func (h *CatalogHandler) getProductByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	productDto, err := h.catalogService.GetProductByID(r.Context(), id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(productDto)
	}
}

func (h *CatalogHandler) getSkuByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	sku, err := h.catalogService.getSkuByID(r.Context(), id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(sku)
	}
}
