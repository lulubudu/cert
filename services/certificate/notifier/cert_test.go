package notifier_test

import (
	"certificate/notifier"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

const mockCertUUID = "mock_cert_uuid"

func TestCertImpl_SendCertToggled(t *testing.T) {
	mn := &MockNotifier{}
	n := notifier.New(mn)
	t.Run("happy_path", func(t *testing.T) {
		assert.Nil(t, n.SendCertToggled(mockCertUUID, true))
		assert.Nil(t, n.SendCertToggled(mockCertUUID, false))
		assert.Equal(t, 2, len(mn.Messages))
		assert.Regexp(t, "{\"uuid\":\"mock_cert_uuid\",\"active\":true,\"updated_at\":\"(.+)T(.+)\"}", string(mn.Messages[0]))
		assert.Regexp(t, "{\"uuid\":\"mock_cert_uuid\",\"active\":false,\"updated_at\":\"(.+)T(.+)\"}", string(mn.Messages[1]))
	})
	t.Run("err_write_message", func(t *testing.T) {
		mn.Err = errors.New("mock_error")
		defer func() {
			mn.Err = nil
		}()
		assert.NotNil(t, n.SendCertToggled(mockCertUUID, true))
		assert.NotNil(t, n.SendCertToggled(mockCertUUID, false))
	})
}
