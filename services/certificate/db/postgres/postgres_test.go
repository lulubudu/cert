package postgres_test

import (
	"certificate/db"
	"certificate/db/postgres"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func MockConnect(t *testing.T) (db.Database, sqlmock.Sqlmock, error) {
	mockDB, mock, err := sqlmock.New()
	assert.Nil(t, err)
	assert.NotNil(t, mockDB)
	assert.NotNil(t, mock)
	return &postgres.Postgres{DB: mockDB}, mock, err
}
