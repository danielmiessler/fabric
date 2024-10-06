package restapi

import (
	"github.com/danielmiessler/fabric/db"
	"github.com/labstack/echo/v4"
	"net/http"
)

// PatternsHandler defines the handler for patterns-related operations
type PatternsHandler struct {
	*StorageHandler[db.Pattern]
	patterns *db.PatternsEntity
}

// NewPatternsHandler creates a new PatternsHandler
func NewPatternsHandler(e *echo.Echo, patterns *db.PatternsEntity) (ret *PatternsHandler) {
	ret = &PatternsHandler{
		StorageHandler: NewStorageHandler[db.Pattern](e, "patterns", patterns), patterns: patterns}
	e.GET("/patterns/:name", ret.GetPattern)
	return
}

// GetPattern handles the GET /patterns/:name route
func (h *PatternsHandler) GetPattern(c echo.Context) error {
	name := c.Param("name")
	variables := make(map[string]string) // Assuming variables are passed somehow
	pattern, err := h.patterns.GetApplyVariables(name, variables)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, pattern)
}
