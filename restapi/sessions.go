package restapi

import (
	"github.com/danielmiessler/fabric/plugins/db/fs"
	"github.com/gin-gonic/gin"
)

// SessionsHandler defines the handler for sessions-related operations
type SessionsHandler struct {
	*StorageHandler[fs.Session]
	sessions *fs.SessionsEntity
}

// NewSessionsHandler creates a new SessionsHandler
func NewSessionsHandler(r *gin.Engine, sessions *fs.SessionsEntity) (ret *SessionsHandler) {
	ret = &SessionsHandler{
		StorageHandler: NewStorageHandler[fs.Session](r, "sessions", sessions), sessions: sessions}
	return ret
}
