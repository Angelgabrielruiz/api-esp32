//File: createDatos_controller.go

package infraestructure

import (
	"API/src/Sensores/application" // Depende solo de la capa de aplicación
	"strings"                          // Para formatear errores
	"log"
	"net/http"

	//"strings" // Para validación simple si es necesario

	"github.com/gin-gonic/gin"
)

type CreateDatosController struct {
	useCase application.CreateDatos // Referencia al caso de uso
}

func NewCreateDatosController(useCase application.CreateDatos) *CreateDatosController {
	return &CreateDatosController{useCase: useCase}
}

// La request ahora solo necesita los campos que vienen del ESP32/Consumidor
type CreateDatosRequest struct {
	Temperatura string `json:"temperatura"` // Quitar binding:"required" si algunos pueden faltar
	Movimiento  string `json:"movimiento"`
	Distancia   string `json:"distancia"`
	Peso        string `json:"peso"`
	Mac         string `json:"mac"` // ¡Esencial!
}

// Este endpoint será llamado por tu CONSUMIDOR
func (csc *CreateDatosController) Execute(c *gin.Context) {
	var requestBody CreateDatosRequest

	// Parsear el cuerpo JSON. ShouldBindJSON es suficiente.
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("ERROR: [CreateCtrl] Datos inválidos en la solicitud del consumidor: %v. Body: %s", err, c.Request.Body)
		// Error 400 indica que el consumidor envió algo malformado
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Payload JSON inválido o incompleto recibido del consumidor",
			"detail": err.Error(),
		})
		return
	}

	// Validar que la MAC no esté vacía (importante!)
	if requestBody.Mac == "" {
		log.Printf("ERROR: [CreateCtrl] Payload recibido sin MAC address: %+v", requestBody)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Falta la dirección MAC en el payload"})
		return
	}

	// Llamar al caso de uso pasando los datos recibidos
	err := csc.useCase.Execute(
		requestBody.Temperatura,
		requestBody.Movimiento,
		requestBody.Distancia,
		requestBody.Peso,
		requestBody.Mac, // Pasar la MAC
	)

	if err != nil {
		// Analizar el tipo de error devuelto por el caso de uso
		if strings.HasPrefix(err.Error(), "mac_no_asignada:") {
			// MAC válida pero no asignada. Esto no es un error del servidor.
			// Respondemos 200 OK o 202 Accepted al consumidor para que haga ACK,
			// pero informamos en el log o cuerpo de respuesta (opcional).
			log.Printf("INFO: [CreateCtrl] Datos de MAC no asignada (%s) descartados como esperado.", requestBody.Mac)
			c.JSON(http.StatusOK, gin.H{"message": "Datos recibidos pero MAC no asignada a un usuario.", "mac": requestBody.Mac})
			// O simplemente: c.Status(http.StatusNoContent) // 204
		} else {
			// Otro error (problema de DB, etc.) -> Error 500
			log.Printf("ERROR: [CreateCtrl] Falló la ejecución del caso de uso CreateDatos: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error interno al procesar los datos del sensor",
			})
		}
		return
	}

	// Éxito: el caso de uso guardó y notificó (o lo intentó)
	log.Printf("INFO: [CreateCtrl] Datos procesados exitosamente para MAC: %s", requestBody.Mac)
	// 201 Created es apropiado si se creó un recurso nuevo
	c.JSON(http.StatusCreated, gin.H{"message": "Datos del sensor procesados exitosamente"})
}