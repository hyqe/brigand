package storage

import "time"

type Metadata struct {
	Id       string
	FileName string
	CreateAt time.Time
}

func NewMetadata(filename string) *Metadata {
	return &Metadata{
		Id:       "none",
		FileName: filename,
		CreateAt: time.Now(),
	}
}
