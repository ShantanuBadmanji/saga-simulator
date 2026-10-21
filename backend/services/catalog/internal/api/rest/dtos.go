package rest

import "github.com/ShantanuBadmanji/saga-simulator/backend/services/catalog/internal/datasource/sqlite/sqlcgen"

type ProductDTO struct {
	Product *sqlcgen.Product
	Media   []sqlcgen.ProductMedium
	Skus    []sqlcgen.Sku
}

type ProductOverview struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	CategoryID  string `json:"category_id"`
	Brand       string `json:"brand"`
	Status      string `json:"status"`
	Price       int64  `json:"price"`
}

func NewProductOverviewFromRow(r *sqlcgen.ListProductsWithSkusAscRow) *ProductOverview {
	return &ProductOverview{
		ID:          r.ID,
		Name:        r.Name,
		Slug:        r.Slug,
		Description: r.Description.String,
		CategoryID:  r.CategoryID.String,
		Brand:       r.Brand.String,
		Status:      r.Status,
		Price:       r.MinPrice,
	}
}
