package restapi

import (
	"github.com/danielmiessler/fabric/db"
	"github.com/gin-gonic/gin"
)

func Serve(fabricDb *db.Db, address string) (err error) {
	r := gin.Default()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Register routes
	NewPatternsHandler(r, fabricDb.Patterns)
	NewContextsHandler(r, fabricDb.Contexts)
	NewSessionsHandler(r, fabricDb.Sessions)

	// Start server
	err = r.Run(address)
	if err != nil {
		return err
	}

	return
}
