package store_test

import (
	"errors"
	"testing"

	"github.com/sing3demons/go-example/store"
	"github.com/stretchr/testify/assert"
)

func TestMockStoreFind(t *testing.T) {
	expected := User{Name: "Alice", Age: 30}

	t.Run("successfully finds data", func(t *testing.T) {
		mock := &store.MockStore{Data: expected}
		var result User

		err := mock.Find(&result)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("returns error if set", func(t *testing.T) {
		mock := &store.MockStore{Err: errors.New("db error")}
		var result User

		err := mock.Find(&result)
		assert.EqualError(t, err, "db error")
	})
}

func TestMockStoreFirst(t *testing.T) {
	expected := User{Name: "Bob", Age: 25}

	t.Run("successfully retrieves first record", func(t *testing.T) {
		mock := &store.MockStore{Data: expected}
		var result User

		err := mock.First(&result)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestMockStoreCreate(t *testing.T) {
	t.Run("create without error", func(t *testing.T) {
		mock := &store.MockStore{}
		err := mock.Create(User{Name: "Charlie"})
		assert.NoError(t, err)
	})

	t.Run("create with error", func(t *testing.T) {
		mock := &store.MockStore{Err: errors.New("insert error")}
		err := mock.Create(User{Name: "Charlie"})
		assert.EqualError(t, err, "insert error")
	})
}

func TestMockStoreSave(t *testing.T) {
	expected := User{Name: "Dave", Age: 40}

	t.Run("successfully saves data", func(t *testing.T) {
		mock := &store.MockStore{Data: expected}
		var result User

		err := mock.Save(&result)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("save with error", func(t *testing.T) {
		mock := &store.MockStore{Err: errors.New("save error")}
		var result User

		err := mock.Save(&result)
		assert.EqualError(t, err, "save error")
	})
}
