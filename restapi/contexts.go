package restapi

import (
	"github.com/danielmiessler/fabric/db"
	"github.com/labstack/echo/v4"
)

// ContextsHandler defines the handler for contexts-related operations
type ContextsHandler struct {
	*StorageHandler[db.Context]
	contexts *db.ContextsEntity
}

// NewContextsHandler creates a new ContextsHandler
func NewContextsHandler(e *echo.Echo, contexts *db.ContextsEntity) (ret *ContextsHandler) {
	ret = &ContextsHandler{
		StorageHandler: NewStorageHandler[db.Context](e, "contexts", contexts), contexts: contexts}
	return
}
