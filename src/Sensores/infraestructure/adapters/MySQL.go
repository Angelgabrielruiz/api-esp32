// File: MySQL.go

package adapters

import (
	"API/src/core"
	"API/src/Sensores/domain/entities"
	"database/sql" // Necesario para sql.ErrNoRows
	"fmt"
	"log"
)

type MySQLRutas struct {
	conn *core.Conn_MySQL
}

func NewMySQLRutas(conn *core.Conn_MySQL) *MySQLRutas { // Ahora recibe la conexión
    if conn == nil || conn.DB == nil {
        log.Fatal("CRÍTICO: MySQLRutas recibió una conexión DB nula.")
    }
	return &MySQLRutas{conn: conn}
}


// Save AHORA incluye user_id
func (mysql *MySQLRutas) Save(userID int, temperatura string, movimiento string, distancia string, peso string, mac string) error {
    query := "INSERT INTO rutas (user_id, temperatura, movimiento, distancia, peso, mac) VALUES (?, ?, ?, ?, ?, ?)"
    result, err := mysql.conn.ExecutePreparedQuery(query, userID, temperatura, movimiento, distancia, peso, mac)
    // ... manejo de errores y logs como antes ...
			if err != nil {
	log.Printf("ERROR: [MySQLAdapter] Error al ejecutar INSERT: %v", err)
	return fmt.Errorf("error al guardar datos en MySQL: %w", err) // Envolver error
}

rowsAffected, _ := result.RowsAffected()
lastInsertId, _ := result.LastInsertId() // Útil si necesitas el ID

if rowsAffected == 1 {
	log.Printf("INFO: [MySQLAdapter] Datos insertados exitosamente para UserID %d. ID: %d, Filas afectadas: %d", userID, lastInsertId, rowsAffected)
} else {
	log.Printf("ADVERTENCIA: [MySQLAdapter] INSERT ejecutado pero filas afectadas no fue 1 (%d).", rowsAffected)
}
return nil
}

// GetAll (sin cambios si lo necesitas para admin)
func (mysql *MySQLRutas) GetAll() ([]entities.Datos, error) {
    // ... tu código existente ...
			query := "SELECT id, user_id, temperatura, movimiento, distancia, peso, mac FROM rutas ORDER BY id DESC" // Añadir user_id al SELECT
rows, err := mysql.conn.FetchRows(query)
if err != nil {
	log.Printf("ERROR: [MySQLAdapter] Error al ejecutar SELECT *: %v", err)
	return nil, fmt.Errorf("error al obtener datos de MySQL: %w", err)
}
defer rows.Close()

var datosList []entities.Datos

for rows.Next() {
	var dato entities.Datos
	// Añadir &dato.UserID al Scan en la posición correcta
	var userId sql.NullInt32 // Usar NullInt32 por si user_id es NULL en la BD
	if err := rows.Scan(&dato.ID, &userId, &dato.Temperatura, &dato.Movimiento, &dato.Distancia, &dato.Peso, &dato.Mac); err != nil {
		log.Printf("ERROR: [MySQLAdapter] Error al escanear fila (GetAll): %v", err)
		return nil, fmt.Errorf("error al procesar fila de datos MySQL: %w", err)
	}
	// Opcional: asignar UserID a tu struct si tiene el campo y no es NULL
	if userId.Valid {
		 // dato.UserID = userId.Int32 // Asumiendo que tienes un campo UserID en entities.Datos
	}
	datosList = append(datosList, dato)
}

if err := rows.Err(); err != nil {
	log.Printf("ERROR: [MySQLAdapter] Error durante la iteración de filas (GetAll): %v", err)
	return nil, fmt.Errorf("error final al leer datos de MySQL: %w", err)
}

log.Printf("INFO: [MySQLAdapter] Se recuperaron %d registros (GetAll).", len(datosList))
return datosList, nil
}

