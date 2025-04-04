// File: src/Sensores/application/login_useCase.go

package application

import (
	userDomain "API/src/Sensores/domain" // Repositorio de usuarios
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt" // Para comparar hashes
)

// LoginCredentials DTO para la entrada
type LoginCredentials struct {
	Username string
	Password string
}

// LoginResult DTO para la salida
type LoginResult struct {
	Token string `json:"token"`
}

// LoginUseCase maneja la lógica de negocio del login
type LoginUseCase struct {
	userRepo userDomain.UserRepository
}

// NewLoginUseCase crea una instancia del caso de uso de login
func NewLoginUseCase(userRepo userDomain.UserRepository) *LoginUseCase {
	if userRepo == nil {
		log.Fatal("CRITICO: NewLoginUseCase recibió un userRepo nulo.")
	}
	return &LoginUseCase{userRepo: userRepo}
}

// Función de ayuda para comparar hash de contraseña
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Printf("DEBUG: Error comparando hash: %v", err) // Loguea si falla la comparación
	}
	return err == nil // Devuelve true si coinciden (err es nil)
}

// Execute procesa el intento de login
func (uc *LoginUseCase) Execute(credentials LoginCredentials) (*LoginResult, error) {
	// 1. Buscar usuario por username
	user, err := uc.userRepo.FindByUsername(credentials.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("WARN: [LoginUseCase] Intento de login fallido (usuario no encontrado): %s", credentials.Username)
			// Devolver un error específico que el controlador pueda interpretar como 401
			return nil, errors.New("credenciales inválidas")
		}
		log.Printf("ERROR: [LoginUseCase] Error buscando usuario %s: %v", credentials.Username, err)
		return nil, fmt.Errorf("error interno de autenticación") // Error interno 500
	}

	// 2. Comparar contraseña (hash)
	if !checkPasswordHash(credentials.Password, user.PasswordHash) {
		log.Printf("WARN: [LoginUseCase] Intento de login fallido (contraseña incorrecta): %s", credentials.Username)
		return nil, errors.New("credenciales inválidas") // Error específico 401
	}

	// 3. Generar Token JWT
	tokenString, err := generateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		log.Printf("ERROR: [LoginUseCase] Error generando JWT para usuario %s: %v", credentials.Username, err)
		return nil, fmt.Errorf("error interno al generar sesión") // Error interno 500
	}

	log.Printf("INFO: [LoginUseCase] Login exitoso para usuario: %s (ID: %d)", user.Username, user.ID)
	return &LoginResult{Token: tokenString}, nil
}

// generateJWT crea un nuevo token JWT para el usuario
func generateJWT(userID int, username string, role string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Println("ERROR CRITICO: JWT_SECRET_KEY no está configurada en el entorno.")
		return "", errors.New("JWT_SECRET_KEY no está configurada")
	}

	expirationMinutesStr := os.Getenv("JWT_EXPIRATION_MINUTES")
	expirationMinutes, err := strconv.Atoi(expirationMinutesStr)
	if err != nil || expirationMinutes <= 0 {
		expirationMinutes = 60 // Default a 60 minutos si no está o es inválido
	}
	expirationTime := time.Now().Add(time.Duration(expirationMinutes) * time.Minute)

	// Crear los Claims (Payload del token)
	claims := jwt.MapClaims{
		"user_id":  userID,                // Identificador único del usuario
		"username": username,              // Nombre de usuario (informativo)
		"role":     role,                  // Rol del usuario para autorización
		"exp":      expirationTime.Unix(), // Tiempo de expiración (timestamp Unix)
		"iat":      time.Now().Unix(),     // Tiempo de emisión (timestamp Unix)
	}

	// Crear el token con método de firma HS256 y los claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("error al firmar el token JWT: %w", err)
	}

	return tokenString, nil
}