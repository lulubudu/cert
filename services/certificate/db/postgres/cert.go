package postgres

import (
	"certificate/db"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// checkUser checks if userUUID is valid and active
func checkUser(tx *sql.Tx, userUUID string) error {
	query := `
SELECT uuid FROM users
WHERE uuid = $1 AND active`
	err := tx.QueryRow(query, userUUID).Scan(&userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user does not exist or is not active: %w", err)
		}
		return fmt.Errorf("failed to query for user_uuid: %w", err)
	}
	return nil
}

// AddCert adds cert to the database if `cert.UserUUID` exists and is active,
// and fills `cert` with db-generated fields like `UUID` and `CreatedAt`.
func (pg *Postgres) AddCert(cert *db.Cert) error {
	// use transaction for atomicity
	tx, err := pg.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}

	if err := checkUser(tx, cert.UserUUID); err != nil {
		return errors.Join(err, tx.Rollback())
	}

	// insert cert and fills the auto generated fields in Cert
	query := `
INSERT INTO certificates (user_uuid, private_key, body)
VALUES ($1, $2, $3)
RETURNING uuid, active, created_at`
	if err := tx.QueryRow(query, cert.UserUUID, cert.PrivateKey, cert.Body).
		Scan(&cert.UUID, &cert.Active, &cert.CreatedAt); err != nil {
		return errors.Join(fmt.Errorf("failed to insert certificate: %w", err), tx.Rollback())
	}

	// commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

// GetCerts returns all active certificates belonging to `userUUID`, it errors
// out if the user does not exist or is not active.
func (pg *Postgres) GetCerts(userUUID string) ([]*db.Cert, error) {
	// use transaction for atomicity
	tx, err := pg.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}

	if err := checkUser(tx, userUUID); err != nil {
		return nil, errors.Join(err, tx.Rollback())
	}

	// query for active certificates
	query := `
SELECT uuid, private_key, body, active, created_at FROM certificates
WHERE user_uuid = $1 AND active`
	rows, err := tx.Query(query, userUUID)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to query: %w", err), tx.Rollback())
	}
	var certs []*db.Cert
	for rows.Next() {
		cert := &db.Cert{UserUUID: userUUID}
		if errScan := rows.Scan(&cert.UUID, &cert.PrivateKey, &cert.Body, &cert.Active, &cert.CreatedAt); errScan != nil {
			err = errors.Join(err, fmt.Errorf("failed to scan row: %w", errScan))
		} else {
			certs = append(certs, cert)
		}
	}
	if err := rows.Close(); err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to close rows: %w", err), tx.Rollback())
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to commit tx: %w", err))
	}
	return certs, err
}

func updateCertActiveStatus(tx *sql.Tx, uuid string, active bool) error {
	// update db only if active status is different from cert.active
	query := `
UPDATE certificates
SET active = $2
WHERE uuid = $1 AND active != $2`
	res, err := tx.Exec(query, uuid, active)
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

// SetCertActiveStatus updates the active field of a certificate if needed, it
// errors out if the user does not exist or is not active.
// TODO: assumption - cert status cannot be changed after user deletion
func (pg *Postgres) SetCertActiveStatus(uuid, userUUID string, active bool) error {
	// use transaction for atomicity
	tx, err := pg.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}

	if err = checkUser(tx, userUUID); err != nil {
		return errors.Join(err, tx.Rollback())
	}

	if err = updateCertActiveStatus(tx, uuid, active); err != nil {
		return errors.Join(err, tx.Rollback())
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return errors.Join(err, fmt.Errorf("failed to commit tx: %w", err))
	}
	return nil
}
