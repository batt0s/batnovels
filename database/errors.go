package database

import "errors"

var (
	// Database Connection
	ErrorDatabaseDriverInvalid = errors.New("given database driver is invalid")
	ErrorDatabaseSourceInvalid = errors.New("given database source is invalid")
	ErrorConnectionFailed      = errors.New("can not connect to database")
	// Migrations
	ErrorMigrationFailed = errors.New("failed to migrate database")
	// Query operations
	ErrorRecordNotFound  = errors.New("record with given query not found")
	ErrorOperationFailed = errors.New("operation failed")
	// context
	ErrorOperationCanceled = errors.New("operation canceled")
	// Validations
	ErrorInvalidUser    = errors.New("invalid user")
	ErrorInvalidProject = errors.New("invalid project")
	//
	ErrorNotImplemented = errors.New("not yet implemented")
)
