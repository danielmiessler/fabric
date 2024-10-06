package restapi

import (
	"github.com/danielmiessler/fabric/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Serve(fabricDb *db.Db, address string) (err error) {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register routes
	NewPatternsHandler(e, fabricDb.Patterns)
	NewContextsHandler(e, fabricDb.Contexts)
	NewSessionsHandler(e, fabricDb.Sessions)

	// Start server
	e.Logger.Fatal(e.Start(address))

	return
}
