//File: routes.go

package infraestructure

import (
	"API/src/core"                                          // Importar core para la conexión
	sensorApp "API/src/Sensores/application"                // Alias para claridad
	sensorAdapters "API/src/Sensores/infraestructure/adapters" // Alias
	infraWS "API/src/Sensores/infraestructure/websocket"
	userDomain "API/src/Sensores/domain" // Importar el paquete que define UserRepository
	// La dependencia de userAdapters puede ser necesaria aquí si se instancia aquí
	// o si se pasa el repo ya creado desde main.go
	"log"

	"github.com/gin-gonic/gin" // Necesario para gin.HandlerFunc
)

// SetupRoutesDatos configura las rutas para Sensores, AHORA recibe el middleware de Auth
func SetupRoutesDatos(r *gin.Engine, wsManager *infraWS.Manager, dbConn *core.Conn_MySQL, userRepo userDomain.UserRepository, authMiddleware gin.HandlerFunc) {

	log.Println("INFO: Configurando rutas y dependencias para Sensores...")

	if dbConn == nil || dbConn.DB == nil {
		log.Fatal("CRÍTICO: SetupRoutesDatos recibió una conexión DB nula.")
	}
	if userRepo == nil {
		log.Fatal("CRITICO: SetupRoutesDatos recibió un userRepo nulo.")
	}
	if authMiddleware == nil {
		log.Fatal("CRITICO: SetupRoutesDatos recibió un authMiddleware nulo.")
	}


	// --- 1. Crear Adaptadores ---
	dbSensorAdapter := sensorAdapters.NewMySQLRutas(dbConn)
	log.Println("INFO: Adaptador MySQL para Sensores creado.")

	// userRepo ya viene inyectado desde main.go

	wsNotifierAdapter := sensorAdapters.NewWebSocketNotifier(wsManager)
	log.Println("INFO: Adaptador WebSocketNotifier creado.")

	// --- 2. Crear Casos de Uso ---
	// CreateDatos necesita el userRepo (que ya recibimos)
	createDatosUseCase := sensorApp.NewCreateDatos(dbSensorAdapter, userRepo, wsNotifierAdapter)
	getDatosUseCase := sensorApp.NewGetDatos(dbSensorAdapter)
	updateDatosUseCase := sensorApp.NewUpdateDatos(dbSensorAdapter) // Podría necesitar userRepo si valida pertenencia
	deleteDatosUseCase := sensorApp.NewDeleteDatos(dbSensorAdapter) // Podría necesitar userRepo si valida pertenencia
	log.Println("INFO: Casos de uso de Sensores creados e inyectados.")

	// --- 3. Crear Controladores ---
	createDatosController := NewCreateDatosController(*createDatosUseCase)
	getDatosController := NewGetDatosController(*getDatosUseCase)
	updateDatosController := NewUpdateDatosController(*updateDatosUseCase)
	deleteDatosController := NewDeleteDatosController(*deleteDatosUseCase)
	log.Println("INFO: Controladores HTTP de Sensores creados.")

	// --- 4. Definir Rutas HTTP ---
	// Endpoint de ingesta (llamado por el consumidor, usualmente no requiere auth de usuario final)
	sensorDataIngestPath := "/api/sensor-data"
	r.POST(sensorDataIngestPath, createDatosController.Execute)
	log.Printf("INFO: Ruta POST %s configurada para ingesta de datos (sin auth JWT usuario).", sensorDataIngestPath)

	// Grupo para las rutas del FRONTEND (protegidas por JWT)
	datosGroup := r.Group("/datos")
	datosGroup.Use(authMiddleware) // <--- APLICAR MIDDLEWARE JWT A ESTE GRUPO
	{
		datosGroup.GET("", getDatosController.Execute)          // Protegido
		datosGroup.PUT("/:id", updateDatosController.Execute)   // Protegido
		datosGroup.DELETE("/:id", deleteDatosController.Execute) // Protegido

		// Opcional: Ruta admin (si la implementas)
		// datosGroup.GET("/all", getDatosController.ExecuteAll) // Necesitaría check de rol adicional
	}
	log.Println("INFO: Rutas HTTP para /datos (frontend) configuradas y protegidas por JWT.")
}