// NUEVO: GetByUserID
func (mysql *MySQLRutas) GetByUserID(userID int) ([]entities.Datos, error) {
    query := "SELECT id, user_id, temperatura, movimiento, distancia, peso, mac FROM rutas WHERE user_id = ? ORDER BY id DESC"
    rows, err := mysql.conn.FetchRows(query, userID)
    if err != nil {
        log.Printf("ERROR: [MySQLAdapter] Error al ejecutar SELECT por UserID %d: %v", userID, err)
        return nil, fmt.Errorf("error al obtener datos por usuario de MySQL: %w", err)
    }
    defer rows.Close()

    var datosList []entities.Datos
    for rows.Next() {
        var dato entities.Datos
					var dbUserId sql.NullInt32 // Necesario para Scan aunque ya filtramos
        if err := rows.Scan(&dato.ID, &dbUserId, &dato.Temperatura, &dato.Movimiento, &dato.Distancia, &dato.Peso, &dato.Mac); err != nil {
            log.Printf("ERROR: [MySQLAdapter] Error al escanear fila (GetByUserID: %d): %v", userID, err)
            return nil, fmt.Errorf("error al procesar fila de datos MySQL por usuario: %w", err)
        }
					// if dbUserId.Valid {
					// 	 dato.UserID = dbUserId.Int32 // Asumiendo que tienes UserID en tu struct
					// }
        datosList = append(datosList, dato)
    }

    if err := rows.Err(); err != nil {
        log.Printf("ERROR: [MySQLAdapter] Error durante la iteración de filas (GetByUserID: %d): %v", userID, err)
        return nil, fmt.Errorf("error final al leer datos por usuario de MySQL: %w", err)
    }

    log.Printf("INFO: [MySQLAdapter] Se recuperaron %d registros para UserID %d.", len(datosList), userID)
    return datosList, nil
}

// Update - Adaptar para recibir y potencialmente usar userID
func (mysql *MySQLRutas) Update(id int, userID int, temperatura string, movimiento string, distancia string, peso string, mac string) error {
    // Podrías añadir 'AND user_id = ?' al WHERE si solo el dueño puede actualizar
    query := "UPDATE rutas SET temperatura = ?, movimiento = ?, distancia = ?, peso = ?, mac = ? WHERE id = ?" // AND user_id = ?
    result, err := mysql.conn.ExecutePreparedQuery(query, temperatura, movimiento, distancia, peso, mac, id) // , userID
     // ... manejo de errores y logs como antes ...
			if err != nil {
	log.Printf("ERROR: [MySQLAdapter] Error al ejecutar UPDATE (ID: %d): %v", id, err)
	return fmt.Errorf("error al actualizar datos en MySQL (ID: %d): %w", id, err)
}

rowsAffected, _ := result.RowsAffected()
if rowsAffected == 1 {
	log.Printf("INFO: [MySQLAdapter] Datos actualizados exitosamente (ID: %d, UserID check: %d).", id, userID)
} else if rowsAffected == 0 {
	log.Printf("ADVERTENCIA: [MySQLAdapter] UPDATE ejecutado pero no se encontró el registro (ID: %d) o no pertenecía al usuario (UserID: %d) o los datos eran iguales.", id, userID)
	// Podrías devolver un error específico aquí si no se encontró
	// return fmt.Errorf("registro no encontrado o no perteneciente al usuario")
} else {
	log.Printf("ADVERTENCIA: [MySQLAdapter] UPDATE afectó un número inesperado de filas (%d) para ID %d.", rowsAffected, id)
}

return nil
}

// Delete - Adaptar para recibir y potencialmente usar userID
func (mysql *MySQLRutas) Delete(id int, userID int) error {
    // Podrías añadir 'AND user_id = ?' al WHERE si solo el dueño puede borrar
    query := "DELETE FROM rutas WHERE id = ?" // AND user_id = ?
    result, err := mysql.conn.ExecutePreparedQuery(query, id) // , userID
     // ... manejo de errores y logs como antes ...
			if err != nil {
	log.Printf("ERROR: [MySQLAdapter] Error al ejecutar DELETE (ID: %d): %v", id, err)
	return fmt.Errorf("error al eliminar datos en MySQL (ID: %d): %w", id, err)
}

rowsAffected, _ := result.RowsAffected()
if rowsAffected == 1 {
	log.Printf("INFO: [MySQLAdapter] Datos eliminados exitosamente (ID: %d, UserID check: %d).", id, userID)
} else if rowsAffected == 0 {
	log.Printf("ADVERTENCIA: [MySQLAdapter] DELETE ejecutado pero no se encontró el registro (ID: %d) o no pertenecía al usuario (UserID: %d).", id, userID)
	 // Podrías devolver un error específico aquí si no se encontró
	 // return fmt.Errorf("registro no encontrado o no perteneciente al usuario")
} else {
	log.Printf("ADVERTENCIA: [MySQLAdapter] DELETE afectó un número inesperado de filas (%d) para ID %d.", rowsAffected, id)
}

return nil
}