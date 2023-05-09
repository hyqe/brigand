package storage

import (
	"github.com/google/uuid"
	"time"
)

type Metadata struct {
	Id        string    `bson:"id"`
	FileName  string    `bson:"filename"`
	CreatedAt time.Time `bson:"createdAt"`
}

func NewMetadata(filename string) *Metadata {
	return &Metadata{
		Id:        uuid.New().String(),
		FileName:  filename,
		CreatedAt: time.Now().UTC(),
	}

}
