package rest

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ShantanuBadmanji/saga-simulator/backend/services/catalog/internal/datasource"
	"github.com/ShantanuBadmanji/saga-simulator/backend/services/catalog/internal/datasource/sqlite/sqlcgen"
)

type CatalogService struct {
	categoryRepo datasource.CategoryRepository
	productRepo  datasource.ProductRepository
	skuRepo      datasource.SkuRepository
}

func NewCatalogService(categoryRepo datasource.CategoryRepository, productRepo datasource.ProductRepository, skuRepo datasource.SkuRepository) *CatalogService {
	return &CatalogService{categoryRepo, productRepo, skuRepo}
}

func (s *CatalogService) GetCategories(ctx context.Context, parentID string) ([]sqlcgen.Category, error) {
	if len(parentID) != 0 {
		return s.categoryRepo.ListSubCategories(ctx, sql.NullString{String: parentID, Valid: true})
	}

	return s.categoryRepo.ListCategories(ctx)
}

func (s *CatalogService) GetProducts(ctx context.Context, categoryID string, limit, offset int64) ([]sqlcgen.ListProductsWithSkusAscRow, error) {

	var nullableCategoryID sql.NullString
	if len(categoryID) != 0 {
		nullableCategoryID = sql.NullString{String: categoryID, Valid: true}
	} else {
		nullableCategoryID = sql.NullString{Valid: false}
	}

	return s.productRepo.ListProductsWithSkusAsc(ctx, sqlcgen.ListProductsWithSkusAscParams{CategoryID: nullableCategoryID, Limit: limit, Offset: offset})
}

func (s *CatalogService) GetProductByID(ctx context.Context, id string) (*ProductDTO, error) {
	product, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	media, err1 := s.productRepo.ListProductMediaByProductID(ctx, id)
	skus, err2 := s.skuRepo.ListSkusByProductID(ctx, id)
	if err := errors.Join(err1, err2); err != nil {
		return nil, err
	}

	return &ProductDTO{&product, media, skus}, nil
}

func (s *CatalogService) getSkuByID(ctx context.Context, id string) (*sqlcgen.Sku, error) {
	sku, err := s.skuRepo.GetSkuByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &sku, nil
}
