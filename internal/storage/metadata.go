package storage

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MetadataClient interface {
	Create(ctx context.Context, md *Metadata) error
	GetById(ctx context.Context, id string) (*Metadata, error)
	DeleteById(ctx context.Context, id string) error
}

func NewMongoMetadataClient(c *mongo.Client) MetadataClient {
	return &mongoMetadataClient{
		Client: c,
	}
}

// MockMetadataClient is a mock of MetadataClient which can be
// used by tests.
type MockMetadataClient struct {
	CreateFunc func(ctx context.Context, md *Metadata) error
}

func (m *MockMetadataClient) Create(ctx context.Context, md *Metadata) error {
	return m.CreateFunc(ctx, md)
}
func (m *MockMetadataClient) GetById(ctx context.Context, id string) (*Metadata, error) {
	return NewMetadata(id), nil
}

func (m *MockMetadataClient) DeleteById(ctx context.Context, id string) error {
	return m.DeleteById(ctx, id)
}

type mongoMetadataClient struct {
	*mongo.Client
}

func (m *mongoMetadataClient) GetById(ctx context.Context, id string) (*Metadata, error) {
	filter := bson.D{{Key: "id", Value: id}}
	coll := metadataColl(m.Client)

	var result Metadata
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (m *mongoMetadataClient) Create(ctx context.Context, md *Metadata) error {
	_, err := metadataColl(m.Client).InsertOne(ctx, md)
	return err
}

func (m *mongoMetadataClient) DeleteById(ctx context.Context, id string) error {
	filter := bson.D{{Key: "id", Value: id}}
	coll := metadataColl(m.Client)

	result, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount < 1 {
		return fmt.Errorf("there is no value to delete by this id: %s", id)
	}

	return nil
}

func brigandDB(m *mongo.Client) *mongo.Database {
	return m.Database("brigand")
}

func metadataColl(m *mongo.Client) *mongo.Collection {
	return brigandDB(m).Collection("metadata")
}

func CreateMetadataIndex(ctx context.Context, coll *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := coll.Indexes().CreateOne(context.TODO(), indexModel)

	return err
}
