package infraestructure

import (
	//"Send/src/Pagos/application"
	//"Send/src/Pagos/infraestructure/adapters"
	"API/src/Sensores/application"
	"API/src/Sensores/infraestructure/adapters"

	"github.com/gin-gonic/gin"
)

func SetupRoutesDatos(r *gin.Engine) {

	
	ps := adapters.NewMySQLRutas()

	
	createDatosController := NewCreateDatosController(*application.NewCreateDatos(ps))
	getDatosController := NewGetDatosController(*application.NewGetDatos(ps))
	updateDatosController := NewUpdateDatosController(*application.NewUpdateDatos(ps))
	deleteDatosController := NewDeleteDatosController(*application.NewDeleteDatos(ps))

	
	r.POST("/datos", createDatosController.Execute)
	r.GET("/datos", getDatosController.Execute)
	r.PUT("/datos/:id", updateDatosController.Execute)
	r.DELETE("/datos/:id", deleteDatosController.Execute)
}
