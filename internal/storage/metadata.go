package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate moq -out metadata_mock.go . MetadataClient
type MetadataClient interface {
	Create(ctx context.Context, md *Metadata) error
}

func NewMongoMetadataClient(c *mongo.Client) MetadataClient {
	return &mongoMetadataClient{
		Client: c,
	}
}

type mongoMetadataClient struct {
	*mongo.Client
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
