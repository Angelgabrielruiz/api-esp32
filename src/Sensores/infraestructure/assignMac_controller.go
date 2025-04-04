package infraestructure // O el paquete de infraestructura que uses

import (
	"API/src/Sensores/application" // Ruta a tu paquete application
	"database/sql"                 // Para sql.ErrNoRows
	//"fmt"
	"log"
	"net/http"
	"strconv"
	//"strings"

	"github.com/gin-gonic/gin"
)

// AssignMacController maneja las peticiones HTTP para asignar MAC
type AssignMacController struct {
	useCase application.AssignMacToUserUseCase // Dependencia del caso de uso
}

// NewAssignMacController crea la instancia
func NewAssignMacController(uc application.AssignMacToUserUseCase) *AssignMacController {
	return &AssignMacController{useCase: uc}
}

// assignMacRequest define el cuerpo JSON esperado
type assignMacRequest struct {
	// Permitir cadena vacía para desasignar, `binding:"required"` fallaría
	MacAddress string `json:"mac_address"`
}

// Execute es el manejador Gin para la ruta PUT /admin/users/:userId/assign-mac
func (ctrl *AssignMacController) Execute(c *gin.Context) {
	// 1. Verificar si el usuario es Administrador
	userRoleValue, _ := c.Get("userRole") // Obtiene rol del middleware JWT
	userRole, _ := userRoleValue.(string)
	if userRole != "admin" { // <-- CHEQUEO DE ROL
		log.Printf("WARN: [AssignMacCtrl] Intento de acceso no autorizado por rol: '%s'", userRole)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Acceso denegado: requiere rol de administrador"})
		return
	}

	// 2. Obtener el ID del usuario objetivo de la URL
	userIdParam := c.Param("userId") // El :userId de la ruta
	targetUserID, err := strconv.Atoi(userIdParam)
	if err != nil || targetUserID <= 0 {
		log.Printf("ERROR: [AssignMacCtrl] UserID inválido en URL: '%s'", userIdParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido en la ruta"})
		return
	}

	// 3. Obtener la nueva MAC del cuerpo JSON
	var req assignMacRequest
	// Usar BindJSON en lugar de ShouldBindJSON si MacAddress puede estar ausente u omitida
	if err := c.ShouldBindJSON(&req); err != nil {
        // Podría ser un JSON mal formado o faltar el campo si usas binding:required
		log.Printf("ERROR: [AssignMacCtrl] JSON inválido en la petición: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de la petición inválido o falta 'mac_address'"})
		return
	}
    // Nota: req.MacAddress puede ser "" si el JSON es `{"mac_address": ""}` o si se omite el campo y no hay `binding:"required"`

	// 4. Preparar y ejecutar el caso de uso
	input := application.AssignMacInput{
		TargetUserID: targetUserID,
		MacAddress:   req.MacAddress, // Pasamos el valor recibido (puede ser vacío)
	}
	err = ctrl.useCase.Execute(input)

	// 5. Manejar la respuesta basada en el error del caso de uso
	if err != nil {
		log.Printf("ERROR: [AssignMacCtrl] Error del caso de uso al asignar MAC a UserID %d: %v", targetUserID, err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		} else if err.Error() == "mac_address_duplicado" { // Error específico del repo
			c.JSON(http.StatusConflict, gin.H{"error": "La dirección MAC ya está asignada a otro usuario"})
		} else if err.Error() == "formato_mac_invalido" { // Error específico del use case
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de dirección MAC inválido"})
		} else {
			// Otro error (probablemente DB o interno)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al asignar la MAC"})
		}
		return
	}

	// Éxito
	log.Printf("INFO: [AssignMacCtrl] MAC '%s' asignada/actualizada para UserID %d por admin.", req.MacAddress, targetUserID)
	c.JSON(http.StatusOK, gin.H{"message": "Dirección MAC asignada/actualizada exitosamente"})
}