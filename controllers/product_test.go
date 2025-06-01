package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/go-example/models"
	"github.com/sing3demons/go-example/store"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	pathProducts = "/api/products"
)

func setupProductApp(db store.Storer) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	productController := NewProductController(db)

	r := gin.New()
	r.GET(pathProducts, productController.Find)
	r.GET(pathProducts+"/:id", productController.FindOne)
	r.POST(pathProducts, productController.Create)

	req, _ := http.NewRequest(http.MethodGet, pathProducts, nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	return rec
}

func TestFindProducts(t *testing.T) {
	t.Run("Find Products success", func(t *testing.T) {
		products := []models.Product{
			{
				ID:          primitive.NewObjectID(),
				Name:        "Product 1",
				Price:       99,
				Description: "Description for Product 1",
			},
		}

		db := store.MockStore{
			Data: products,
		}
		rec := setupProductApp(&db)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", rec.Code)
		}

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Len(t, response["data"].([]interface{}), len(products), "Expected data length to match products length")
	})

	t.Run("Find Products empty", func(t *testing.T) {
		db := store.MockStore{
			Data: []models.Product{},
			Err:  mongo.ErrNoDocuments,
		}
		rec := setupProductApp(&db)

		assert.Equal(t, http.StatusNotFound, rec.Code, "Expected status code 200")

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err, "failed to unmarshal response")

		assert.Empty(t, response["data"], "Expected empty data in response")
	})

	t.Run("Find Products error", func(t *testing.T) {
		db := store.MockStore{
			Err: mongo.ErrClientDisconnected,
		}
		rec := setupProductApp(&db)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code 500, got %d", rec.Code)
		}

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err, "failed to unmarshal response")

		assert.Equal(t, db.Err.Error(), response["message"], "Expected error message to match")
	})
}

func setupProductPost(db store.Storer, body *strings.Reader) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	productController := NewProductController(db)

	r := gin.New()
	r.POST(pathProducts, productController.Create)

	req, _ := http.NewRequest(http.MethodPost, pathProducts, body)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	return rec
}
func setupProductGetByID(db store.Storer, id string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	productController := NewProductController(db)

	r := gin.New()
	r.GET(pathProducts+"/:id", productController.FindOne)

	req, _ := http.NewRequest(http.MethodGet, pathProducts+"/"+id, nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	return rec
}

func TestFindOneProduct(t *testing.T) {
	t.Run("Find One Product success", func(t *testing.T) {
		productID := primitive.NewObjectID()
		product := models.Product{
			ID:          productID,
			Name:        "Product 1",
			Price:       99,
			Description: "Description for Product 1",
		}

		db := store.MockStore{
			Data: []models.Product{product},
		}
		rec := setupProductGetByID(&db, productID.Hex())

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200")

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err, "Failed to unmarshal response")

		assert.NotNil(t, response["data"], "expected 'data' key in response")

	})

	t.Run("Find One Product Invalid ID format", func(t *testing.T) {
		db := store.MockStore{}
		rec := setupProductGetByID(&db, "invalid-id")

		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code 400")

		var response map[string]string
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal err: %v", err)
		}

		if response["message"] != "Invalid ID format" {
			t.Errorf("Expected error message 'Invalid ID format', got '%s'", response["message"])
		}
	})

	t.Run("Find One Product not found", func(t *testing.T) {
		db := store.MockStore{
			Err: mongo.ErrNoDocuments,
		}
		rec := setupProductGetByID(&db, primitive.NewObjectID().Hex())

		assert.Equal(t, http.StatusNotFound, rec.Code, "Expected status code 404")

		var response map[string]interface{}
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if _, ok := response["message"]; !ok {
			t.Error("Expected 'data' key in response")
		}
	})

	t.Run("Find One Product error", func(t *testing.T) {
		db := store.MockStore{
			Err: mongo.ErrClientDisconnected,
		}
		rec := setupProductGetByID(&db, primitive.NewObjectID().Hex())

		assert.Equal(t, http.StatusInternalServerError, rec.Code, "expected status code 500")

		var response map[string]string
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["message"] != db.Err.Error() {
			t.Errorf("Expected error message '%s', got '%s'", db.Err.Error(), response["message"])
		}
	})
}

func TestCreateProduct(t *testing.T) {
	t.Run("Create Product success", func(t *testing.T) {
		product := models.Product{
			Name:        "New Product",
			Price:       100,
			Description: "Description for New Product",
		}

		body, _ := json.Marshal(product)
		db := store.MockStore{}
		rec := setupProductPost(&db, strings.NewReader(string(body)))

		assert.Equal(t, http.StatusCreated, rec.Code, "Expected status code 201")

		var response map[string]interface{}
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if _, ok := response["data"]; !ok {
			t.Error("Expected 'data' key in response")
		}
	})

	t.Run("Create Product error bad request", func(t *testing.T) {
		body := strings.NewReader(`{"name": "Invalid Product"}`) // Missing required fields
		db := store.MockStore{
			Err: mongo.ErrClientDisconnected,
		}
		rec := setupProductPost(&db, body)

		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code 500")

		var response map[string]string
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			assert.NoError(t, err, "Failed to unmarshal response")
		}

		assert.Equal(t, "Key: 'Product.Price' Error:Field validation for 'Price' failed on the 'required' tag\nKey: 'Product.Description' Error:Field validation for 'Description' failed on the 'required' tag", response["message"], "Expected error message 'insert error'")
	})

	t.Run("Create Product error", func(t *testing.T) {
		product := models.Product{
			Name:        "New Product",
			Price:       100,
			Description: "Description for New Product",
		}

		body, _ := json.Marshal(product)
		db := store.MockStore{
			Err: mongo.ErrClientDisconnected,
		}
		rec := setupProductPost(&db, strings.NewReader(string(body)))

		assert.Equal(t, http.StatusInternalServerError, rec.Code, "Expected status code 500")

		var response map[string]string
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response["message"] != db.Err.Error() {
			t.Errorf("Expected error message '%s', got '%s'", db.Err.Error(), response["message"])
		}
	})
}
