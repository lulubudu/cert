package notifier

import (
	"encoding/json"
	"fmt"
	"time"
)

// CertNotifier is the interface that wraps the `SendCertToggled` method.
type CertNotifier interface {
	SendCertToggled(uuid string, active bool) error
}

// SendCertToggled writes a JSON message using its Writer.
func (n *Notifier) SendCertToggled(uuid string, active bool) error {
	// include an UpdatedAt timestamp to message
	jsonCert, err := json.Marshal(struct {
		UUID      string    `json:"uuid"`
		Active    bool      `json:"active"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		UUID:      uuid,
		Active:    active,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal cert: %w", err)
	}
	// TODO: add retry logic if write fails
	if err := n.Writer.WriteMessage(jsonCert); err != nil {
		return fmt.Errorf("failed to send kafka message: %w", err)
	}
	return nil
}
