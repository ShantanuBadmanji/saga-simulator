package datasource

import (
	"context"
	"database/sql"

	"github.com/ShantanuBadmanji/saga-simulator/backend/services/catalog/internal/datasource/sqlite/sqlcgen"
)

type CategoryRepository interface {
	GetCategoryByID(ctx context.Context, id string) (sqlcgen.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (sqlcgen.Category, error)
	ListCategories(ctx context.Context) ([]sqlcgen.Category, error)
	ListSubCategories(ctx context.Context, parentID sql.NullString) ([]sqlcgen.Category, error)
}

var _ CategoryRepository = (*sqlcgen.Queries)(nil)

type ProductRepository interface {
	ListProductsWithSkusAsc(ctx context.Context, arg sqlcgen.ListProductsWithSkusAscParams) ([]sqlcgen.ListProductsWithSkusAscRow, error)
	GetProductByID(ctx context.Context, id string) (sqlcgen.Product, error)
	ListProductMediaByProductID(ctx context.Context, productID string) ([]sqlcgen.ProductMedium, error)
}

var _ ProductRepository = (*sqlcgen.Queries)(nil)

type SkuRepository interface {
	GetSkuByID(ctx context.Context, id string) (sqlcgen.Sku, error)
	GetSkuBySkuCode(ctx context.Context, skuCode string) (sqlcgen.Sku, error)
	ListSkusByProductID(ctx context.Context, productID string) ([]sqlcgen.Sku, error)
}

var _ SkuRepository = (*sqlcgen.Queries)(nil)
