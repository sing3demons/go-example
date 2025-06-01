package store_test

import (
	"testing"

	"github.com/sing3demons/go-example/store"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

type UserMock struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string
	Age  int
}

func TestMongoStoreCreate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("create document and assign ID", func(mt *mtest.T) {
		col := mt.Coll
		s := store.NewMongoStore(col)

		user := &UserMock{Name: "Alice", Age: 30}

		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{Key: "insertedId", Value: primitive.NewObjectID()},
		))

		err := s.Create(user)
		assert.NoError(t, err)
		assert.NotEqual(t, primitive.NilObjectID, user.ID)
	})
}

func TestMongoStoreFind(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("find multiple users", func(mt *mtest.T) {
		col := mt.Coll
		s := store.NewMongoStore(col)

		expected := bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "name", Value: "Bob"},
			{Key: "age", Value: 25},
		}

		// 1. First response = cursor with one document
		first := mtest.CreateCursorResponse(1, "test.users", mtest.FirstBatch, expected)

		// 2. Second response = cursor end (empty batch)
		second := mtest.CreateCursorResponse(1, "test.users", mtest.NextBatch)

		mt.AddMockResponses(first, second)

		var results []User
		err := s.Find(&results, bson.M{})
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Bob", results[0].Name)
	})
}

func TestMongoStoreFirst(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("find one user", func(mt *mtest.T) {
		col := mt.Coll
		s := store.NewMongoStore(col)

		expected := bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "name", Value: "Carol"},
			{Key: "age", Value: 22},
		}
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.users", mtest.FirstBatch, expected))

		var result UserMock
		err := s.First(&result, bson.M{"name": "Carol"})
		assert.NoError(t, err)
		assert.Equal(t, "Carol", result.Name)
	})
}

func TestMongoStoreSave(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("update document by ID", func(mt *mtest.T) {
		col := mt.Coll
		s := store.NewMongoStore(col)

		user := &UserMock{
			ID:   primitive.NewObjectID(),
			Name: "UpdatedName",
			Age:  35,
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := s.Save(user)
		assert.NoError(t, err)
	})
}
