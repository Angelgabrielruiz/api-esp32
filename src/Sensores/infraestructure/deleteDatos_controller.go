//File: deleteDatos_controller.go

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
	// Obtener ID del parámetro de la URL
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 { // Validar que sea un entero positivo
		log.Printf("ERROR: [DeleteCtrl] ID inválido en la URL: '%s'", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido proporcionado"})
		return
	}

	// Ejecutar el caso de uso
	err = dsc.useCase.Execute(id)
	if err != nil {
			log.Printf("ERROR: [DeleteCtrl] Falló la ejecución del caso de uso DeleteDatos (ID: %d): %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al eliminar los datos"})
		return
	}

	log.Printf("INFO: [DeleteCtrl] Solicitud de eliminación procesada para ID %d.", id)
	c.JSON(http.StatusOK, gin.H{"message": "Datos eliminados exitosamente"})
}