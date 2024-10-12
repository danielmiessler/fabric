package restapi

import (
	"github.com/danielmiessler/fabric/db"
	"github.com/gin-gonic/gin"
)

// SessionsHandler defines the handler for sessions-related operations
type SessionsHandler struct {
	*StorageHandler[db.Session]
	sessions *db.SessionsEntity
}

// NewSessionsHandler creates a new SessionsHandler
func NewSessionsHandler(r *gin.Engine, sessions *db.SessionsEntity) (ret *SessionsHandler) {
	ret = &SessionsHandler{
		StorageHandler: NewStorageHandler[db.Session](r, "sessions", sessions), sessions: sessions}
	return ret
}
