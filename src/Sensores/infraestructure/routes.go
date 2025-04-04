package infraestructure

import (
	"API/src/core"                                          // Importar core para la conexión
	sensorApp "API/src/Sensores/application"                // Alias para claridad
	sensorAdapters "API/src/Sensores/infraestructure/adapters" // Alias
	infraWS "API/src/Sensores/infraestructure/websocket"
	userAdapters "API/src/Sensores/infraestructure/adapters" // Importar adaptador de usuarios
	"log"

    // Importar tu middleware de autenticación
    // "API/src/Auth/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutesDatos(r *gin.Engine, wsManager *infraWS.Manager, dbConn *core.Conn_MySQL /*, authMiddleware gin.HandlerFunc */) { // Recibe conexión y middleware Auth

	log.Println("INFO: Configurando rutas y dependencias para Sensores...")

    if dbConn == nil || dbConn.DB == nil {
         log.Fatal("CRÍTICO: SetupRoutesDatos recibió una conexión DB nula.")
    }

	// --- 1. Crear Adaptadores ---
	// Adaptador DB Sensores (AHORA recibe la conexión)
	dbSensorAdapter := sensorAdapters.NewMySQLRutas(dbConn)
	log.Println("INFO: Adaptador MySQL para Sensores creado.")

	// Adaptador DB Usuarios (NUEVO, también recibe la conexión)
	dbUserAdapter := userAdapters.NewMySQLUserRepository(dbConn)
	log.Println("INFO: Adaptador MySQL para Usuarios creado.")

	// Adaptador Notificaciones WebSocket
	wsNotifierAdapter := sensorAdapters.NewWebSocketNotifier(wsManager)
	log.Println("INFO: Adaptador WebSocketNotifier creado.")

	// --- 2. Crear Casos de Uso ---
	// CreateDatos AHORA necesita userRepo
	createDatosUseCase := sensorApp.NewCreateDatos(dbSensorAdapter, dbUserAdapter, wsNotifierAdapter)
	getDatosUseCase := sensorApp.NewGetDatos(dbSensorAdapter)
	updateDatosUseCase := sensorApp.NewUpdateDatos(dbSensorAdapter)
	deleteDatosUseCase := sensorApp.NewDeleteDatos(dbSensorAdapter)
	log.Println("INFO: Casos de uso creados e inyectados.")

	// --- 3. Crear Controladores ---
	createDatosController := NewCreateDatosController(*createDatosUseCase) // Llamado por el consumidor
	getDatosController := NewGetDatosController(*getDatosUseCase)         // Llamado por el frontend (necesita Auth)
	updateDatosController := NewUpdateDatosController(*updateDatosUseCase) // Llamado por el frontend (necesita Auth)
	deleteDatosController := NewDeleteDatosController(*deleteDatosUseCase) // Llamado por el frontend (necesita Auth)
	log.Println("INFO: Controladores HTTP creados.")

	// --- 4. Definir Rutas HTTP ---
	// Endpoint para que el CONSUMIDOR envíe datos (puede o no necesitar auth)
	// Si necesita auth, sería un token fijo del consumidor, no de un usuario final.
    // Por ahora, asumimos que es un endpoint "interno" o protegido por red.
	sensorDataIngestPath := "/api/sensor-data" // O usa "/datos" si prefieres, pero separa conceptualmente
	r.POST(sensorDataIngestPath, createDatosController.Execute)
    log.Printf("INFO: Ruta POST %s configurada para ingesta de datos.", sensorDataIngestPath)


	// Grupo para las rutas que el FRONTEND consume (protegidas por Auth)
	datosGroup := r.Group("/datos")
    // Aplicar middleware de autenticación a este grupo
    // datosGroup.Use(authMiddleware) // DESCOMENTA cuando tengas el middleware
	{
		datosGroup.GET("", getDatosController.Execute)          // GET /datos (filtrado por usuario logueado)
		datosGroup.PUT("/:id", updateDatosController.Execute)   // PUT /datos/:id (validando usuario logueado)
		datosGroup.DELETE("/:id", deleteDatosController.Execute) // DELETE /datos/:id (validando usuario logueado)

        // Opcional: Ruta admin para ver todo
         // adminGroup := r.Group("/admin/datos")
         // adminGroup.Use(authMiddleware) // Y un check de rol admin
         // {
         //     adminGroup.GET("", getDatosController.ExecuteAll)
         // }
	}
	log.Println("INFO: Rutas HTTP para /datos (frontend) configuradas.")
    // Nota: La ruta /ws se configura en main.go
}