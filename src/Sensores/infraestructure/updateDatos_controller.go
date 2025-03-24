package infraestructure

import (
	"API/src/Sensores/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UpdateDatosController struct {
	useCase application.UpdateDatos
}

func NewUpdateDatosController(useCase application.UpdateDatos) *UpdateDatosController {
	return &UpdateDatosController{useCase: useCase}
}

func (usc *UpdateDatosController) Execute(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de sensor inválido"})
		return
	}

	var input struct {
		Temperatura string `json:"temperatura"`
		Movimiento  string `json:"movimiento"`
		Distancia 	string `json:"distancia"`
		Peso	  	string `json:"peso"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := usc.useCase.Execute(id, input.Temperatura, input.Movimiento, input.Distancia, input.Peso); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar los datos del sensor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Datos del sensor actualizados exitosamente"})
}
