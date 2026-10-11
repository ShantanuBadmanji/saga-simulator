package rest

import "github.com/ShantanuBadmanji/saga-simulator/backend/services/catalog/internal/datasource/repository"

type CatalogService struct {
	productRepository repository.ProductRepository
}

func NewCatalogService(productRepository repository.ProductRepository) *CatalogService {
	return &CatalogService{
		productRepository: productRepository,
	}
}

func (c *CatalogService) Get() {

}
