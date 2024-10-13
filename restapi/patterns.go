package restapi

import (
	"github.com/danielmiessler/fabric/plugins/db/db_fs"
	"github.com/gin-gonic/gin"
	"net/http"
)

// PatternsHandler defines the handler for patterns-related operations
type PatternsHandler struct {
	*StorageHandler[db_fs.Pattern]
	patterns *db_fs.PatternsEntity
}

// NewPatternsHandler creates a new PatternsHandler
func NewPatternsHandler(r *gin.Engine, patterns *db_fs.PatternsEntity) (ret *PatternsHandler) {
	ret = &PatternsHandler{
		StorageHandler: NewStorageHandler[db_fs.Pattern](r, "patterns", patterns), patterns: patterns}

	// TODO: Add custom, replacement routes here
	//r.GET("/patterns/:name", ret.Get)
	return
}

// Get handles the GET /patterns/:name route
func (h *PatternsHandler) Get(c *gin.Context) {
	name := c.Param("name")
	variables := make(map[string]string) // Assuming variables are passed somehow
	pattern, err := h.patterns.GetApplyVariables(name, variables)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, pattern)
}
