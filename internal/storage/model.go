package storage

import "time"

type Metadata struct {
	Id        string    `bson:"id"`
	FileName  string    `bson:"filename"`
	CreatedAt time.Time `bson:"createdAt"`
}
