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
	// Ejecutar el caso de uso para obtener los datos
	datos, err := gdc.useCase.Execute()
	if err != nil {
		log.Printf("ERROR: [GetCtrl] Falló la ejecución del caso de uso GetDatos: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al obtener los datos del sensor"})
		return
	}

	// Si no hay datos, devolver un array vacío, no un error.
	// La capa de aplicación/repositorio ya debería manejar el caso "no encontrado" devolviendo un slice vacío.
	if datos == nil {
        // Esto no debería ocurrir si el use case devuelve slice vacío en lugar de nil
        log.Println("ADVERTENCIA: [GetCtrl] El caso de uso GetDatos devolvió nil en lugar de un slice vacío.")
        datos = []entities.Datos{} // Asegurar que siempre sea un array JSON
    }

	log.Printf("INFO: [GetCtrl] Devolviendo %d registros.", len(datos))
	// Devolver los datos como JSON. Gin se encarga de la serialización.
	c.JSON(http.StatusOK, datos)
}