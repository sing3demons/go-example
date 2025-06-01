package store

import (
	"gorm.io/gorm"
)

type Storer interface {
	Find(dest any, conds ...any) error
	Create(value any) error
	First(dest any, conds ...any) error
	Save(value any) error
}

type gormStore struct {
	db *gorm.DB
}
