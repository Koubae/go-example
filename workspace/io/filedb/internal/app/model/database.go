package model

import "time"

type Manifest struct {
	Name    string    `json:"name"`
	Created time.Time `json:"created"` // time.RFC3339Nano
	Updated time.Time `json:"updated"`
}

type Database struct {
	Manifest
}
