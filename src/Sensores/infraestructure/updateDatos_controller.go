//File: updateDatos_controller.go

package infraestructure

import (
	"API/src/Sensores/application"
	"net/http"
	"strconv"
	"log"

	"github.com/gin-gonic/gin"
)

type UpdateDatosController struct {
	useCase application.UpdateDatos
}

func NewUpdateDatosController(useCase application.UpdateDatos) *UpdateDatosController {
	return &UpdateDatosController{useCase: useCase}
}

// UpdateDatosRequest define la estructura esperada para la actualización.
// Podría ser la misma que CreateDatosRequest si todos los campos son actualizables.
type UpdateDatosRequest struct {
	Temperatura string `json:"temperatura" binding:"required"`
	Movimiento  string `json:"movimiento" binding:"required"`
	Distancia   string `json:"distancia" binding:"required"`
	Peso        string `json:"peso" binding:"required"`
	Mac         string `json:"mac" binding:"required"`
}

func (udc *UpdateDatosController) Execute(c *gin.Context) {
	// Obtener ID del parámetro de la URL
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		log.Printf("ERROR: [UpdateCtrl] ID inválido en la URL: '%s'", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido proporcionado"})
		return
	}

	// Validar y parsear el cuerpo JSON
	var requestBody UpdateDatosRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("ERROR: [UpdateCtrl] Datos inválidos en la solicitud (ID: %d): %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Datos inválidos o faltantes en la solicitud de actualización",
			"detail": err.Error(),
		})
		return
	}

	// Ejecutar el caso de uso de actualización
	err = udc.useCase.Execute(
		id,
		requestBody.Temperatura,
		requestBody.Movimiento,
		requestBody.Distancia,
		requestBody.Peso,
		requestBody.Mac,
	)
	if err != nil {
		// Diferenciar errores si es posible (ej: NotFound vs InternalError)
		log.Printf("ERROR: [UpdateCtrl] Falló la ejecución del caso de uso UpdateDatos (ID: %d): %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al actualizar los datos del sensor"})
		return
	}

	log.Printf("INFO: [UpdateCtrl] Datos actualizados exitosamente para ID %d.", id)
	c.JSON(http.StatusOK, gin.H{"message": "Datos del sensor actualizados exitosamente"})
}