package infraestructure

import (
	"API/src/Sensores/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleteDatosController struct {
	useCase application.DeleteDatos
}

func NewDeleteDatosController(useCase application.DeleteDatos) *DeleteDatosController {
	return &DeleteDatosController{useCase: useCase}
}

func (ds_c *DeleteDatosController) Execute(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pago ID"})
		return
	}

	ds_c.useCase.Execute(id)
	c.JSON(http.StatusOK, gin.H{"message": "Pago deleted successfully"})
}
