package infraestructure

import (
	"API/src/Sensores/application"
	"API/src/Sensores/domain/entities"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetDatosController struct {
	useCase application.GetDatos
}

func NewGetDatosController(useCase application.GetDatos) *GetDatosController {
	return &GetDatosController{useCase: useCase}
}

func (gdc *GetDatosController) Execute(c *gin.Context) {
	// --- OBTENER USER ID DEL CONTEXTO (Puesto por el middleware de Auth) ---
	userIDValue, exists := c.Get("userID")
	if !exists {
		log.Printf("ERROR: [GetCtrl] No se encontró userID en el contexto. ¿Falta middleware de Auth?")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado o sesión inválida"})
		return
	}

	userID, ok := userIDValue.(int) // Asegúrate que el tipo sea correcto (int, int64, string?)
	if !ok || userID <= 0 {
		log.Printf("ERROR: [GetCtrl] userID en contexto tiene tipo inválido o es <= 0.")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error de autorización interna"})
		return
	}
	// --- FIN OBTENER USER ID ---

	// Pasar el userID al caso de uso
	datos, err := gdc.useCase.Execute(userID)
	if err != nil {
		log.Printf("ERROR: [GetCtrl] Falló la ejecución del caso de uso GetDatos para UserID %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al obtener los datos del sensor"})
		return
	}

	if datos == nil {
		log.Println("INFO: [GetCtrl] No se encontraron datos para UserID %d. Devolviendo array vacío.", userID)
		datos = []entities.Datos{}
	}

	log.Printf("INFO: [GetCtrl] Devolviendo %d registros para UserID %d.", len(datos), userID)
	c.JSON(http.StatusOK, datos)
}

 // Opcional: Endpoint para Admin (sin filtro de usuario)
 func (gdc *GetDatosController) ExecuteAll(c *gin.Context) {
    // Aquí podrías verificar si el usuario tiene rol de 'admin' usando el contexto
    // if !isAdmin(c) { c.JSON(http.StatusForbidden, ...); return }

    datos, err := gdc.useCase.ExecuteAll()
    if err != nil {
         // ... manejo de error ...
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al obtener todos los datos"})
        return
    }
    c.JSON(http.StatusOK, datos)
}