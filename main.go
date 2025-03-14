package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/gin-contrib/cors"
    sensoresInfra "API/src/Sensores/infraestructure"
)

func main() {
    // Cargar variables de entorno
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error al cargar el archivo .env: %v", err)
    }

    // Crear servidor Gin
    r := gin.Default()

    // Habilitar CORS para aceptar solicitudes del ESP32
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"}, // Permitir todas las conexiones (puedes cambiarlo a tu IP específica)
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        AllowCredentials: true,
    }))

    // Configurar rutas de sensores
    sensoresInfra.SetupRoutesDatos(r)

    // Iniciar el servidor
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Error al iniciar el servidor: %v", err)
    }
}