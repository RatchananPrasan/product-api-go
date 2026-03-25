package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/example/product-api/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type productRepository struct {
	db *sqlx.DB
}

// NewProductRepository creates a new ProductRepository
func NewProductRepository(db *sqlx.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *domain.Product) (*domain.Product, error) {
	query := `
		INSERT INTO products (id, name, description, price, sale_price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, description, price, sale_price`

	var result domain.Product
	err := r.db.QueryRowx(query,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.SalePrice,
	).StructScan(&result)
	if err != nil {
		return nil, fmt.Errorf("repository.Create: %w", err)
	}
	return &result, nil
}

func (r *productRepository) GetByID(id uuid.UUID) (*domain.Product, error) {
	query := `SELECT id, name, description, price, sale_price FROM products WHERE id = $1`

	var product domain.Product
	err := r.db.QueryRowx(query, id).StructScan(&product)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrProductNotFound
		}
		return nil, fmt.Errorf("repository.GetByID: %w", err)
	}
	return &product, nil
}

func (r *productRepository) Patch(id uuid.UUID, input *domain.PatchProductInput) error {
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if input.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, *input.Name)
		argIdx++
	}
	if input.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, *input.Description) // *input.Description is *string (nullable)
		argIdx++
	}
	if input.Price != nil {
		setClauses = append(setClauses, fmt.Sprintf("price = $%d", argIdx))
		args = append(args, *input.Price)
		argIdx++
	}
	if input.SalePrice != nil {
		setClauses = append(setClauses, fmt.Sprintf("sale_price = $%d", argIdx))
		args = append(args, *input.SalePrice) // *input.SalePrice is *float64 (nullable)
		argIdx++
	}

	if len(setClauses) == 0 {
		return nil // nothing to update
	}

	args = append(args, id)
	query := fmt.Sprintf(
		"UPDATE products SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "),
		argIdx,
	)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("repository.Patch: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}
