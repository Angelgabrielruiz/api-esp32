package infraestructure

import (
	"API/src/Sensores/application" // Depende solo de la capa de aplicación
	"log"
	"net/http"
	"strings" // Para validación simple

	"github.com/gin-gonic/gin"
)

// CreateDatosController maneja las solicitudes HTTP para crear datos.
// Depende del caso de uso correspondiente.
type CreateDatosController struct {
	useCase application.CreateDatos // Referencia al caso de uso
}

// NewCreateDatosController crea una instancia del controlador.
func NewCreateDatosController(useCase application.CreateDatos) *CreateDatosController {
	// No necesitamos validar useCase aquí porque los constructores de use case ya lo hacen.
	return &CreateDatosController{useCase: useCase}
}

// CreateDatosRequest define la estructura esperada en el cuerpo JSON de la solicitud POST.
// Usar `binding:"required"` para validación automática de Gin.
type CreateDatosRequest struct {
	Temperatura string `json:"temperatura" binding:"required"`
	Movimiento  string `json:"movimiento" binding:"required"`
	Distancia   string `json:"distancia" binding:"required"`
	Peso        string `json:"peso" binding:"required"`
	Mac         string `json:"mac" binding:"required"`
}

// Execute es el manejador de Gin para la ruta POST /datos.
func (csc *CreateDatosController) Execute(c *gin.Context) {
	var requestBody CreateDatosRequest

	// Validar y parsear el cuerpo JSON de la solicitud.
	// ShouldBindJSON ya incluye validación si usas `binding:"required"`.
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("ERROR: [CreateCtrl] Datos inválidos en la solicitud: %v", err)
		// Devolver un error claro al cliente.
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Datos inválidos o faltantes en la solicitud",
			"detail": err.Error(), // Proporciona detalles del error de validación
		})
		return
	}

    // Validación adicional (opcional, podría estar en el use case o dominio)
    if strings.TrimSpace(requestBody.Temperatura) == "" || /* otras validaciones */ false {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Valores no pueden estar vacíos"})
        return
    }



	err := csc.useCase.Execute(
		requestBody.Temperatura,
		requestBody.Movimiento,
		requestBody.Distancia,
		requestBody.Peso,
		requestBody.Mac,  // Fixed: removed "mac:" prefix
	)

	// Manejar el resultado del caso de uso.
	if err != nil {
		// Si el caso de uso falla (probablemente al guardar en DB), retornar error interno.
		log.Printf("ERROR: [CreateCtrl] Falló la ejecución del caso de uso CreateDatos: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno al registrar los datos del sensor",
			// "detail": err.Error(), // Considera no exponer errores internos detallados al cliente final
		})
		return
	}

	// Si el caso de uso se completó sin error, responder con éxito.
	log.Printf("INFO: [CreateCtrl] Datos registrados exitosamente: %+v", requestBody)
	c.JSON(http.StatusCreated, gin.H{"message": "Datos del sensor registrados exitosamente"}) // Usar 201 Created
}