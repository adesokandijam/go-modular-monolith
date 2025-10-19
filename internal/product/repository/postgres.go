package repository

import (
	"context"
	"database/sql"
	"dijam-ecommerce/internal/product"
	"dijam-ecommerce/shared"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ProductPostgresRepository struct {
	DB *sql.DB
}

var _ product.ProductRepository = (*ProductPostgresRepository)(nil)

// NewProductRepository is the constructor.
func NewProductRepository(db *sql.DB) *ProductPostgresRepository {
	return &ProductPostgresRepository{
		DB: db,
	}
}

func calculateProductMetadata(totalRecords, page, pageSize int64) product.ProductMetadata {
	if totalRecords == 0 {
		return product.ProductMetadata{}
	}

	return product.ProductMetadata{
		CurrentPage:   page,
		PageSize:      pageSize,
		FirstPage:     1,
		LastPage:      (totalRecords + pageSize - 1) / pageSize,
		TotalProducts: totalRecords,
	}
}

// FindByID retrieves a single product from the database.
func (r *ProductPostgresRepository) FindByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	query := `
		SELECT id, name, description, price_in_pence, currency, stock_quantity, created_at, updated_at
		FROM products
		WHERE id = $1`

	p := &product.Product{}
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.PriceInPence,
		&p.Currency,
		&p.StockQuantity,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("postgres: could not scan product: %w", err)
	}
	return p, nil
}

func (r *ProductPostgresRepository) List(ctx context.Context, category string, filters shared.ProductFilters) ([]*product.Product, product.ProductMetadata, error) {
	var totalRecords int64
	countQuery := `SELECT count(*) FROM products WHERE (LOWER(category) = LOWER($1) OR $1 = '')`
	err := r.DB.QueryRowContext(ctx, countQuery, category).Scan(&totalRecords)
	if err != nil {
		return nil, product.ProductMetadata{}, fmt.Errorf("postgres: could not count products: %w", err)
	}
	query := fmt.Sprintf(
		`SELECT id, name, description, price_in_pence, currency, stock_quantity, category, created_at, updated_at
         FROM products 
         WHERE (LOWER(category) = LOWER($1) OR $1 = '')
         ORDER BY %s %s, id ASC
         LIMIT $2
         OFFSET $3`,
		filters.SortColumns(),
		filters.SortDirection(),
	)
	start := time.Now()
	rows, err := r.DB.QueryContext(ctx, query, category, filters.Limit(), filters.Offset())
	if err != nil {
		return nil, product.ProductMetadata{}, fmt.Errorf("postgres: could not query products: %w", err)
	}
	defer rows.Close()
	fmt.Println(time.Since(start))
	// fmt.Println(start)

	var productList = make([]*product.Product, 0)
	for rows.Next() {
		p := &product.Product{}
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.PriceInPence,
			&p.Currency,
			&p.StockQuantity,
			&p.Category,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, product.ProductMetadata{}, fmt.Errorf("postgres: could not scan product row: %w", err)
		}
		productList = append(productList, p)
	}
	metadata := calculateProductMetadata(int64(totalRecords), filters.Page, filters.PageSize)
	return productList, metadata, nil
}

func (r *ProductPostgresRepository) Save(ctx context.Context, product *product.Product) error {
	return fmt.Errorf("not implemented")
}

func (r *ProductPostgresRepository) UpdateStock(ctx context.Context, id uuid.UUID, newQuantity int) error {
	return fmt.Errorf("not implemented")
}
