package component_test

// Component (E2E within service) tests
// These wire up the full stack (handler -> usecase -> repo) against a real test DB.
// Run with: TEST_DB_DSN=... go test ./internal/component/...

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/example/product-api/internal/handler"
	"github.com/example/product-api/internal/repository"
	"github.com/example/product-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupApp(t *testing.T) *gin.Engine {
	t.Helper()
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skip("TEST_DB_DSN not set, skipping component tests")
	}

	db, err := sqlx.Connect("postgres", dsn)
	require.NoError(t, err)
	db.MustExec("DELETE FROM products")

	repo := repository.NewProductRepository(db)
	uc := usecase.NewProductUsecase(repo)
	h := handler.NewProductHandler(uc)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h.RegisterRoutes(r)
	return r
}

func TestCreateAndPatchProduct_E2E(t *testing.T) {
	r := setupApp(t)

	// Step 1: Create product
	createBody := `{"name":"E2E Widget","price":200.0,"sale_price":150.0}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/product", bytes.NewBufferString(createBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	assert.True(t, createResp["successful"].(bool))
	data := createResp["data"].(map[string]interface{})
	productID := data["data1"].(string)

	// Step 2: Patch the product
	patchBody := `{"name":"Updated E2E Widget"}`
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/product/%s", productID), bytes.NewBufferString(patchBody))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var patchResp map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &patchResp)
	assert.True(t, patchResp["successful"].(bool))
}
