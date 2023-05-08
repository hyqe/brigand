package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MetadataClient interface {
	Create(ctx context.Context, md *Metadata) error
	GetById(ctx context.Context, id string) (*Metadata, error)
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
	return m.GetById(ctx, id)
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

func brigandDB(m *mongo.Client) *mongo.Database {
	return m.Database("brigand")
}

func metadataColl(m *mongo.Client) *mongo.Collection {
	return brigandDB(m).Collection("metadata")
}
