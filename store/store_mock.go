package store

import "reflect"

type MockStore struct {
	Err  error
	Data any
}

func (m *MockStore) Find(dest any, conds ...any) error {
	if m.Err != nil {
		return m.Err
	}

	if m.Data != nil {
		// Ensure dest is a pointer
		destVal := reflect.ValueOf(dest)
		if destVal.Kind() != reflect.Ptr {
			return nil
		}

		// Set the value pointed to by dest
		reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(m.Data))
	}
	return nil
}

func (m *MockStore) Create(value any) error {

	if m.Err != nil {
		return m.Err
	}

	return nil
}

func (m *MockStore) First(dest any, conds ...any) error {

	return nil
}

func (m *MockStore) Save(value any) error {
	return nil
}
