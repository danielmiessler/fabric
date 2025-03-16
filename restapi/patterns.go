package restapi

import (
	"net/http"

	"github.com/danielmiessler/fabric/plugins/db/fsdb"
	"github.com/gin-gonic/gin"
)

// PatternsHandler defines the handler for patterns-related operations
type PatternsHandler struct {
	*StorageHandler[fsdb.Pattern]
	patterns *fsdb.PatternsEntity
}

// NewPatternsHandler creates a new PatternsHandler
func NewPatternsHandler(r *gin.Engine, patterns *fsdb.PatternsEntity) (ret *PatternsHandler) {
	ret = &PatternsHandler{
		StorageHandler: NewStorageHandler(r, "patterns", patterns), patterns: patterns}

	// TODO: Add custom, replacement routes here
	//r.GET("/patterns/:name", ret.Get)
	return
}

// Get handles the GET /patterns/:name route
func (h *PatternsHandler) Get(c *gin.Context) {
	name := c.Param("name")
	variables := make(map[string]string) // Assuming variables are passed somehow
	input := ""                          // Assuming input is passed somehow
	pattern, err := h.patterns.GetApplyVariables(name, variables, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, pattern)
}
