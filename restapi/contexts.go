package restapi

import (
	"github.com/danielmiessler/fabric/plugins/db/fs"
	"github.com/gin-gonic/gin"
)

// ContextsHandler defines the handler for contexts-related operations
type ContextsHandler struct {
	*StorageHandler[fs.Context]
	contexts *fs.ContextsEntity
}

// NewContextsHandler creates a new ContextsHandler
func NewContextsHandler(r *gin.Engine, contexts *fs.ContextsEntity) (ret *ContextsHandler) {
	ret = &ContextsHandler{
		StorageHandler: NewStorageHandler[fs.Context](r, "contexts", contexts), contexts: contexts}
	return
}
