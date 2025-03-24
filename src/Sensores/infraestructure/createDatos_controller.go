package infraestructure

import (
	"API/src/Sensores/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateDatosController struct {
	useCase application.CreateDatos
}

func NewCreateDatosController(useCase application.CreateDatos) *CreateDatosController {
	return &CreateDatosController{useCase: useCase}
}

func (csc *CreateDatosController) Execute(c *gin.Context) {
	var requestBody struct {
		Temperatura string `json:"temperatura"`
		Movimiento  string `json:"movimiento"`
		Distancia   string `json:"distancia"`
		Peso        string `json:"peso"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	if err := csc.useCase.Execute(requestBody.Temperatura, requestBody.Movimiento, requestBody.Distancia, requestBody.Peso); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al registrar los datos del sensor",
			"detail": err.Error(), // Aquí mostramos el error detallado
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Datos del sensor registrados exitosamente"})
}
