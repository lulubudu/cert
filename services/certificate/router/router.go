package router

import (
	"certificate/db"
	"certificate/notifier"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router struct {
	db       db.Database
	notifier *notifier.Notifier
	*echo.Echo
}

func New() *Router {
	r := &Router{Echo: echo.New()}
	r.Use(middleware.Logger())
	r.routeCert()
	r.routeUser()
	return r
}

func (r *Router) Start(address string) error {
	return r.Echo.Start(address)
}

func (r *Router) WithDatabase(db db.Database) *Router {
	r.db = db
	return r
}

func (r *Router) WithNotifier(notifier *notifier.Notifier) *Router {
	r.notifier = notifier
	return r
}
