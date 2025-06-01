package store_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sing3demons/go-example/store"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID   int
	Name string
	Age  int
}

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	cleanup := func() {
		db.Close()
	}

	return gormDB, mock, cleanup
}

func TestGormStoreFind(t *testing.T) {
	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT .* FROM "users"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "age"}).
			AddRow(1, "Alice", 30))

	s := store.NewGormStore(gdb)
	var users []User
	err := s.Find(&users)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "Alice", users[0].Name)
}

func TestGormStoreFirst(t *testing.T) {
	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT .* FROM "users"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "age"}).
			AddRow(1, "Bob", 25))

	s := store.NewGormStore(gdb)
	var user User
	err := s.First(&user)
	assert.NoError(t, err)
	assert.Equal(t, "Bob", user.Name)
}

func TestGormStoreCreate(t *testing.T) {
	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users" \("name","age"\) VALUES \(\$1,\$2\) RETURNING "id"`).
		WithArgs("Charlie", 40).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	s := store.NewGormStore(gdb)
	user := User{Name: "Charlie", Age: 40}
	err := s.Create(&user)
	assert.NoError(t, err)
}

func TestGormStoreSave(t *testing.T) {
	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "users"`).
		WithArgs("Dave", 35, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	s := store.NewGormStore(gdb)
	user := User{ID: 1, Name: "Dave", Age: 35}
	err := s.Save(&user)
	assert.NoError(t, err)
}
