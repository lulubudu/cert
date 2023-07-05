package postgres_test

import (
	"certificate/db"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var mockCert0 = &db.Cert{
	UUID:       "mock_cert_uuid_0",
	UserUUID:   mockUser.UUID,
	PrivateKey: "private_key",
	Body:       "cert_body",
	Active:     true,
	CreatedAt:  time.Now(),
}

var mockCert1 = &db.Cert{
	UUID:       "mock_cert_uuid_1",
	UserUUID:   mockUser.UUID,
	PrivateKey: "private_key",
	Body:       "cert_body",
	Active:     true,
	CreatedAt:  time.Now(),
}

func TestPostgres_AddCert(t *testing.T) {
	pg, mock, _ := MockConnect(t)

	t.Run("happy_path", func(t *testing.T) {
		cert := &db.Cert{
			UserUUID:   mockCert0.UserUUID,
			PrivateKey: mockCert0.PrivateKey,
			Body:       mockCert0.Body,
		}

		mock.ExpectBegin()

		rows := sqlmock.NewRows([]string{"uuid"}).
			AddRow(mockUser.UUID)

		mock.ExpectQuery(`
^SELECT uuid FROM users
WHERE (.+)*`).
			WithArgs(mockUser.UUID).
			WillReturnRows(rows)

		rows = sqlmock.NewRows([]string{"uuid", "active", "created_at"}).
			AddRow(mockCert0.UUID, mockCert0.Active, mockCert0.CreatedAt)

		mock.ExpectQuery(`
^INSERT INTO certificates (.+)
VALUES (.+)
RETURNING (.+)*`).
			WithArgs(cert.UserUUID, cert.PrivateKey, cert.Body).
			WillReturnRows(rows)
		mock.ExpectCommit()

		assert.Nil(t, pg.AddCert(cert))
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Equal(t, mockCert0, cert)
	})
	t.Run("error_invalid_user_uuid_with_tx_rollback", func(t *testing.T) {
		cert := &db.Cert{
			UserUUID:   mockCert0.UserUUID,
			PrivateKey: mockCert0.PrivateKey,
			Body:       mockCert0.Body,
		}

		mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"uuid"})
		mock.ExpectQuery(`
^SELECT uuid FROM users
WHERE (.+)*`).
			WithArgs(mockUser.UUID).
			WillReturnRows(rows)
		mock.ExpectRollback()

		assert.NotNil(t, pg.AddCert(cert))
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.NotEqual(t, mockCert0, cert)
	})
}

func TestPostgres_GetCerts(t *testing.T) {
	pg, mock, _ := MockConnect(t)

	t.Run("happy_path", func(t *testing.T) {
		mock.ExpectBegin()

		rows := sqlmock.NewRows([]string{"uuid"}).
			AddRow(mockUser.UUID)

		mock.ExpectQuery(`
^SELECT uuid FROM users
WHERE (.+)*`).
			WithArgs(mockUser.UUID).
			WillReturnRows(rows)

		rows = sqlmock.NewRows([]string{"uuid", "private_key", "body", "active", "created_at"}).
			AddRow(mockCert0.UUID, mockCert0.PrivateKey, mockCert0.Body, mockCert0.Active, mockCert0.CreatedAt).
			AddRow(mockCert1.UUID, mockCert1.PrivateKey, mockCert1.Body, mockCert1.Active, mockCert1.CreatedAt)
		mock.ExpectQuery(`
^SELECT (.+) FROM certificates
WHERE (.+)*`).
			WithArgs(mockUser.UUID).
			WillReturnRows(rows)

		mock.ExpectCommit()

		certs, err := pg.GetCerts(mockUser.UUID)
		assert.Nil(t, err)
		assert.EqualValues(t, []*db.Cert{mockCert0, mockCert1}, certs)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("error_invalid_user_uuid_with_tx_rollback", func(t *testing.T) {
		mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"uuid"})
		mock.ExpectQuery(`
^SELECT uuid FROM users
WHERE (.+)*`).
			WithArgs(mockUser.UUID).
			WillReturnRows(rows)
		mock.ExpectRollback()

		certs, err := pg.GetCerts(mockUser.UUID)
		assert.NotNil(t, err)
		assert.Nil(t, certs)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPostgres_SetCertActiveStatus(t *testing.T) {
	pg, mock, _ := MockConnect(t)

	t.Run("happy_path", func(t *testing.T) {
		mock.ExpectBegin()

		rows := sqlmock.NewRows([]string{"uuid"}).
			AddRow(mockUser.UUID)

		mock.ExpectQuery(`
^SELECT uuid FROM users
WHERE (.+)*`).
			WithArgs(mockUser.UUID).
			WillReturnRows(rows)

		mock.ExpectExec(`
^UPDATE certificates
SET (.+)*
WHERE (.+)*`).
			WithArgs(mockCert0.UUID, mockCert0.Active).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		err := pg.SetCertActiveStatus(mockCert0.UUID, mockCert0.UserUUID, mockCert0.Active)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("error_invalid_user_uuid_with_tx_rollback", func(t *testing.T) {
		mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"uuid"})
		mock.ExpectQuery(`
^SELECT uuid FROM users
WHERE (.+)*`).
			WithArgs(mockUser.UUID).
			WillReturnRows(rows)
		mock.ExpectRollback()

		err := pg.SetCertActiveStatus(mockCert0.UUID, mockCert0.UserUUID, mockCert0.Active)
		assert.NotNil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
