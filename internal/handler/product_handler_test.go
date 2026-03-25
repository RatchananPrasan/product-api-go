package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/product-api/internal/domain"
	"github.com/example/product-api/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductUsecase mocks domain.ProductUsecase
type MockProductUsecase struct {
	mock.Mock
}

func (m *MockProductUsecase) CreateProduct(input *domain.CreateProductInput) (*domain.Product, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductUsecase) PatchProduct(id uuid.UUID, input *domain.PatchProductInput) error {
	args := m.Called(id, input)
	return args.Error(0)
}

func setupRouter(uc domain.ProductUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := handler.NewProductHandler(uc)
	h.RegisterRoutes(r)
	return r
}

func ptr[T any](v T) *T { return &v }

// --- POST /product ---

func TestCreateProduct_Handler_Success(t *testing.T) {
	mockUC := new(MockProductUsecase)
	r := setupRouter(mockUC)

	id := uuid.New()
	mockUC.On("CreateProduct", mock.Anything).
		Return(&domain.Product{ID: id, Name: "Widget", Price: 100.0}, nil)

	body := `{"name":"Widget","price":100.0}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/product", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["successful"].(bool))
}

func TestCreateProduct_Handler_InvalidBody(t *testing.T) {
	mockUC := new(MockProductUsecase)
	r := setupRouter(mockUC)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/product", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateProduct_Handler_SalePriceError(t *testing.T) {
	mockUC := new(MockProductUsecase)
	r := setupRouter(mockUC)

	mockUC.On("CreateProduct", mock.Anything).
		Return(nil, domain.ErrSalePriceMustBeLessThanPrice)

	body := `{"name":"Widget","price":10.0,"sale_price":99.0}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/product", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "SALE_PRICE_MUST_BE_LESS_THAN_PRICE", resp["error_code"])
}

// --- PATCH /product/:id ---

func TestPatchProduct_Handler_Success(t *testing.T) {
	mockUC := new(MockProductUsecase)
	r := setupRouter(mockUC)

	id := uuid.New()
	mockUC.On("PatchProduct", id, mock.Anything).Return(nil)

	body := `{"name":"Updated"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/product/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPatchProduct_Handler_InvalidID(t *testing.T) {
	mockUC := new(MockProductUsecase)
	r := setupRouter(mockUC)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/product/not-a-uuid", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPatchProduct_Handler_NotFound(t *testing.T) {
	mockUC := new(MockProductUsecase)
	r := setupRouter(mockUC)

	id := uuid.New()
	mockUC.On("PatchProduct", id, mock.Anything).Return(domain.ErrProductNotFound)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/product/"+id.String(), bytes.NewBufferString(`{"name":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
