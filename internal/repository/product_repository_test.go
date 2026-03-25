package repository_test

// Repository integration tests require a real PostgreSQL connection.
// Run with: DB_HOST=localhost DB_PORT=5432 ... go test ./internal/repository/...
//
// These tests are tagged as integration and skipped unless the DB env vars are set.

import (
	"os"
	"testing"

	"github.com/example/product-api/internal/domain"
	"github.com/example/product-api/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skip("TEST_DB_DSN not set, skipping integration tests")
	}
	db, err := sqlx.Connect("postgres", dsn)
	require.NoError(t, err)

	// Clean up before test
	db.MustExec("DELETE FROM products")
	return db
}

func ptr[T any](v T) *T { return &v }

func TestProductRepository_Create_Integration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewProductRepository(db)
	product := &domain.Product{
		ID:    uuid.New(),
		Name:  "Integration Widget",
		Price: 29.99,
	}

	created, err := repo.Create(product)

	require.NoError(t, err)
	assert.Equal(t, product.ID, created.ID)
	assert.Equal(t, "Integration Widget", created.Name)
}

func TestProductRepository_GetByID_NotFound_Integration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewProductRepository(db)

	_, err := repo.GetByID(uuid.New())

	assert.ErrorIs(t, err, domain.ErrProductNotFound)
}

func TestProductRepository_Patch_Integration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewProductRepository(db)
	id := uuid.New()

	_, err := repo.Create(&domain.Product{ID: id, Name: "Before", Price: 100.0})
	require.NoError(t, err)

	err = repo.Patch(id, &domain.PatchProductInput{Name: ptr("After")})
	require.NoError(t, err)

	updated, err := repo.GetByID(id)
	require.NoError(t, err)
	assert.Equal(t, "After", updated.Name)
}

func TestProductRepository_Patch_NotFound_Integration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewProductRepository(db)
	err := repo.Patch(uuid.New(), &domain.PatchProductInput{Name: ptr("X")})
	assert.ErrorIs(t, err, domain.ErrProductNotFound)
}
