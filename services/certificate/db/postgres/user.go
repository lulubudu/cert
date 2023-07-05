package postgres

import (
	"certificate/db"
	"fmt"
)

// AddUser adds `user` into the database if there's no existing user with the
// same email address, and fills the db-generated fields like `UUID` and
// `CreatedAt` for `user`.
func (pg *Postgres) AddUser(user *db.User) error {
	// use blowfish and crypt to process password
	query := `
INSERT INTO users (name, email, password)
VALUES ($1, $2, crypt($3, gen_salt('bf')))
RETURNING uuid, created_at`
	if err := pg.QueryRow(query, user.Name, user.Email, user.Password).
		Scan(&user.UUID, &user.CreatedAt); err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

// GetUser returns the user with UUID `user.UUID`, if `user` does not exist or
// is not active, it returns an error.
func (pg *Postgres) GetUser(userUUID string) (*db.User, error) {
	user := &db.User{}
	query := `
SELECT (uuid, name, email, created_at) FROM users
WHERE uuid=$1`
	if err := pg.QueryRow(query, userUUID).
		Scan(&user.UUID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
		return nil, fmt.Errorf("failed to query for user: %w", err)
	}
	return user, nil
}

// DeleteUser sets the user with UUID `userUUID` as inactive.
func (pg *Postgres) DeleteUser(userUUID string) error {
	query := `
UPDATE users
SET active = False
WHERE uuid = $1`
	res, err := pg.Exec(query, userUUID)
	if err != nil {
		return fmt.Errorf("failed to execute sql statement: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if count != 1 {
		return fmt.Errorf("rows affected = %d, should be 1", count)
	}
	return nil
}
