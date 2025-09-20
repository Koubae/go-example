package model

import "time"

type Manifest struct {
	Name    string    `json:"name"`
	Created time.Time `json:"created"` // time.RFC3339Nano
	Updated time.Time `json:"updated"`
}

type Database struct {
	Manifest
	path string
}

func (d Database) Path() string {
	return d.path
}

func NewDatabase(manifest Manifest, path string) *Database {
	return &Database{
		Manifest: manifest,
		path:     path,
	}
}

type Record[T any] struct {
	Manifest Manifest `json:"manifest"`
	ID       string   `json:"id"`
	Data     T        `json:"data"`
}

type Table[T any] struct {
	Manifest Manifest             `json:"manifest"`
	Records  map[string]Record[T] `json:"records"`
}
