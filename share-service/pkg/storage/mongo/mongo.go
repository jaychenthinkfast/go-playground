package mongo

import (
	"context"
	"time"

	"github.com/playground/share-service/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewMongoStorage 创建新的 MongoDB 存储实例
func NewMongoStorage(ctx context.Context, uri, database, collection string) (*MongoStorage, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// 测试连接
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	// 创建索引
	col := client.Database(database).Collection(collection)
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "shareId", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "expires_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	}

	if _, err := col.Indexes().CreateMany(ctx, indexes); err != nil {
		return nil, err
	}

	return &MongoStorage{
		client:     client,
		collection: col,
	}, nil
}

// CreateShare 实现 Storage 接口
func (s *MongoStorage) CreateShare(ctx context.Context, share *models.Share) error {
	_, err := s.collection.InsertOne(ctx, share)
	return err
}

// GetShare 实现 Storage 接口
func (s *MongoStorage) GetShare(ctx context.Context, shareId string) (*models.Share, error) {
	var share models.Share
	err := s.collection.FindOne(ctx, bson.M{"shareId": shareId}).Decode(&share)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &share, err
}

// IncrementViews 实现 Storage 接口
func (s *MongoStorage) IncrementViews(ctx context.Context, shareId string) error {
	update := bson.M{
		"$inc": bson.M{"views": 1},
		"$set": bson.M{"last_viewed": time.Now()},
	}
	_, err := s.collection.UpdateOne(ctx, bson.M{"shareId": shareId}, update)
	return err
}

// DeleteExpiredShares 实现 Storage 接口
func (s *MongoStorage) DeleteExpiredShares(ctx context.Context) error {
	_, err := s.collection.DeleteMany(ctx, bson.M{
		"expires_at": bson.M{"$lt": time.Now()},
	})
	return err
}

// Close 实现 Storage 接口
func (s *MongoStorage) Close(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}
