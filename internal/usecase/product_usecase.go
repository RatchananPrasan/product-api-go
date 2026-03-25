package usecase

import (
	"github.com/example/product-api/internal/domain"
	"github.com/google/uuid"
)

type productUsecase struct {
	productRepo domain.ProductRepository
}

// NewProductUsecase creates a new ProductUsecase with dependency injection
func NewProductUsecase(repo domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{productRepo: repo}
}

// CreateProduct validates and creates a new product
func (u *productUsecase) CreateProduct(input *domain.CreateProductInput) (*domain.Product, error) {
	if input.SalePrice != nil && *input.SalePrice >= input.Price {
		return nil, domain.ErrSalePriceMustBeLessThanPrice
	}

	product := &domain.Product{
		ID:          uuid.New(),
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		SalePrice:   input.SalePrice,
	}

	return u.productRepo.Create(product)
}

// PatchProduct applies partial updates to a product
func (u *productUsecase) PatchProduct(id uuid.UUID, input *domain.PatchProductInput) error {
	existing, err := u.productRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Validate sale_price vs price consistency
	finalPrice := existing.Price
	if input.Price != nil {
		finalPrice = *input.Price
	}

	if input.SalePrice != nil && *input.SalePrice != nil {
		if **input.SalePrice >= finalPrice {
			return domain.ErrSalePriceMustBeLessThanPrice
		}
	}

	return u.productRepo.Patch(id, input)
}
