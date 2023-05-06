package storage

import "time"

type Metadata struct {
	Id       string
	FileName string
	CreateAt time.Time
}
