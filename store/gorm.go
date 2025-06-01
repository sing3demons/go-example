package store

import "gorm.io/gorm"

func NewGormStore(db *gorm.DB) Storer {
	return &gormStore{
		db: db,
	}
}

func (s *gormStore) Find(dest any, conds ...any) error {
	r := s.db.Find(dest, conds...)
	return r.Error
}

func (s *gormStore) Create(value any) error {
	return s.db.Create(value).Error
}
func (s *gormStore) First(dest any, conds ...any) error {
	return s.db.First(dest, conds...).Error
}
func (s *gormStore) Save(value any) error {
	return s.db.Save(value).Error
}
