package restapi

import (
	"github.com/danielmiessler/fabric/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

// PatternsHandler defines the handler for patterns-related operations
type PatternsHandler struct {
	*StorageHandler[db.Pattern]
	patterns *db.PatternsEntity
}

// NewPatternsHandler creates a new PatternsHandler
func NewPatternsHandler(r *gin.Engine, patterns *db.PatternsEntity) (ret *PatternsHandler) {
	ret = &PatternsHandler{
		StorageHandler: NewStorageHandler[db.Pattern](r, "patterns", patterns), patterns: patterns}
	r.GET("/patterns/:name", ret.GetPattern)
	return
}

// GetPattern handles the GET /patterns/:name route
func (h *PatternsHandler) GetPattern(c *gin.Context) {
	name := c.Param("name")
	variables := make(map[string]string) // Assuming variables are passed somehow
	pattern, err := h.patterns.GetApplyVariables(name, variables)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, pattern)
}
