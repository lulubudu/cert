package db

import (
	"time"
)

// Cert represents the database schema for certificates.
type Cert struct {
	UUID       string    `json:"uuid"`
	UserUUID   string    `json:"user_uuid"`
	PrivateKey string    `json:"private_key,omitempty"`
	Body       string    `json:"body,omitempty"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}

// CertDatabase is the interface that wraps all certificate related database
// operations.
type CertDatabase interface {
	AddCert(cert *Cert) error
	GetCerts(userUUID string) ([]*Cert, error)
	SetCertActiveStatus(certUUID, userUUID string, active bool) error
}
