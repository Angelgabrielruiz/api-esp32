// File: deleteDatos_controller.go

package infraestructure

import (
	"API/src/Sensores/application"
	"net/http"
	"strconv"
	"log"

	"github.com/gin-gonic/gin"
)

type DeleteDatosController struct {
	useCase application.DeleteDatos
}

func NewDeleteDatosController(useCase application.DeleteDatos) *DeleteDatosController {
	return &DeleteDatosController{useCase: useCase}
}

func (dsc *DeleteDatosController) Execute(c *gin.Context) {
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
		log.Printf("ERROR: [DeleteCtrl] ID inválido en la URL: '%s'", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido proporcionado"})
		return
	}

	// Ejecutar el caso de uso, pasando userID para validación
	// Asume que DeleteUseCase.Execute ahora toma userID como parámetro
	err = dsc.useCase.Execute(id /*, userID */) // Pasa userID si el caso de uso/repo lo usa
	if err != nil {
			log.Printf("ERROR: [DeleteCtrl] Falló la ejecución del caso de uso DeleteDatos (ID: %d, UserID: %d): %v", id, userID, err)
			// Podrías devolver 404 Not Found o 403 Forbidden si el error indica que no existe o no pertenece al usuario
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al eliminar los datos"})
		return
	}

	log.Printf("INFO: [DeleteCtrl] Solicitud de eliminación procesada para ID %d (por UserID: %d).", id, userID)
	// Devolver 200 OK o 204 No Content en caso de éxito
	c.JSON(http.StatusOK, gin.H{"message": "Datos eliminados exitosamente"})
	// Alternativa: c.Status(http.StatusNoContent)
}