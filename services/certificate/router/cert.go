package router

import (
	"certificate/db"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

const certPath = "/cert"

func (r *Router) routeCert() {
	r.POST(certPath, r.addCert)
	r.GET(certPath, r.getCerts)
	r.PATCH(certPath, r.setCertActiveStatus)
}

// addCert adds a certificate that belongs to an existing user.
func (r *Router) addCert(c echo.Context) error {
	// decode the request body into `cert`
	cert := &db.Cert{}
	if err := c.Bind(cert); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Errorf("failed to decode cert: %w", err))
	}

	// add cert to database and let it fill db-generated fields
	if err := r.db.AddCert(cert); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Errorf("failed to add cert: %w", err))
	}

	// send message to notifier
	if err := r.notifier.SendCertToggled(cert.UUID, cert.Active); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Errorf("failed to send cert toggled message: %w", err))
	}

	// write to response with generated fields
	return c.JSON(http.StatusOK, cert)
}

// getCerts returns all active certificates belonging to an existing user.
func (r *Router) getCerts(c echo.Context) error {
	// decode request body to get the querying user's UUID
	cert := &db.Cert{}
	if err := c.Bind(cert); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Errorf("failed to decode cert: %w", err))
	}

	// query the database for certificates belonging to this user
	certs, err := r.db.GetCerts(cert.UserUUID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Errorf("failed to get certs: %w", err))
	}

	// write `certs` to response
	return c.JSON(http.StatusOK, certs)
}

// setCertActiveStatus activates/deactivates an existing user's certificate
// according the `active` field in the request body, and sends a message
// through notifier.
func (r *Router) setCertActiveStatus(c echo.Context) error {
	// decode the request body into `cert`
	cert := &db.Cert{}
	if err := c.Bind(cert); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Errorf("failed to decode cert: %w", err))
	}

	// update the certificate's status in database to active
	if err := r.db.SetCertActiveStatus(cert.UUID, cert.UserUUID, cert.Active); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Errorf("failed to toggle cert status: %w", err))
	}

	// send message to notifier
	if err := r.notifier.SendCertToggled(cert.UUID, cert.Active); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Errorf("failed to send cert toggled message: %w", err))
	}

	return c.String(http.StatusOK, "success!")
}
