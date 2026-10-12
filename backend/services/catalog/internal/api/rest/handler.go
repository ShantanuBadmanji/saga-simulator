package rest

import (
	"encoding/json"
	"net/http"
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

func (c *CatalogHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /catalog", c.get)
}

// GetCatalog handles the request to retrieve the catalog.
func (c *CatalogHandler) get(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to retrieve the catalog data.
	// For example, you can fetch data from a database or any other source.

	// Example response (replace with actual catalog data):
	catalogData := []string{"Item 1", "Item 2", "Item 3"}

	// Convert the catalog data to JSON and write it to the response.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(catalogData)
}
