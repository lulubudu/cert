package postgres_test

import (
	"certificate/db"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var mockUser = &db.User{
	UUID:      "mock_user_uuid",
	Name:      "Dog",
	Email:     "dog@cat.com",
	CreatedAt: time.Now(),
}

func TestPostgres_AddUser(t *testing.T) {
	pg, mock, _ := MockConnect(t)

	mockUser.Password = "tuna"
	defer func() {
		mockUser.Password = ""
	}()

	t.Run("happy_path", func(t *testing.T) {
		user := &db.User{
			Name:     mockUser.Name,
			Email:    mockUser.Email,
			Password: mockUser.Password,
		}
		rows := sqlmock.NewRows([]string{"uuid", "created_at"}).
			AddRow(mockUser.UUID, mockUser.CreatedAt)
		mock.ExpectQuery(`
^INSERT INTO users (.+)
VALUES (.+)
RETURNING uuid, created_at*`).
			WithArgs(user.Name, user.Email, user.Password).
			WillReturnRows(rows)

		assert.Nil(t, pg.AddUser(user))
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Equal(t, mockUser, user)
	})

	t.Run("error_no_rows_returned", func(t *testing.T) {
		user := &db.User{
			Name:     mockUser.Name,
			Email:    mockUser.Email,
			Password: mockUser.Password,
		}
		rows := sqlmock.NewRows([]string{"uuid", "created_at"})
		mock.ExpectQuery(`
^INSERT INTO users (.+)
VALUES (.+)
RETURNING uuid, created_at*`).
			WithArgs(user.Name, user.Email, user.Password).
			WillReturnRows(rows)

		assert.NotNil(t, pg.AddUser(user))
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.NotEqual(t, mockUser, user)
	})
}

func TestPostgres_GetUser(t *testing.T) {
	pg, mock, _ := MockConnect(t)

	t.Run("happy_path", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"uuid", "name", "email", "created_at"}).
			AddRow(mockUser.UUID, mockUser.Name, mockUser.Email, mockUser.CreatedAt)

		mock.ExpectQuery(`
^SELECT (.+) FROM users
WHERE (.+)*`).
			WithArgs(mockUser.UUID).
			WillReturnRows(rows)

		resultUser, err := pg.GetUser(mockUser.UUID)
		assert.Nil(t, err)
		assert.NotNil(t, resultUser)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Equal(t, mockUser, resultUser)
	})

	t.Run("error_no_rows_returned", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"uuid", "name", "email", "created_at"})

		mock.ExpectQuery(`
^SELECT (.+) FROM users
WHERE (.+)*`).
			WithArgs(mockUser.UUID).
			WillReturnRows(rows)

		resultUser, err := pg.GetUser(mockUser.UUID)
		assert.NotNil(t, err)
		assert.Nil(t, resultUser)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPostgres_DeleteUser(t *testing.T) {
	pg, mock, _ := MockConnect(t)

	t.Run("happy_path", func(t *testing.T) {
		mock.ExpectExec(`
^UPDATE users
SET active = False
WHERE (.+)*`).
			WithArgs(mockUser.UUID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := pg.DeleteUser(mockUser.UUID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("error_no_row_affected", func(t *testing.T) {
		mock.ExpectExec(`
^UPDATE users
SET active = False
WHERE (.+)*`).
			WithArgs(mockUser.UUID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := pg.DeleteUser(mockUser.UUID)
		assert.NotNil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
