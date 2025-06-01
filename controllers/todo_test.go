package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/go-example/models"
	"github.com/sing3demons/go-example/store"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type Todo struct {
	ID        int       `json:"ID"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	DeletedAt any       `json:"DeletedAt"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
}

func setupApp(db store.Storer) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	todoController := NewTodoController(db)

	r := gin.New()
	r.GET("/api/todos", todoController.Index)

	req, _ := http.NewRequest(http.MethodGet, "/api/todos", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	return rec
}

func setupPost(db store.Storer, body *strings.Reader) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	todoController := NewTodoController(db)

	r := gin.New()
	r.POST("/api/todos", todoController.Create)

	req, _ := http.NewRequest(http.MethodPost, "/api/todos", body)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	return rec
}

func TestFindTodo(t *testing.T) {
	mockTodos := []models.Todo{
		{
			Model:     gorm.Model{ID: 1},
			Title:     "Buy groceries",
			Completed: false,
		},
		{
			Model:     gorm.Model{ID: 2},
			Title:     "Read a book",
			Completed: true,
		},
	}

	t.Run("Find Todo success", func(t *testing.T) {
		rec := setupApp(&store.MockStore{
			Data: mockTodos,
		})

		assert.Equal(t, http.StatusOK, rec.Code)
		var response struct {
			Data []models.Todo `json:"data"`
		}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, len(mockTodos), len(response.Data))

	})

	t.Run("Find Todo data not found", func(t *testing.T) {
		rec := setupApp(&store.MockStore{
			Data: []models.Todo{},
			Err:  gorm.ErrRecordNotFound,
		})

		assert.Equal(t, http.StatusNotFound, rec.Code)
		var response struct {
			Message string `json:"message"`
		}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "No todos found", response.Message)
	})

	t.Run("Find Todo error", func(t *testing.T) {
		rec := setupApp(&store.MockStore{
			Err: gorm.ErrInvalidData,
		})

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		var response struct {
			Message string `json:"message"`
		}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, gorm.ErrInvalidData.Error(), response.Message)
	})
}

func TestCreateTodo(t *testing.T) {
	t.Run("Create Todo success", func(t *testing.T) {
		userInput := models.Todo{
			Title:     "Buy groceries",
			Completed: false,
		}

		userInputJSON, err := json.Marshal(userInput)
		assert.NoError(t, err)

		rec := setupPost(&store.MockStore{
			Data: userInput,
			Err:  nil,
		}, strings.NewReader(string(userInputJSON)))

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.NotEmpty(t, rec.Body.String())
	})

	t.Run("Create Todo validation error", func(t *testing.T) {
		rec := setupPost(&store.MockStore{
			Data: nil,
			Err:  gorm.ErrInvalidData,
		}, strings.NewReader(`{"title":""}`))

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var response struct {
			Message string `json:"message"`
		}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Key: 'TodoCreateRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag", response.Message)
	})

	t.Run("Create Todo error", func(t *testing.T) {
		userInput := models.Todo{
			Title:     "Buy groceries",
			Completed: false,
		}

		userInputJSON, err := json.Marshal(userInput)
		assert.NoError(t, err)

		rec := setupPost(&store.MockStore{
			Data: nil,
			Err:  gorm.ErrInvalidData,
		}, strings.NewReader(string(userInputJSON)))

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		var response struct {
			Message string `json:"message"`
		}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, gorm.ErrInvalidData.Error(), response.Message)
	})
}
