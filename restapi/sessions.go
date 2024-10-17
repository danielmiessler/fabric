package restapi

import (
	"github.com/danielmiessler/fabric/plugins/db/fsdb"
	"github.com/gin-gonic/gin"
)

// SessionsHandler defines the handler for sessions-related operations
type SessionsHandler struct {
	*StorageHandler[fsdb.Session]
	sessions *fsdb.SessionsEntity
}

// NewSessionsHandler creates a new SessionsHandler
func NewSessionsHandler(r *gin.Engine, sessions *fsdb.SessionsEntity) (ret *SessionsHandler) {
	ret = &SessionsHandler{
		StorageHandler: NewStorageHandler[fsdb.Session](r, "sessions", sessions), sessions: sessions}
	return ret
}
