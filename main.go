package main

import (
	sensoresInfra "API/src/Sensores/infraestructure"         // Importa el paquete de infraestructura de Sensores
	infraWS "API/src/Sensores/infraestructure/websocket" // Importa el paquete del WebSocket manager
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		
		log.Printf("Advertencia: No se pudo cargar el archivo .env: %v", err)
	}


	r := gin.Default()


	wsManager := infraWS.NewManager() 
	go wsManager.Run()                

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))


	sensoresInfra.SetupRoutesDatos(r, wsManager)

	// Define la ruta específica para las conexiones WebSocket
	r.GET("/ws", func(c *gin.Context) {
		// Delega el manejo de la conexión WebSocket al manager
		wsManager.HandleConnections(c.Writer, c.Request)
	})

	// --- Iniciar Servidor ---
	port := ":8080" 
	log.Printf("Servidor iniciando en el puerto %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Error fatal al iniciar el servidor: %v", err)
	}
}