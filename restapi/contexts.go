package restapi

import (
	"github.com/danielmiessler/fabric/plugins/db/fsdb"
	"github.com/gin-gonic/gin"
)

// ContextsHandler defines the handler for contexts-related operations
type ContextsHandler struct {
	*StorageHandler[fsdb.Context]
	contexts *fsdb.ContextsEntity
}

// NewContextsHandler creates a new ContextsHandler
func NewContextsHandler(r *gin.Engine, contexts *fsdb.ContextsEntity) (ret *ContextsHandler) {
	ret = &ContextsHandler{
		StorageHandler: NewStorageHandler[fsdb.Context](r, "contexts", contexts), contexts: contexts}
	return
}
