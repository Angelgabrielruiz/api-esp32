package infraestructure

import (
	"API/src/Sensores/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetDatosController struct {
	useCase application.GetDatos
}

func NewGetDatosController(useCase application.GetDatos) *GetDatosController {
	return &GetDatosController{useCase: useCase}
}

func (gp_c *GetDatosController) Execute(c *gin.Context) {
	pagos, err := gp_c.useCase.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los datos del sensor"})
		return
	}

	c.JSON(http.StatusOK, pagos)
}
