package main

import (
	"API/src/core" // Importar core para la conexión
	sensoresInfra "API/src/Sensores/infraestructure"
	infraWS "API/src/Sensores/infraestructure/websocket"
	"log"
    // Importar tu middleware
    // auth "API/src/Auth/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Advertencia: No se pudo cargar el archivo .env: %v", err)
	}

	// --- Conexión a Base de Datos ---
	dbConn := core.GetDBPool() // Obtiene la instancia de conexión
    if dbConn.Err != "" {
        log.Fatalf("CRÍTICO: Error inicial al conectar con la BD: %s", dbConn.Err)
    }
    if dbConn.DB == nil {
         log.Fatal("CRÍTICO: GetDBPool devolvió una conexión nula sin error.")
    }
	defer dbConn.Close() // Asegura cerrar la conexión al final
	log.Println("INFO: Pool de conexiones MySQL listo.")


	r := gin.Default()

	// --- WebSocket Manager ---
	wsManager := infraWS.NewManager()
	go wsManager.Run()
	log.Println("INFO: WebSocket Manager iniciado.")


	// --- CORS ---
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Cambia a tu dominio de frontend en producción
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	log.Println("INFO: Middleware CORS configurado.")

    // --- Autenticación Middleware ---
    // authMiddleware := auth.JWTMiddleware() // Necesitas crear esta función/paquete
    // log.Println("INFO: Middleware de autenticación listo.")


	// --- Configurar Rutas de Módulos ---
    // Pasar la conexión DB y el middleware a la configuración de rutas
	sensoresInfra.SetupRoutesDatos(r, wsManager, dbConn /*, authMiddleware */)


	// --- Ruta WebSocket --- (Probablemente necesite Auth también)
	r.GET("/ws", func(c *gin.Context) {
        // Aquí deberías validar el token JWT (quizás pasado como query param en la conexión inicial WS)
        // y luego asociar el userID a la conexión en el wsManager antes de registrarla.
        // Por ahora, lo dejamos como estaba:
		wsManager.HandleConnections(c.Writer, c.Request)
	})
    log.Println("INFO: Ruta /ws configurada.")

	// --- Iniciar Servidor ---
	port := ":8080"
	log.Printf("Servidor iniciando en http://localhost%s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Error fatal al iniciar el servidor Gin: %v", err)
	}
}