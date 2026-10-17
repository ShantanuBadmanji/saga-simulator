package rest

import "github.com/ShantanuBadmanji/saga-simulator/backend/services/catalog/internal/datasource/sqlite/sqlcgen"

type ProductDTO struct {
	Product *sqlcgen.Product
	Media   []sqlcgen.ProductMedium
	Skus    []sqlcgen.Sku
}
