package infraestructure // O el paquete de infraestructura que uses

import (
	"API/src/Sensores/application" // Asegúrate que la ruta al paquete application sea correcta
	"fmt"                          // Para errores
	"log"
	"net/http"
	"strings" // Para manejo de errores

	"github.com/gin-gonic/gin"
)

// CreateUserController maneja las peticiones HTTP para crear usuarios
type CreateUserController struct {
	createUserUseCase application.CreateUserUseCase // Dependencia del caso de uso
}

// NewCreateUserController crea una instancia del controlador
func NewCreateUserController(uc application.CreateUserUseCase) *CreateUserController {
	return &CreateUserController{createUserUseCase: uc}
}

// Estructura para el cuerpo JSON de la petición de registro
type createUserRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	MacAddress string `json:"mac_address"` // MAC es opcional en el request
	Role       string `json:"role"`        // Rol es opcional, el use case pondrá 'user' por defecto
}

// Execute es el manejador Gin para la ruta POST /auth/register
func (ctrl *CreateUserController) Execute(c *gin.Context) {
	var req createUserRequest

	// 1. Validar y parsear el JSON del cuerpo
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ERROR: [CreateUserCtrl] JSON de registro inválido: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de registro inválidos o incompletos", "details": err.Error()})
		return
	}

	// 2. Preparar el DTO de entrada para el caso de uso
	input := application.CreateUserInput{
		Username:   req.Username,
		Password:   req.Password,
		MacAddress: req.MacAddress, // Pasa la MAC (puede ser vacía)
		Role:       req.Role,       // Pasa el Rol (puede ser vacío)
	}

	// 3. Ejecutar el caso de uso
	err := ctrl.createUserUseCase.Execute(input)
	if err != nil {
		// Mapear errores del caso de uso a respuestas HTTP
		if strings.HasPrefix(err.Error(), "la contraseña debe") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if strings.HasPrefix(err.Error(), "conflicto:") {
			// Extraer el campo duplicado si es posible (ej: "username_duplicado")
			fieldName := "recurso"
			if strings.Contains(err.Error(), "username_duplicado") {
				fieldName = "nombre de usuario"
			} else if strings.Contains(err.Error(), "mac_address_duplicado"){
                fieldName = "dirección MAC"
            }
			errMsg := fmt.Sprintf("El %s ya está en uso.", fieldName)
			c.JSON(http.StatusConflict, gin.H{"error": errMsg}) // 409 Conflict
		} else {
			// Errores internos genéricos
			log.Printf("ERROR: [CreateUserCtrl] Error interno al ejecutar caso de uso: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor al crear usuario"})
		}
		return
	}

	// 4. Éxito
	log.Printf("INFO: [CreateUserCtrl] Usuario '%s' registrado exitosamente.", req.Username)
	// Devolver 201 Created
	c.JSON(http.StatusCreated, gin.H{"message": "Usuario registrado exitosamente", "username": req.Username})
}