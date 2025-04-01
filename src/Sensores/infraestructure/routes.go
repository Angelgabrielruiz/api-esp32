package infraestructure

import (
	"API/src/Sensores/application"
	"API/src/Sensores/infraestructure/adapters"
	infraWS "API/src/Sensores/infraestructure/websocket" // Importar tipo del manager
    "log"

	"github.com/gin-gonic/gin"
)

// SetupRoutesDatos configura todas las dependencias y rutas para el módulo Sensores.
// Esta función actúa como el "Composition Root" para este módulo.
func SetupRoutesDatos(r *gin.Engine, wsManager *infraWS.Manager) {

	log.Println("INFO: Configurando rutas y dependencias para Sensores...")

	// --- 1. Crear Adaptadores (Implementaciones de Puertos) ---
	// Adaptador para la Base de Datos (conecta al puerto domain.DatosRepository)
	dbAdapter := adapters.NewMySQLRutas()
	log.Println("INFO: Adaptador MySQL creado.")

	// Adaptador para Notificaciones WebSocket (conecta al puerto application.DatosNotifier)
	wsNotifierAdapter := adapters.NewWebSocketNotifier(wsManager) // Inyecta el WS Manager
	log.Println("INFO: Adaptador WebSocketNotifier creado.")


	// --- 2. Crear Casos de Uso (Lógica de Aplicación) ---
	// Inyectar los adaptadores necesarios a cada caso de uso a través de sus interfaces.
	createDatosUseCase := application.NewCreateDatos(dbAdapter, wsNotifierAdapter)
	getDatosUseCase := application.NewGetDatos(dbAdapter)
	updateDatosUseCase := application.NewUpdateDatos(dbAdapter)
	deleteDatosUseCase := application.NewDeleteDatos(dbAdapter)
	log.Println("INFO: Casos de uso creados e inyectados.")


	// --- 3. Crear Controladores (Manejadores HTTP) ---
	// Inyectar los casos de uso correspondientes a cada controlador.
	createDatosController := NewCreateDatosController(*createDatosUseCase)
	getDatosController := NewGetDatosController(*getDatosUseCase)
	updateDatosController := NewUpdateDatosController(*updateDatosUseCase)
	deleteDatosController := NewDeleteDatosController(*deleteDatosUseCase)
	log.Println("INFO: Controladores HTTP creados.")


	// --- 4. Definir Rutas HTTP ---
	// Agrupar rutas bajo un prefijo (opcional pero bueno para organización)
	// apiGroup := r.Group("/api/v1") // Ejemplo con versionado
	datosGroup := r.Group("/datos") // Grupo para las rutas de datos
	{
		datosGroup.POST("", createDatosController.Execute)      // POST /datos
		datosGroup.GET("", getDatosController.Execute)          // GET /datos
		datosGroup.PUT("/:id", updateDatosController.Execute)   // PUT /datos/:id
		datosGroup.DELETE("/:id", deleteDatosController.Execute) // DELETE /datos/:id
	}
	log.Println("INFO: Rutas HTTP para /datos configuradas.")

    // Nota: La ruta /ws se configura en main.go
}