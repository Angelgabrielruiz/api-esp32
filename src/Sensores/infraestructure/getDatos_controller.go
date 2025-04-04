// File: src/Sensores/infraestructure/getDatos_controller.go

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
	// --- OBTENER USER ID DEL CONTEXTO (Puesto por JWTMiddleware) ---
	userIDValue, exists := c.Get("userID") // Usa la clave "userID"
	if !exists {
		log.Printf("ERROR: [GetCtrl] No se encontró userID en el contexto. Middleware de Auth ausente o falló?")
		// No debería pasar si el middleware está bien configurado y aplicado
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado (contexto inválido)"})
		return
	}

	userID, ok := userIDValue.(int) // Asegúrate que el tipo sea int
	if !ok || userID <= 0 {
		log.Printf("ERROR: [GetCtrl] userID en contexto tiene tipo inválido (%T) o valor no positivo.", userIDValue)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno de autorización"}) // Podría ser 500 si el tipo es erróneo
		return
	}
	// --- FIN OBTENER USER ID ---

	// Pasar el userID al caso de uso para filtrar
	datos, err := gdc.useCase.Execute(userID) // Llama al caso de uso con el ID
	if err != nil {
		log.Printf("ERROR: [GetCtrl] Falló la ejecución del caso de uso GetDatos para UserID %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al obtener los datos del sensor"})
		return
	}

	// Si no hay datos, devolver un array vacío en lugar de null
	if datos == nil {
		log.Println("INFO: [GetCtrl] No se encontraron datos para UserID %d. Devolviendo array vacío.", userID)
		datos = []entities.Datos{}
	}

	log.Printf("INFO: [GetCtrl] Devolviendo %d registros para UserID %d.", len(datos), userID)
	c.JSON(http.StatusOK, datos) // Devuelve los datos filtrados
}

// Opcional: Endpoint para Admin (si lo implementaste en el use case)
/*
func (gdc *GetDatosController) ExecuteAll(c *gin.Context) {
    // Aquí verificarías el rol del usuario obtenido del contexto:
    // userRoleValue, _ := c.Get("userRole")
    // userRole, _ := userRoleValue.(string)
    // if userRole != "admin" {
    //     c.JSON(http.StatusForbidden, gin.H{"error": "Acceso denegado: requiere rol de administrador"})
    //     return
    // }

    datos, err := gdc.useCase.ExecuteAll() // Llama al método sin filtro
    if err != nil {
        log.Printf("ERROR: [GetCtrl] Falló la ejecución del caso de uso ExecuteAll (admin): %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al obtener todos los datos"})
        return
    }
     if datos == nil {
        datos = []entities.Datos{}
    }
    log.Printf("INFO: [GetCtrl] Devolviendo %d registros en total (admin).", len(datos))
    c.JSON(http.StatusOK, datos)
}
*/