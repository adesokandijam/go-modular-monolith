package product

import (
	"context"
	"dijam-ecommerce/shared"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	PriceInPence  int64     `json:"price_in_pence"`
	Currency      string    `json:"currency"`
	StockQuantity int64     `json:"stock_quantity"`
	Category      string    `json:"category"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ProductMetadata struct {
	FirstPage     int64 `json:"first_page,omitempty"`
	PageSize      int64 `json:"page_size,omitempty"`
	TotalProducts int64 `json:"total_products,omitempty"`
	CurrentPage   int64 `json:"current_page,omitempty"`
	LastPage      int64 `json:"last_page,omitempty"`
}

// ProductRepository defines the contract for product persistence.
type ProductRepository interface {
	Save(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, id uuid.UUID) (*Product, error)
	List(ctx context.Context, category string, filters shared.ProductFilters) ([]*Product, ProductMetadata, error)
	UpdateStock(ctx context.Context, id uuid.UUID, newQuantity int) error
}
