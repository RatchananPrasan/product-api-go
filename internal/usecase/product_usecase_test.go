package usecase_test

import (
	"errors"
	"testing"

	"github.com/example/product-api/internal/domain"
	"github.com/example/product-api/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository is a mock for domain.ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *domain.Product) (*domain.Product, error) {
	args := m.Called(product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) GetByID(id uuid.UUID) (*domain.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) Patch(id uuid.UUID, input *domain.PatchProductInput) error {
	args := m.Called(id, input)
	return args.Error(0)
}

func ptr[T any](v T) *T { return &v }

// --- CreateProduct Tests ---

func TestCreateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := usecase.NewProductUsecase(mockRepo)

	input := &domain.CreateProductInput{
		Name:  "Widget",
		Price: 100.0,
	}

	mockRepo.On("Create", mock.AnythingOfType("*domain.Product")).
		Return(&domain.Product{ID: uuid.New(), Name: "Widget", Price: 100.0}, nil)

	product, err := uc.CreateProduct(input)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "Widget", product.Name)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_SalePriceGreaterThanPrice_ReturnsError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := usecase.NewProductUsecase(mockRepo)

	input := &domain.CreateProductInput{
		Name:      "Widget",
		Price:     50.0,
		SalePrice: ptr(99.0),
	}

	_, err := uc.CreateProduct(input)

	assert.ErrorIs(t, err, domain.ErrSalePriceMustBeLessThanPrice)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateProduct_SalePriceEqualToPrice_ReturnsError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := usecase.NewProductUsecase(mockRepo)

	input := &domain.CreateProductInput{
		Name:      "Widget",
		Price:     50.0,
		SalePrice: ptr(50.0),
	}

	_, err := uc.CreateProduct(input)

	assert.ErrorIs(t, err, domain.ErrSalePriceMustBeLessThanPrice)
}

func TestCreateProduct_RepoError_PropagatesError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := usecase.NewProductUsecase(mockRepo)

	input := &domain.CreateProductInput{Name: "Widget", Price: 100.0}
	mockRepo.On("Create", mock.Anything).Return(nil, errors.New("db error"))

	_, err := uc.CreateProduct(input)

	assert.Error(t, err)
}

// --- PatchProduct Tests ---

func TestPatchProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := usecase.NewProductUsecase(mockRepo)

	id := uuid.New()
	existing := &domain.Product{ID: id, Name: "Old", Price: 100.0}
	input := &domain.PatchProductInput{Name: ptr("New Name")}

	mockRepo.On("GetByID", id).Return(existing, nil)
	mockRepo.On("Patch", id, input).Return(nil)

	err := uc.PatchProduct(id, input)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPatchProduct_ProductNotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := usecase.NewProductUsecase(mockRepo)

	id := uuid.New()
	mockRepo.On("GetByID", id).Return(nil, domain.ErrProductNotFound)

	err := uc.PatchProduct(id, &domain.PatchProductInput{})

	assert.ErrorIs(t, err, domain.ErrProductNotFound)
}

func TestPatchProduct_NewSalePriceExceedsExistingPrice_ReturnsError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := usecase.NewProductUsecase(mockRepo)

	id := uuid.New()
	existing := &domain.Product{ID: id, Price: 50.0}
	salePtr := ptr(99.0)
	input := &domain.PatchProductInput{SalePrice: &salePtr}

	mockRepo.On("GetByID", id).Return(existing, nil)

	err := uc.PatchProduct(id, input)

	assert.ErrorIs(t, err, domain.ErrSalePriceMustBeLessThanPrice)
	mockRepo.AssertNotCalled(t, "Patch")
}
