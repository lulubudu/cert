package db

// Database wraps all database operations.
type Database interface {
	UserDatabase
	CertDatabase
}
