// File: jwt_middleware.go

package middleware // O el nombre que le des a tu paquete de middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware crea un manejador Gin para verificar tokens JWT
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Obtener la cabecera Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("WARN: [AuthMW] Falta cabecera Authorization")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Se requiere token de autenticación"})
			return
		}

		// 2. Validar formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			log.Println("WARN: [AuthMW] Formato de cabecera Authorization inválido")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			return
		}
		tokenString := parts[1]

		// 3. Obtener clave secreta del entorno
		secretKey := os.Getenv("JWT_SECRET_KEY")
		if secretKey == "" {
			log.Println("ERROR: [AuthMW] JWT_SECRET_KEY no configurada en el servidor")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error de configuración del servidor"})
			return
		}

		// 4. Parsear y validar el token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Valida que el alg sea el esperado (HS256 en este caso)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
			}
			// Devuelve la clave secreta como []byte
			return []byte(secretKey), nil
		})

		// 5. Manejar errores de parseo/validación
		if err != nil {
			log.Printf("WARN: [AuthMW] Error al parsear/validar token: %v", err)
			errMsg := "Token inválido"
			// Ser más específico si el token expiró
			if errors.Is(err, jwt.ErrTokenExpired) {
				errMsg = "Token expirado"
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errMsg})
			return
		}

		// 6. Extraer claims si el token es válido
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extraer user_id (viene como float64 de JSON)
			userIDFloat, okUserID := claims["user_id"].(float64)
			if !okUserID {
				log.Println("ERROR: [AuthMW] user_id no encontrado o tipo inválido en claims JWT")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido (claim user_id)"})
				return
			}
			userID := int(userIDFloat) // Convertir a int

			// Extraer otros claims opcionales (username, role)
			username, _ := claims["username"].(string) // Ignorar error si no está
			role, _ := claims["role"].(string)         // Ignorar error si no está

			// 7. Guardar información en el contexto de Gin
			c.Set("userID", userID)     // Clave "userID"
			c.Set("username", username) // Clave "username"
			c.Set("userRole", role)     // Clave "userRole"

			log.Printf("INFO: [AuthMW] Token válido. Usuario autenticado: ID=%d, Username=%s, Role=%s", userID, username, role)

			// 8. Continuar con el siguiente manejador en la cadena
			c.Next()
		} else {
			// Caso raro donde el token no es válido después de parsear sin error
			log.Println("WARN: [AuthMW] Token JWT inválido (claims no ok o token no válido)")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			return
		}
	}
}