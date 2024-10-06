package restapi

import (
	"github.com/danielmiessler/fabric/db"
	"github.com/labstack/echo/v4"
)

// SessionsHandler defines the handler for sessions-related operations
type SessionsHandler struct {
	*StorageHandler[db.Session]
	sessions *db.SessionsEntity
}

// NewSessionsHandler creates a new SessionsHandler
func NewSessionsHandler(e *echo.Echo, sessions *db.SessionsEntity) (ret *SessionsHandler) {
	ret = &SessionsHandler{
		StorageHandler: NewStorageHandler[db.Session](e, "sessions", sessions), sessions: sessions}
	return ret
}
