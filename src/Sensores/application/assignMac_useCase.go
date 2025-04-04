package application // O el paquete que uses para casos de uso

import (
	userDomain "API/src/Sensores/domain" // Ruta a tu paquete domain
	"database/sql"
	"fmt"
	"log"
	"net" // Para validar formato MAC (opcional)
)

// AssignMacInput DTO para la entrada
type AssignMacInput struct {
	TargetUserID int    // ID del usuario a modificar
	MacAddress   string // Nueva MAC address (puede ser vacía para desasignar)
}

// AssignMacToUserUseCase maneja la lógica de asignar MAC
type AssignMacToUserUseCase struct {
	userRepo userDomain.UserRepository
}

// NewAssignMacToUserUseCase crea la instancia
func NewAssignMacToUserUseCase(userRepo userDomain.UserRepository) *AssignMacToUserUseCase {
	if userRepo == nil {
		log.Fatal("CRITICO: AssignMacToUserUseCase recibió userRepo nulo.")
	}
	return &AssignMacToUserUseCase{userRepo: userRepo}
}

// Función de ayuda para validar formato MAC (opcional)
func isValidMacAddress(mac string) bool {
	if mac == "" { // Permitir desasignar con cadena vacía
		return true
	}
	_, err := net.ParseMAC(mac)
	return err == nil
}

// Execute realiza la asignación
func (uc *AssignMacToUserUseCase) Execute(input AssignMacInput) error {
	log.Printf("INFO: [AssignMacUC] Intentando asignar MAC '%s' a UserID %d", input.MacAddress, input.TargetUserID)

	// 1. Validar formato de MAC (Opcional pero recomendado)
	if !isValidMacAddress(input.MacAddress) {
		log.Printf("WARN: [AssignMacUC] Formato de MAC inválido: '%s'", input.MacAddress)
		return fmt.Errorf("formato_mac_invalido")
	}

	// 2. Preparar sql.NullString
	// Si la MAC está vacía, Valid será false, guardando NULL.
	macSQLString := sql.NullString{String: input.MacAddress, Valid: input.MacAddress != ""}

	// 3. Llamar al repositorio para actualizar
	err := uc.userRepo.UpdateMacAddress(input.TargetUserID, macSQLString)
	if err != nil {
		// Propagar errores específicos del repo (not found, duplicado) u otros
		log.Printf("ERROR: [AssignMacUC] Error del repositorio al actualizar MAC para UserID %d: %v", input.TargetUserID, err)
		return err // Devolver el error original del repo
	}

	log.Printf("INFO: [AssignMacUC] Operación de asignación de MAC completada para UserID %d.", input.TargetUserID)
	return nil // Éxito
}