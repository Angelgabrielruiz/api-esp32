//File: login_controller.go

package infraestructure

import (
	"API/src/Sensores/application" // Importa el caso de uso                // Para comparar errores
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginController maneja las peticiones HTTP para el login
type LoginController struct {
	loginUseCase application.LoginUseCase
}

// NewLoginController crea una instancia del controlador de login
func NewLoginController(uc application.LoginUseCase) *LoginController {
	return &LoginController{loginUseCase: uc}
}

// Estructura para recibir el JSON del request de login
type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Execute es el manejador de Gin para la ruta POST /auth/login
func (ctrl *LoginController) Execute(c *gin.Context) {
	var req loginRequest
	// Validar que el JSON tenga username y password
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ERROR: [LoginCtrl] JSON de login inválido: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Usuario y contraseña son requeridos"})
		return
	}

	credentials := application.LoginCredentials{
		Username: req.Username,
		Password: req.Password,
	}

	// Ejecutar el caso de uso de login
	result, err := ctrl.loginUseCase.Execute(credentials)
	if err != nil {
		// Verificar el tipo de error devuelto por el caso de uso
		if err.Error() == "credenciales inválidas" {
			// Error esperado por credenciales incorrectas -> 401 Unauthorized
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			// Cualquier otro error se considera interno -> 500 Internal Server Error
			log.Printf("ERROR: [LoginCtrl] Error interno durante el login: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		}
		return
	}

	// Si no hubo error, login exitoso, devolver el token
	c.JSON(http.StatusOK, result) // Devuelve {"token": "ey..."}
}