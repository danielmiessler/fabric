package restapi

import (
	"github.com/danielmiessler/fabric/core"
	"github.com/gin-gonic/gin"
)

func Serve(registry *core.PluginRegistry, address string) (err error) {
	r := gin.Default()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Register routes
	fabricDb := registry.Db
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
