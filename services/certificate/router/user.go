package router

import (
	"certificate/db"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

const userPath = "/user"

func (r *Router) routeUser() {
	r.POST(userPath, r.addUser)
	r.GET(userPath, r.getUser)
	r.DELETE(userPath, r.deleteUser)
}

// addUser adds a new user if the provided email address does not exist in the
// database.
func (r *Router) addUser(c echo.Context) error {
	// decode request body into `user`
	user := &db.User{}
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// add user to database and let it fill the db-generated fields of `user`
	if err := r.db.AddUser(user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Errorf("failed to add user: %w", err))
	}

	// wipe password and return user with db-generated fields
	user.Password = ""
	return c.JSON(http.StatusOK, user)
}

// getUser gets an existing user.
func (r *Router) getUser(c echo.Context) error {
	// decode request body into `user`
	user := &db.User{}
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// query database for user
	user, err := r.db.GetUser(user.UUID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Errorf("failed to get user: %w", err))
	}

	// write user to response
	return c.JSON(http.StatusOK, user)
}

// deleteUser deletes an existing user.
func (r *Router) deleteUser(c echo.Context) error {
	// decode request body into `user`
	user := &db.User{}
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// ask the database to delete user
	if err := r.db.DeleteUser(user.UUID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Errorf("failed to delete user %s: %w", user.UUID, err))
	}

	return c.String(http.StatusOK, "success!")
}
