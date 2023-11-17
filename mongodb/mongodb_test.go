package mongodb

import (
	"context"
	"testing"

	common_utils "github.com/dispenal/go-common/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func loadBaseConfig() *common_utils.BaseConfig {
	cfg, err := common_utils.LoadBaseConfig("../", "test")
	if err != nil {
		common_utils.PanicIfError(err)
	}
	return cfg
}

func TestConnectionMongoDB(t *testing.T) {
	ctx := context.Background()
	config := loadBaseConfig()
	client, err := NewMongoDBConn(ctx, config)
	assert.Nil(t, err)
	assert.NotNil(t, client)

	t.Run("Ping MongoDB Connectiom", func(t *testing.T) {
		err = client.Ping(ctx, nil)
		assert.Nil(t, err)
	})

	t.Run("Disconnect MongoDB Connectiom", func(t *testing.T) {
		err = client.Disconnect(ctx)
		assert.Nil(t, err)
	})
}

func TestMongoClientDB(t *testing.T) {
	ctx := context.Background()
	config := loadBaseConfig()
	client, err := NewMongoDBConn(ctx, config)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	defer client.Disconnect(ctx)

	db := NewMongoDBClient(client, config)
	assert.NotNil(t, db)

	documents := []interface{}{
		map[string]interface{}{
			"aggregate_id": 2,
			"name":         "test2",
		},
		map[string]interface{}{
			"aggregate_id": 3,
			"name":         "test3",
		},
	}

	t.Run("Create Collection", func(t *testing.T) {
		err = db.CreateCollection(ctx, "test")
		assert.Nil(t, err)
	})

	t.Run("Check Collection", func(t *testing.T) {
		indexModel := mongo.IndexModel{
			Keys:    bson.M{"aggregate_id": 1},
			Options: options.Index().SetUnique(true),
		}
		collection := db.Collection("test")

		collection.Indexes().CreateOne(ctx, indexModel)

		assert.NotNil(t, collection)
	})

	t.Run("Insert One Document", func(t *testing.T) {
		result, err := db.InsertOne(ctx, "test", map[string]interface{}{
			"aggregate_id": 1,
			"name":         "test1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Insert Many Documents", func(t *testing.T) {
		result, err := db.InsertMany(ctx, "test", documents)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Find One Document", func(t *testing.T) {
		result := db.FindOne(ctx, "test", map[string]interface{}{
			"aggregate_id": 1,
		})
		assert.NotNil(t, result)

		var output map[string]interface{}
		err = result.Decode(&output)
		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, "test1", output["name"])

	})

	t.Run("Find Documents", func(t *testing.T) {
		result, err := db.Find(ctx, "test", map[string]interface{}{
			"aggregate_id": bson.M{"$gt": 1},
		})
		assert.Nil(t, err)
		assert.NotNil(t, result)

		var output []map[string]interface{}
		err = result.All(ctx, &output)
		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, 2, len(output))
	})

	// t.Run("Find One And Update Document", func(t *testing.T) {
	// 	result := db.FindOneAndUpdate(ctx, "test", map[string]interface{}{
	// 		"aggregate_id": 1,
	// 	}, map[string]interface{}{
	// 		"$set": map[string]interface{}{
	// 			"name": "test1_updated",
	// 		},
	// 	})
	// 	assert.NotNil(t, result)

	// 	updatedResult := db.FindOne(ctx, "test", map[string]interface{}{
	// 		"aggregate_id": 1,
	// 	})

	// 	var output map[string]interface{}
	// 	err = result.Decode(&updatedResult)
	// 	assert.Nil(t, err)
	// 	assert.NotNil(t, output)
	// 	assert.Equal(t, "test1_updated", output["name"])
	// })

	// t.Run("Find One And Replace Document", func(t *testing.T) {
	// 	result := db.FindOneAndReplace(ctx, "test", map[string]interface{}{
	// 		"aggregate_id": 1,
	// 	}, map[string]interface{}{
	// 		"aggregate_id": 1,
	// 		"name":         "test1_replaced",
	// 	})
	// 	assert.NotNil(t, result)

	// 	replacedResult := db.FindOne(ctx, "test", map[string]interface{}{
	// 		"aggregate_id": 1,
	// 	})

	// 	var output map[string]interface{}
	// 	err = result.Decode(&replacedResult)
	// 	assert.Nil(t, err)
	// 	assert.NotNil(t, output)
	// 	assert.Equal(t, "test1_replaced", output["name"])
	// })

	t.Run("Find One And Delete Document", func(t *testing.T) {
		result := db.FindOneAndDelete(ctx, "test", map[string]interface{}{
			"aggregate_id": 1,
		})
		assert.NotNil(t, result)

		var output map[string]interface{}
		err = result.Decode(&output)
		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, "test1", output["name"])
	})

	t.Run("Update One Document", func(t *testing.T) {
		result, err := db.UpdateOne(ctx, "test", map[string]interface{}{
			"aggregate_id": 2,
		}, map[string]interface{}{
			"$set": map[string]interface{}{
				"name": "test2",
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, result)

		resultUpdated := db.FindOne(ctx, "test", map[string]interface{}{
			"aggregate_id": 2,
		})
		assert.NotNil(t, resultUpdated)

		var output map[string]interface{}
		err = resultUpdated.Decode(&output)
		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, "test2", output["name"])

	})

	t.Run("Update Many Documents", func(t *testing.T) {
		result, err := db.UpdateMany(ctx, "test", map[string]interface{}{
			"aggregate_id": bson.M{"$gt": 1},
		}, map[string]interface{}{
			"$set": map[string]interface{}{
				"name": "test1",
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Delete One Document", func(t *testing.T) {
		result, err := db.DeleteOne(ctx, "test", map[string]interface{}{
			"aggregate_id": 2,
		})
		assert.Nil(t, err)
		assert.NotNil(t, result)

		resultDeleted := db.FindOne(ctx, "test", map[string]interface{}{
			"aggregate_id": 2,
		})
		assert.NotNil(t, resultDeleted)
		assert.Equal(t, mongo.ErrNoDocuments, resultDeleted.Err())

	})

	t.Run("Delete Many Documents", func(t *testing.T) {
		result, err := db.DeleteMany(ctx, "test", map[string]interface{}{
			"aggregate_id": bson.M{"$gt": 1},
		})
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Drop Collection", func(t *testing.T) {
		err = db.DropCollection(ctx, "test")
		assert.Nil(t, err)
	})

}
