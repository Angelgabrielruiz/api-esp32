package adapters

import (
	"API/src/core"
	"API/src/Sensores/domain" // Asegúrate que la ruta sea correcta
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type MySQLUserRepository struct {
	conn *core.Conn_MySQL
}

func NewMySQLUserRepository(conn *core.Conn_MySQL) *MySQLUserRepository {
	if conn == nil || conn.DB == nil {
		log.Fatal("CRÍTICO: MySQLUserRepository recibió una conexión DB nula.")
	}
	return &MySQLUserRepository{conn: conn}
}

// --- IMPLEMENTACIÓN MÉTODO Create ---
func (repo *MySQLUserRepository) Create(user *domain.User) error {
	query := "INSERT INTO users (username, password_hash, role, mac_address) VALUES (?, ?, ?, ?)"
	result, err := repo.conn.ExecutePreparedQuery(query,
		user.Username,
		user.PasswordHash,
		user.Role,
		user.MacAddress, // sql.NullString
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			fieldName := "desconocido"
			if strings.Contains(mysqlErr.Message, "users.username") {
				fieldName = "username"
			} else if strings.Contains(mysqlErr.Message, "mac_address") {
				fieldName = "mac_address" // Ajusta a tu constraint unique
			}
			log.Printf("WARN: [UserRepo] Intento de crear usuario con %s duplicado: %s", fieldName, user.Username)
			return fmt.Errorf("%s_duplicado", fieldName)
		}
		log.Printf("ERROR: [UserRepo] Error al ejecutar INSERT para usuario %s: %v", user.Username, err)
		return fmt.Errorf("error al guardar usuario en la base de datos: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected != 1 {
		log.Printf("ADVERTENCIA: [UserRepo] INSERT ejecutado para usuario %s pero filas afectadas fue %d (se esperaba 1).", user.Username, rowsAffected)
	} else {
		lastInsertId, _ := result.LastInsertId()
		log.Printf("INFO: [UserRepo] Usuario '%s' creado exitosamente con ID: %d", user.Username, lastInsertId)
	}
	return nil
}

// --- IMPLEMENTACIÓN MÉTODO FindUserIDByMAC ---
func (repo *MySQLUserRepository) FindUserIDByMAC(macAddress string) (int, error) {
	var userID int
	query := "SELECT id FROM users WHERE mac_address = ? LIMIT 1"
	err := repo.conn.DB.QueryRow(query, macAddress).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("INFO: [UserRepo] No se encontró usuario para MAC: %s", macAddress)
			return 0, sql.ErrNoRows
		}
		log.Printf("ERROR: [UserRepo] Error al buscar usuario por MAC %s: %v", macAddress, err)
		return 0, fmt.Errorf("error al consultar usuario por MAC: %w", err)
	}
	log.Printf("INFO: [UserRepo] Usuario encontrado (ID: %d) para MAC: %s", userID, macAddress)
	return userID, nil
}

// --- IMPLEMENTACIÓN MÉTODO FindByUsername ---
func (repo *MySQLUserRepository) FindByUsername(username string) (*domain.User, error) {
	user := &domain.User{}
	query := "SELECT id, username, password_hash, mac_address, role FROM users WHERE username = ? LIMIT 1"
	var macAddress sql.NullString
	err := repo.conn.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&macAddress,
		&user.Role,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("INFO: [UserRepo] No se encontró usuario: %s", username)
			return nil, sql.ErrNoRows
		}
		log.Printf("ERROR: [UserRepo] Error al buscar usuario %s: %v", username, err)
		return nil, fmt.Errorf("error al consultar usuario: %w", err)
	}
	user.MacAddress = macAddress
	log.Printf("INFO: [UserRepo] Usuario encontrado: %s (ID: %d)", username, user.ID)
	return user, nil
}

// --- IMPLEMENTACIÓN MÉTODO UpdateMacAddress ---
func (repo *MySQLUserRepository) UpdateMacAddress(userID int, macAddress sql.NullString) error {
	query := "UPDATE users SET mac_address = ? WHERE id = ?"
	result, err := repo.conn.ExecutePreparedQuery(query, macAddress, userID)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
            if strings.Contains(mysqlErr.Message, "mac_address") { // Ajusta al nombre de tu constraint
                log.Printf("WARN: [UserRepo] Intento de asignar MAC duplicada (%v) a UserID %d", macAddress, userID)
                return fmt.Errorf("mac_address_duplicado")
            }
		}
		log.Printf("ERROR: [UserRepo] Error al ejecutar UPDATE mac_address para UserID %d: %v", userID, err)
		return fmt.Errorf("error al actualizar MAC del usuario: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("WARN: [UserRepo] UPDATE mac_address ejecutado pero no se encontró UserID %d.", userID)
		return sql.ErrNoRows
	} else if rowsAffected == 1 {
        if macAddress.Valid {
			log.Printf("INFO: [UserRepo] MAC Address '%s' asignada/actualizada exitosamente para UserID %d.", macAddress.String, userID)
		} else {
            log.Printf("INFO: [UserRepo] MAC Address desasignada (NULL) exitosamente para UserID %d.", userID)
        }
	} else {
		log.Printf("ADVERTENCIA: [UserRepo] UPDATE mac_address afectó %d filas para UserID %d (se esperaba 0 o 1).", rowsAffected, userID)
	}
	return nil
}