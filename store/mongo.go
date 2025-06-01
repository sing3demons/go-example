package store

import (
	"context"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoStore struct {
	col *mongo.Collection
}

func NewMongoStore(col *mongo.Collection) Storer {
	return &mongoStore{
		col: col,
	}
}

func (s *mongoStore) Find(dest any, conds ...any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.col.Find(ctx, conds[0])
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, dest); err != nil {
		return err
	}

	return nil

}

func (s *mongoStore) Create(value any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := s.col.InsertOne(ctx, value)
	if err != nil {
		return err
	}

	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil // Only process if it's a struct
	}
	idField := val.FieldByName("ID")
	if !idField.IsValid() || !idField.CanSet() {
		idField = val.FieldByName("_id")
	}
	if idField.IsValid() && idField.CanSet() {
		idField.Set(reflect.ValueOf(r.InsertedID))
	}

	return nil
}

func (s *mongoStore) First(result any, filter ...any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.col.FindOne(ctx, filter[0]).Decode(result); err != nil {
		return err
	}

	return nil
}

func (s *mongoStore) Save(value any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	val := reflect.ValueOf(value)

	// If value is a pointer, get the underlying struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Only process if it's a struct
	if val.Kind() != reflect.Struct {
		return nil
	}

	idField := val.FieldByName("ID")
	if !idField.IsValid() || !idField.CanSet() {
		idField = val.FieldByName("_id")
	}

	if idField.IsValid() && idField.CanSet() && isZero(idField) {
		return nil // Cannot save without a valid ID
	}

	filter := primitive.M{"_id": idField.Interface()}
	update := primitive.M{"$set": val.Interface()}

	_, err := s.col.UpdateOne(ctx, filter, update)
	return err
}

// isZero checks if a reflect.Value is the zero value for its type
func isZero(v reflect.Value) bool {
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
