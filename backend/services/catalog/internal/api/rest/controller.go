package rest

import (
	"encoding/json"
	"net/http"
)

type CatalogController struct {
	catalogService *CatalogService
}

// NewCatalogController creates a new instance of CatalogController.
func NewCatalogController(catalogService *CatalogService) *CatalogController {
	return &CatalogController{
		catalogService: catalogService,
	}
}

func (c *CatalogController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /catalog", c.get)
}

// GetCatalog handles the request to retrieve the catalog.
func (c *CatalogController) get(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to retrieve the catalog data.
	// For example, you can fetch data from a database or any other source.

	// Example response (replace with actual catalog data):
	catalogData := []string{"Item 1", "Item 2", "Item 3"}

	// Convert the catalog data to JSON and write it to the response.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(catalogData)
}
