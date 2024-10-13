package restapi

import (
	"github.com/danielmiessler/fabric/plugins/db/db_fs"
	"github.com/gin-gonic/gin"
)

// ContextsHandler defines the handler for contexts-related operations
type ContextsHandler struct {
	*StorageHandler[db_fs.Context]
	contexts *db_fs.ContextsEntity
}

// NewContextsHandler creates a new ContextsHandler
func NewContextsHandler(r *gin.Engine, contexts *db_fs.ContextsEntity) (ret *ContextsHandler) {
	ret = &ContextsHandler{
		StorageHandler: NewStorageHandler[db_fs.Context](r, "contexts", contexts), contexts: contexts}
	return
}
