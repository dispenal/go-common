package mongodb

import (
	"context"
	"fmt"
	"time"

	common_utils "github.com/dispenal/go-common/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectTimeout  = 30 * time.Second
	maxConnIdleTime = 3 * time.Minute
	minPoolSize     = 20
	maxPoolSize     = 300
)

func NewMongoDBConn(ctx context.Context, cfg *common_utils.BaseConfig) (*mongo.Client, error) {
	mongoUri := fmt.Sprintf("mongodb://%s:%s@%s:%s/?replicaSet=%s",
		cfg.MongoUser,
		cfg.MongoPassword,
		cfg.MongoHost,
		cfg.MongoPort,
		cfg.MongoReplicaSet)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri).
		SetConnectTimeout(connectTimeout).
		SetMaxConnIdleTime(maxConnIdleTime).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(maxPoolSize))
	if err != nil {
		common_utils.PanicIfAppError(err, "mongo.Connect", 500)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}

type MongoDB interface {
	InsertOne(ctx context.Context, collection string, document any) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, collection string, documents []any) (*mongo.InsertManyResult, error)
	FindOne(ctx context.Context, collection string, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult
	Find(ctx context.Context, collection string, filter any, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOneAndUpdate(ctx context.Context, collection string, filter any, update any, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
	FindOneAndReplace(ctx context.Context, collection string, filter any, replacement any, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult
	FindOneAndDelete(ctx context.Context, collection string, filter any, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult
	UpdateOne(ctx context.Context, collection string, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, collection string, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, collection string, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(ctx context.Context, collection string, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	Aggregate(ctx context.Context, collection string, pipeline any, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
	CountDocuments(ctx context.Context, collection string, filter any, opts ...*options.CountOptions) (int64, error)
	EstimatedDocumentCount(ctx context.Context, collection string, opts ...*options.EstimatedDocumentCountOptions) (int64, error)
	Distinct(ctx context.Context, collection string, fieldName string, filter any, opts ...*options.DistinctOptions) ([]any, error)
	Watch(ctx context.Context, collection string, pipeline any, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error)
	Indexes(collection string) mongo.IndexView
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
	DropCollection(ctx context.Context, name string) error
	Client() *mongo.Client
	UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error
	StartSession() (mongo.Session, error)
	CreateCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error
}

type MongoDBClient struct {
	db  *mongo.Database
	cfg *common_utils.BaseConfig
}

func NewMongoDBClient(client *mongo.Client, cfg *common_utils.BaseConfig) MongoDB {

	db := client.Database(cfg.MongoDb)

	return &MongoDBClient{
		db:  db,
		cfg: cfg,
	}
}

func (m *MongoDBClient) InsertOne(ctx context.Context, collection string, document any) (*mongo.InsertOneResult, error) {
	coll := m.db.Collection(collection)
	return coll.InsertOne(ctx, document)
}

func (m *MongoDBClient) InsertMany(ctx context.Context, collection string, documents []any) (*mongo.InsertManyResult, error) {
	coll := m.db.Collection(collection)
	return coll.InsertMany(ctx, documents)
}

func (m *MongoDBClient) FindOne(ctx context.Context, collection string, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult {
	coll := m.db.Collection(collection)
	return coll.FindOne(ctx, filter, opts...)
}

func (m *MongoDBClient) Find(ctx context.Context, collection string, filter any, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	coll := m.db.Collection(collection)
	return coll.Find(ctx, filter, opts...)
}

func (m *MongoDBClient) FindOneAndUpdate(ctx context.Context, collection string, filter any, update any, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	coll := m.db.Collection(collection)
	return coll.FindOneAndUpdate(ctx, filter, update, opts...)
}

func (m *MongoDBClient) FindOneAndReplace(ctx context.Context, collection string, filter any, replacement any, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	coll := m.db.Collection(collection)
	return coll.FindOneAndReplace(ctx, filter, replacement, opts...)
}

func (m *MongoDBClient) FindOneAndDelete(ctx context.Context, collection string, filter any, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	coll := m.db.Collection(collection)
	return coll.FindOneAndDelete(ctx, filter, opts...)
}

func (m *MongoDBClient) UpdateOne(ctx context.Context, collection string, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	coll := m.db.Collection(collection)
	return coll.UpdateOne(ctx, filter, update, opts...)
}

func (m *MongoDBClient) UpdateMany(ctx context.Context, collection string, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	coll := m.db.Collection(collection)
	return coll.UpdateMany(ctx, filter, update, opts...)
}

func (m *MongoDBClient) DeleteOne(ctx context.Context, collection string, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	coll := m.db.Collection(collection)
	return coll.DeleteOne(ctx, filter, opts...)
}

func (m *MongoDBClient) DeleteMany(ctx context.Context, collection string, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	coll := m.db.Collection(collection)
	return coll.DeleteMany(ctx, filter, opts...)
}

func (m *MongoDBClient) Aggregate(ctx context.Context, collection string, pipeline any, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	coll := m.db.Collection(collection)
	return coll.Aggregate(ctx, pipeline, opts...)
}

func (m *MongoDBClient) CountDocuments(ctx context.Context, collection string, filter any, opts ...*options.CountOptions) (int64, error) {
	coll := m.db.Collection(collection)
	return coll.CountDocuments(ctx, filter, opts...)
}

func (m *MongoDBClient) EstimatedDocumentCount(ctx context.Context, collection string, opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	coll := m.db.Collection(collection)
	return coll.EstimatedDocumentCount(ctx, opts...)
}

func (m *MongoDBClient) Distinct(ctx context.Context, collection string, fieldName string, filter any, opts ...*options.DistinctOptions) ([]any, error) {
	coll := m.db.Collection(collection)
	return coll.Distinct(ctx, fieldName, filter, opts...)
}

func (m *MongoDBClient) Watch(ctx context.Context, collection string, pipeline any, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	coll := m.db.Collection(collection)
	return coll.Watch(ctx, pipeline, opts...)
}

func (m *MongoDBClient) Indexes(collection string) mongo.IndexView {
	return m.db.Collection(collection).Indexes()
}

func (m *MongoDBClient) Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	return m.db.Collection(name, opts...)
}

func (m *MongoDBClient) DropCollection(ctx context.Context, name string) error {
	return m.db.Collection(name).Drop(ctx)
}

func (m *MongoDBClient) Client() *mongo.Client {
	return m.db.Client()
}

func (m *MongoDBClient) UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error {
	return m.db.Client().UseSession(ctx, fn)
}

func (m *MongoDBClient) StartSession() (mongo.Session, error) {
	return m.db.Client().StartSession()
}

func (m *MongoDBClient) CreateCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error {
	return m.db.CreateCollection(ctx, name, opts...)
}
