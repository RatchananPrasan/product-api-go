package domain

import "github.com/google/uuid"

// Product represents the product entity
type Product struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Description *string    `db:"description" json:"description"`
	Price       float64    `db:"price" json:"price"`
	SalePrice   *float64   `db:"sale_price" json:"sale_price"`
}

// CreateProductInput holds data for creating a product
type CreateProductInput struct {
	Name        string   `json:"name" binding:"required"`
	Description *string  `json:"description"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	SalePrice   *float64 `json:"sale_price"`
}

// PatchProductInput holds data for patching a product (all fields optional)
type PatchProductInput struct {
	Name        *string  `json:"name"`
	Description **string `json:"description"`
	Price       *float64 `json:"price"`
	SalePrice   **float64 `json:"sale_price"`
}

// ProductRepository defines the data access interface
type ProductRepository interface {
	Create(product *Product) (*Product, error)
	GetByID(id uuid.UUID) (*Product, error)
	Patch(id uuid.UUID, input *PatchProductInput) error
}

// ProductUsecase defines the business logic interface
type ProductUsecase interface {
	CreateProduct(input *CreateProductInput) (*Product, error)
	PatchProduct(id uuid.UUID, input *PatchProductInput) error
}
