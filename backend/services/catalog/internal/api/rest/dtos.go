package rest

import "github.com/ShantanuBadmanji/saga-simulator/backend/services/catalog/internal/datasource/sqlite/sqlcgen"

type ProductDTO struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description string         `json:"description"`
	CategoryID  string         `json:"category_id"`
	Brand       string         `json:"brand"`
	Status      string         `json:"status"`
	Media       []ProductMedia `json:"media"`
	Skus        []SkuDTO       `json:"skus"`
}

type ProductMedia struct {
	ID        string `json:"id"`
	Kind      string `json:"kind"`
	URL       string `json:"url"`
	SortOrder int64  `json:"sort_order"`
}

type SkuDTO struct {
	ID                string `json:"id"`
	SkuCode           string `json:"sku_code"`
	Attributes        string `json:"attributes"`
	PriceMinor        int64  `json:"price_minor"`
	Currency          string `json:"currency"`
	Active            bool   `json:"active"` // convert int64 → bool here
	AvailableQuantity int64  `json:"available_quantity"`
}

func NewProductDTO(p *sqlcgen.Product, media []sqlcgen.ProductMedium, skus []sqlcgen.Sku) *ProductDTO {
	prodMedia := make([]ProductMedia, len(media))
	for i, m := range media {
		prodMedia[i] = ProductMedia{
			ID:        m.ID,
			Kind:      m.Kind,
			URL:       m.Url,
			SortOrder: m.SortOrder,
		}
	}
	skuDtos := make([]SkuDTO, len(skus))
	for i, s := range skus {
		skuDtos[i] = SkuDTO{
			ID:         s.ID,
			SkuCode:    s.SkuCode,
			Attributes: s.Attributes,
			PriceMinor: s.PriceMinor,
			Currency:   s.Currency,
			Active:     s.Active != 0,
		}
	}

	return &ProductDTO{
		ID:          p.ID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description.String,
		CategoryID:  p.CategoryID.String,
		Brand:       p.Brand.String,
		Status:      p.Status,
		Media:       prodMedia,
		Skus:        skuDtos,
	}
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
