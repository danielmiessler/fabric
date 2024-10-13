package restapi

import (
	"github.com/danielmiessler/fabric/plugins/db/db_fs"
	"github.com/gin-gonic/gin"
)

// SessionsHandler defines the handler for sessions-related operations
type SessionsHandler struct {
	*StorageHandler[db_fs.Session]
	sessions *db_fs.SessionsEntity
}

// NewSessionsHandler creates a new SessionsHandler
func NewSessionsHandler(r *gin.Engine, sessions *db_fs.SessionsEntity) (ret *SessionsHandler) {
	ret = &SessionsHandler{
		StorageHandler: NewStorageHandler[db_fs.Session](r, "sessions", sessions), sessions: sessions}
	return ret
}
