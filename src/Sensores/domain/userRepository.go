package domain

import "database/sql"

// User representa la entidad de usuario
type User struct {
	ID            int
	Username      string
	PasswordHash  string
	MacAddress    sql.NullString
	Role          string
}

// UserRepository define las operaciones de persistencia para usuarios
type UserRepository interface {
	FindUserIDByMAC(macAddress string) (int, error)
	FindByUsername(username string) (*User, error)
	Create(user *User) error // Asumiendo que ya añadiste este

	// --- AÑADIR ESTE MÉTODO ---
	UpdateMacAddress(userID int, macAddress sql.NullString) error
	// --------------------------
}