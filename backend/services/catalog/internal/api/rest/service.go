package rest

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/ShantanuBadmanji/saga-simulator/backend/internal/platform/httprequest"
	"github.com/ShantanuBadmanji/saga-simulator/backend/services/catalog/internal/datasource"
	"github.com/ShantanuBadmanji/saga-simulator/backend/services/catalog/internal/datasource/sqlite/sqlcgen"
)

type CatalogService struct {
	categoryRepo datasource.CategoryRepository
	productRepo  datasource.ProductRepository
	skuRepo      datasource.SkuRepository
	log          *slog.Logger
}

func NewCatalogService(categoryRepo datasource.CategoryRepository, productRepo datasource.ProductRepository, skuRepo datasource.SkuRepository, logger *slog.Logger) *CatalogService {
	return &CatalogService{categoryRepo, productRepo, skuRepo, logger}
}

func (s *CatalogService) getCategories(ctx context.Context, parentID string) ([]sqlcgen.Category, error) {
	if len(parentID) != 0 {
		return s.categoryRepo.ListSubCategories(ctx, sql.NullString{String: parentID, Valid: true})
	}

	return s.categoryRepo.ListCategories(ctx)
}

func (s *CatalogService) getProducts(ctx context.Context, categoryID string, pagination httprequest.Pagination) ([]ProductOverview, int64, error) {

	var nullableCategoryID sql.NullString
	if len(categoryID) != 0 {
		nullableCategoryID = sql.NullString{String: categoryID, Valid: true}
	} else {
		nullableCategoryID = sql.NullString{Valid: false}
	}

	limit, offset := pagination.LimitOffset()
	productRows, err := s.productRepo.ListProductsWithSkusAsc(ctx, sqlcgen.ListProductsWithSkusAscParams{
		CategoryID: nullableCategoryID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, 0, err
	}

	var totalItems int64
	if len(productRows) == 0 {
		totalItems = 0
	} else {
		totalItems = productRows[0].TotalCount
	}

	productOverviews := make([]ProductOverview, len(productRows))
	for i, r := range productRows {
		productOverviews[i] = *NewProductOverviewFromRow(&r)
	}

	return productOverviews, totalItems, err
}

func (s *CatalogService) getProductByID(ctx context.Context, id string) (*ProductDTO, error) {
	product, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	media, err1 := s.productRepo.ListProductMediaByProductID(ctx, id)
	skus, err2 := s.skuRepo.ListSkusByProductID(ctx, id)
	if err := errors.Join(err1, err2); err != nil {
		return nil, err
	}

	return NewProductDTO(&product, media, skus), nil
}

func (s *CatalogService) getSkuByID(ctx context.Context, id string) (*sqlcgen.Sku, error) {
	sku, err := s.skuRepo.GetSkuByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &sku, nil
}
