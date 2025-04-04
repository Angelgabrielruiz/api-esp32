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

type UpdateDatosRequest struct {
	Temperatura string `json:"temperatura" binding:"required"`
	Movimiento  string `json:"movimiento" binding:"required"`
	Distancia   string `json:"distancia" binding:"required"`
	Peso        string `json:"peso" binding:"required"`
	Mac         string `json:"mac" binding:"required"` // Puede que no necesites la MAC aquí si usas el ID
}

func (udc *UpdateDatosController) Execute(c *gin.Context) {
	// --- OBTENER USER ID DEL CONTEXTO ---
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado (contexto inválido)"})
		return
	}
	userID, ok := userIDValue.(int)
	if !ok || userID <= 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno de autorización"})
		return
	}
	// --- FIN OBTENER USER ID ---

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

	// Ejecutar el caso de uso de actualización, pasando el userID para validación
	// Asume que UpdateUseCase.Execute ahora toma userID como parámetro
	err = udc.useCase.Execute(
		id,
		// userID, // <--- Pasa el userID si tu caso de uso/repo lo necesita para validar pertenencia
		requestBody.Temperatura,
		requestBody.Movimiento,
		requestBody.Distancia,
		requestBody.Peso,
		requestBody.Mac,
	)
	if err != nil {
		// Aquí podrías diferenciar errores, ej: si el repo devuelve "no encontrado" o "no pertenece al usuario"
		log.Printf("ERROR: [UpdateCtrl] Falló la ejecución del caso de uso UpdateDatos (ID: %d, UserID: %d): %v", id, userID, err)
		// Podrías devolver 404 Not Found o 403 Forbidden si el error indica que el registro no existe o no pertenece al usuario
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al actualizar los datos del sensor"})
		return
	}

	log.Printf("INFO: [UpdateCtrl] Datos actualizados exitosamente para ID %d (por UserID: %d).", id, userID)
	c.JSON(http.StatusOK, gin.H{"message": "Datos del sensor actualizados exitosamente"})
}