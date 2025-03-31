package restapi

import (
	"github.com/danielmiessler/fabric/core"
	"github.com/gin-gonic/gin"
)

func Serve(registry *core.PluginRegistry, address string, apiKey string) (err error) {
	r := gin.New()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	if apiKey != "" {
		r.Use(ApiKeyMiddleware(apiKey))
	}

	// Register routes
	fabricDb := registry.Db
	NewPatternsHandler(r, fabricDb.Patterns)
	NewContextsHandler(r, fabricDb.Contexts)
	NewSessionsHandler(r, fabricDb.Sessions)
	NewChatHandler(r, registry, fabricDb)
	NewConfigHandler(r, fabricDb)
	NewModelsHandler(r, registry.VendorManager)

	// Start server
	err = r.Run(address)
	if err != nil {
		return err
	}

	return
}
