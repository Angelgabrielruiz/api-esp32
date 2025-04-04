package main

import (
	"API/src/core" // Conexión DB
	"net/url"      // <--- IMPORTAR para url.Parse
	"os"
	// Auth (Login y Registro)
	authApp "API/src/Sensores/application"
	authInfra "API/src/Sensores/infraestructure"
	authMW "API/src/Sensores/infraestructure/middleware"
	// Sensores
	sensoresInfra "API/src/Sensores/infraestructure"
	infraWS "API/src/Sensores/infraestructure/websocket"
	// Users
	userDomain "API/src/Sensores/domain"
	userAdapters "API/src/Sensores/infraestructure/adapters"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// --- Cargar .env ---
	err := godotenv.Load()
	if err != nil {
		log.Printf("Advertencia: No se pudo cargar el archivo .env: %v", err)
	}

	// --- Conexión a Base de Datos ---
	dbConn := core.GetDBPool()
	if dbConn.Err != "" {
		log.Fatalf("CRÍTICO: Error inicial al conectar con la BD: %s", dbConn.Err)
	}
	if dbConn.DB == nil {
		log.Fatal("CRÍTICO: GetDBPool devolvió una conexión nula sin error.")
	}
	defer dbConn.Close()
	log.Println("INFO: Pool de conexiones MySQL listo.")

	// --- Instanciar Repositorio de Usuarios ---
	var userRepo userDomain.UserRepository
	userRepo = userAdapters.NewMySQLUserRepository(dbConn)
	log.Println("INFO: Repositorio de Usuarios instanciado.")

	// --- Motor Gin ---
	r := gin.Default()

	// --- WebSocket Manager ---
	wsManager := infraWS.NewManager()
	go wsManager.Run()
	log.Println("INFO: WebSocket Manager iniciado.")

	// --- CORS Middleware (Permitiendo "Todos" los Orígenes con Credenciales - WORKAROUND PELIGROSO) ---
	r.Use(cors.New(cors.Config{
		// NO USAR AllowOrigins ni AllowAllOrigins cuando AllowCredentials es true y quieres "permitir todo"
		AllowOriginFunc: func(origin string) bool {
			// Esta función simplemente refleja cualquier origen válido que venga en la petición.
			// ¡NO FILTRA NADA! Esto es inseguro para producción.
			_, err := url.Parse(origin) // Intenta parsear para validar mínimamente el formato
			if err == nil && origin != "null" { // No permitir el origen "null" que a veces envían los navegadores
                log.Printf("DEBUG CORS (AllowOriginFunc): Permitiendo origen %s", origin)
				return true
			}
            log.Printf("DEBUG CORS (AllowOriginFunc): Rechazando origen %s", origin)
			return false // Rechaza orígenes malformados o "null"
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},        // Métodos permitidos
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"}, // Cabeceras permitidas en la petición
		ExposeHeaders:    []string{"Content-Length"},                                  // Cabeceras expuestas en la respuesta
		AllowCredentials: true,                                                        // <-- PERMITE CREDENCIALES (necesario para JWT en header)
		// MaxAge:           12 * time.Hour,                                              // Opcional: Tiempo de caché para preflight
	}))
	log.Println("INFO: Middleware CORS configurado (usando AllowOriginFunc para simular '*' con credenciales - ¡PRECAUCIÓN!).")
	// -------------------------------------------------------------------------------------------------------


	// --- Instanciar Componentes de Autenticación, Registro y Asignación MAC ---
	loginUseCase := authApp.NewLoginUseCase(userRepo)
	loginController := authInfra.NewLoginController(*loginUseCase)
	createUserUseCase := authApp.NewCreateUserUseCase(userRepo)
	createUserController := authInfra.NewCreateUserController(*createUserUseCase)
	assignMacUseCase := authApp.NewAssignMacToUserUseCase(userRepo)
	assignMacController := authInfra.NewAssignMacController(*assignMacUseCase)
	authMiddleware := authMW.JWTMiddleware()
	log.Println("INFO: Componentes de Autenticación, Registro y Admin listos.")


	// --- Configurar Rutas Públicas (Auth: Login y Registro) ---
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", loginController.Execute)
		authGroup.POST("/register", createUserController.Execute)
	}
	log.Println("INFO: Rutas públicas /auth (login, register) configuradas.")


	// --- Configurar Rutas de Módulos (Sensores) ---
	// Pasa las dependencias necesarias, incluyendo el userRepo y el middleware
	sensoresInfra.SetupRoutesDatos(r, wsManager, dbConn, userRepo, authMiddleware)


	// --- Configurar Rutas de Administración ---
	adminGroup := r.Group("/admin")
	adminGroup.Use(authMiddleware) // Protegido por JWT
	{
		adminGroup.PUT("/users/:userId/assign-mac", assignMacController.Execute)
	}
	log.Println("INFO: Rutas de administración /admin configuradas y protegidas por JWT.")


	// --- Ruta WebSocket ---
	r.GET("/ws", func(c *gin.Context) { wsManager.HandleConnections(c.Writer, c.Request) })
	log.Println("INFO: Ruta /ws configurada.")


	// --- Iniciar Servidor ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Servidor iniciando en http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error fatal al iniciar el servidor Gin: %v", err)
	}
}