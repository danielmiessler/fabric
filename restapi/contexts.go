package restapi

import (
	"github.com/danielmiessler/fabric/db"
	"github.com/gin-gonic/gin"
)

// ContextsHandler defines the handler for contexts-related operations
type ContextsHandler struct {
	*StorageHandler[db.Context]
	contexts *db.ContextsEntity
}

// NewContextsHandler creates a new ContextsHandler
func NewContextsHandler(r *gin.Engine, contexts *db.ContextsEntity) (ret *ContextsHandler) {
	ret = &ContextsHandler{
		StorageHandler: NewStorageHandler[db.Context](r, "contexts", contexts), contexts: contexts}
	return
}
