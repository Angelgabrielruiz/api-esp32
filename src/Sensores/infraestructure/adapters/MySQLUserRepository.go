package adapters // O el nombre de tu paquete de adaptadores

import (
    "API/src/core" // Tu paquete de conexión a BD
    // "API/src/Users/domain" // Si defines la entidad User
    "database/sql"
    "fmt"
    "log"
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

// FindUserIDByMAC implementa la interfaz UserRepository
func (repo *MySQLUserRepository) FindUserIDByMAC(macAddress string) (int, error) {
    var userID int
    query := "SELECT id FROM users WHERE mac_address = ? LIMIT 1"

    // Usaremos QueryRow que es ideal para esperar una sola fila.
    // Accedemos directamente a conn.DB que es el *sql.DB
    err := repo.conn.DB.QueryRow(query, macAddress).Scan(&userID)

    if err != nil {
        if err == sql.ErrNoRows {
            log.Printf("INFO: [UserRepo] No se encontró usuario para MAC: %s", macAddress)
            // Es importante devolver sql.ErrNoRows para que el caso de uso sepa que no existe
            return 0, sql.ErrNoRows
        }
        // Otro tipo de error (conexión, SQL inválido, etc.)
        log.Printf("ERROR: [UserRepo] Error al buscar usuario por MAC %s: %v", macAddress, err)
        return 0, fmt.Errorf("error al consultar usuario por MAC: %w", err)
    }

    log.Printf("INFO: [UserRepo] Usuario encontrado (ID: %d) para MAC: %s", userID, macAddress)
    return userID, nil
}

// Implementa otros métodos de UserRepository si los defines...