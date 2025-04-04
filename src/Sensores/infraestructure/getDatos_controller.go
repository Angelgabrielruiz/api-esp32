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
	datos, err := gdc.useCase.Execute()
	if err != nil {
		log.Printf("ERROR: [GetCtrl] Falló la ejecución del caso de uso GetDatos: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al obtener los datos del sensor"})
		return
	}

	if datos == nil {
        log.Println("ADVERTENCIA: [GetCtrl] El caso de uso GetDatos devolvió nil en lugar de un slice vacío.")
        datos = []entities.Datos{}
    }

	log.Printf("INFO: [GetCtrl] Devolviendo %d registros.", len(datos))
	c.JSON(http.StatusOK, datos)
}