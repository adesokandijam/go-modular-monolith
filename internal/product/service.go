package product

import (
	"context"
	"dijam-ecommerce/shared"
	"fmt"
)

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

// var filter shared.ProductFilters

func (p *ProductService) List(ctx context.Context, category string, filter shared.ProductFilters) ([]*Product, ProductMetadata, error) {
	products, metadata, err := p.repo.List(ctx, category, filter)
	if err != nil {
		return nil, metadata, fmt.Errorf("error listing product: %w", err)
	}
	return products, metadata, nil
}
