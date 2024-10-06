package restapi

import (
	"github.com/danielmiessler/fabric/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Serve(fabricDb *db.Db) {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register routes
	NewPatternsHandler(e, fabricDb.Patterns)
	NewContextsHandler(e, fabricDb.Patterns)
	NewSessionsHandler(e, fabricDb.Patterns)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
