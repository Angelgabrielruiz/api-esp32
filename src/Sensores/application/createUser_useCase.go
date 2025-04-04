package application // O el paquete de aplicación que corresponda

import (
	userDomain "API/src/Sensores/domain" // Asegúrate que la ruta a tu domain sea correcta
	"database/sql"                     // Para sql.NullString
	"fmt"
	"log"
	"strings" // Para manejo de errores

	"golang.org/x/crypto/bcrypt"
)

// Función de ayuda para hashear contraseña (puede estar en un paquete util)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // Costo por defecto es 10
	return string(bytes), err
}

// CreateUserInput DTO para los datos de entrada del caso de uso
type CreateUserInput struct {
	Username   string
	Password   string // Contraseña en texto plano
	MacAddress string // MAC Address (opcional, puede ser string vacío)
	Role       string // Rol (ej: 'user')
}

// CreateUserUseCase maneja la lógica de crear un usuario
type CreateUserUseCase struct {
	userRepo userDomain.UserRepository // Dependencia del repositorio
}

// NewCreateUserUseCase crea una instancia del caso de uso
func NewCreateUserUseCase(userRepo userDomain.UserRepository) *CreateUserUseCase {
	if userRepo == nil {
		log.Fatal("CRITICO: CreateUserUseCase recibió un userRepo nulo.")
	}
	return &CreateUserUseCase{userRepo: userRepo}
}

// Execute procesa la creación del usuario
func (uc *CreateUserUseCase) Execute(input CreateUserInput) error { // Devuelve solo error por simplicidad
	// 1. Validaciones básicas (ej: longitud mínima de contraseña)
	if len(input.Password) < 6 { // Ejemplo de validación simple
		log.Printf("WARN: [CreateUser] Intento de registro con contraseña corta para usuario: %s", input.Username)
		return fmt.Errorf("la contraseña debe tener al menos 6 caracteres")
	}
	if input.Username == "" {
		return fmt.Errorf("el nombre de usuario es requerido")
	}
	if input.Role == "" {
		input.Role = "user" // Rol por defecto si no se especifica
	}

	// 2. Hashear la contraseña
	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		log.Printf("ERROR: [CreateUser] Error al hashear contraseña para %s: %v", input.Username, err)
		return fmt.Errorf("error interno al procesar la contraseña") // Error genérico 500
	}

	// 3. Preparar la entidad User para guardar
	newUser := &userDomain.User{
		Username:     input.Username,
		PasswordHash: hashedPassword, // <- Se guarda el HASH
		Role:         input.Role,
		// Crear sql.NullString para mac_address (será NULL si input.MacAddress está vacío)
		MacAddress:   sql.NullString{String: input.MacAddress, Valid: input.MacAddress != ""},
	}

	// 4. Llamar al repositorio para crear el usuario
	err = uc.userRepo.Create(newUser)
	if err != nil {
		// Manejar errores específicos del repositorio
		if strings.HasSuffix(err.Error(), "_duplicado") { // Error específico que definimos en el repo
			log.Printf("WARN: [CreateUser] Conflicto al crear usuario: %v", err)
			// Devolver un error que el controlador pueda mapear a 409 Conflict
			return fmt.Errorf("conflicto: %w", err)
		}
		// Otros errores de la base de datos
		log.Printf("ERROR: [CreateUser] Error del repositorio al crear usuario %s: %v", input.Username, err)
		return fmt.Errorf("error interno al guardar el usuario") // Error genérico 500
	}

	log.Printf("INFO: [CreateUser] Usuario '%s' registrado exitosamente.", input.Username)
	return nil // Éxito
}