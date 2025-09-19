package service

import "errors"

var (
	DatabaseReadPathError       = errors.New("DATABASE_READ_PATH_ERROR: Error while reading database path")
	DatabaseWritePathError      = errors.New("DATABASE_WRITE_PATH_ERROR: Error while reading database path")
	InvalidDatabaseFormatAtPath = errors.New("INVALID_DATABASE_FORMAT_AT_PATH: Database path is not empty and is not a valid database")

	DatabaseManifestCreateError            = errors.New("DATABASE_MANIFEST_CREATE_ERROR: Error while creating database manifest")
	DatabaseJSONManifestSerializationError = errors.New("DATABASE_JSON_MANIFEST_SERIALIZATION_ERROR: Error while serializing database manifest")
	DatabaseManifestWriteError             = errors.New("DATABASE_MANIFEST_WRITE_ERROR: Error while writing database manifest")
	DatabaseManifestLoadError              = errors.New("DATABASE_MANIFEST_LOAD_ERROR: Error while loading database manifest")

	DatabaseDuplicate = errors.New("DATABASE_DUPLICATE: Database already exists")
)
