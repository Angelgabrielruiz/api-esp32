package adapters

import (
	"API/src/core"
	"API/src/Sensores/domain/entities" // Importar entidad
	"fmt"
	"log"
)

// MySQLRutas (renombrar a MySQLDatosRepository sería más claro) implementa domain.DatosRepository.
type MySQLRutas struct {
	conn *core.Conn_MySQL // Dependencia de la utilidad de conexión
}

// NewMySQLRutas crea una instancia del adaptador de repositorio MySQL.
func NewMySQLRutas() *MySQLRutas {
	conn := core.GetDBPool() // Obtiene la conexión/pool
	if conn.Err != "" {
		// Si la conexión inicial falló, loggear fatalmente o manejar el error.
		log.Fatalf("CRÍTICO: Error al configurar el pool de conexiones MySQL: %v", conn.Err)
	}
	if conn.DB == nil && conn.Err == "" {
        log.Fatal("CRÍTICO: GetDBPool devolvió una conexión nula sin error.")
    }
	return &MySQLRutas{conn: conn}
}

// Save implementa la inserción de datos.
func (mysql *MySQLRutas) Save(temperatura string, movimiento string, distancia string, peso string) error {
	query := "INSERT INTO rutas (temperatura, movimiento, distancia, peso) VALUES (?, ?, ?, ?)"
	result, err := mysql.conn.ExecutePreparedQuery(query, temperatura, movimiento, distancia, peso)
	if err != nil {
		log.Printf("ERROR: [MySQLAdapter] Error al ejecutar INSERT: %v", err)
		return fmt.Errorf("error al guardar datos en MySQL: %w", err) // Envolver error
	}

	rowsAffected, _ := result.RowsAffected()
	lastInsertId, _ := result.LastInsertId() // Útil si necesitas el ID

	if rowsAffected == 1 {
		log.Printf("INFO: [MySQLAdapter] Datos insertados exitosamente. ID: %d, Filas afectadas: %d", lastInsertId, rowsAffected)
	} else {
		log.Printf("ADVERTENCIA: [MySQLAdapter] INSERT ejecutado pero filas afectadas no fue 1 (%d).", rowsAffected)
	}
	return nil
}

// GetAll implementa la recuperación de todos los datos.
func (mysql *MySQLRutas) GetAll() ([]entities.Datos, error) {
	query := "SELECT id, temperatura, movimiento, distancia, peso FROM rutas ORDER BY id DESC" // Añadir un orden
	rows, err := mysql.conn.FetchRows(query)
	if err != nil {
		log.Printf("ERROR: [MySQLAdapter] Error al ejecutar SELECT *: %v", err)
		return nil, fmt.Errorf("error al obtener datos de MySQL: %w", err)
	}
	// MUY IMPORTANTE: Cerrar las filas al terminar.
	defer rows.Close()

	var datosList []entities.Datos // Usar el tipo de la entidad

	// Iterar sobre las filas
	for rows.Next() {
		var dato entities.Datos // Variable para escanear cada fila
		// Escanear en el orden de las columnas del SELECT
		if err := rows.Scan(&dato.ID, &dato.Temperatura, &dato.Movimiento, &dato.Distancia, &dato.Peso); err != nil {
			log.Printf("ERROR: [MySQLAdapter] Error al escanear fila: %v", err)
			// Continuar con las siguientes filas o retornar error aquí? Decisión de diseño.
			// Por ahora, retornamos error si falla el escaneo.
			return nil, fmt.Errorf("error al procesar fila de datos MySQL: %w", err)
		}
		// Añadir el dato escaneado a la lista
		datosList = append(datosList, dato)
	}

	// Verificar si hubo errores durante la iteración (ej: problema de conexión a mitad)
	if err := rows.Err(); err != nil {
		log.Printf("ERROR: [MySQLAdapter] Error durante la iteración de filas: %v", err)
		return nil, fmt.Errorf("error final al leer datos de MySQL: %w", err)
	}

	log.Printf("INFO: [MySQLAdapter] Se recuperaron %d registros.", len(datosList))
	return datosList, nil
}

// Update implementa la actualización de datos.
func (mysql *MySQLRutas) Update(id int, temperatura string, movimiento string, distancia string, peso string) error {
	query := "UPDATE rutas SET temperatura = ?, movimiento = ?, distancia = ?, peso = ? WHERE id = ?"
	result, err := mysql.conn.ExecutePreparedQuery(query, temperatura, movimiento, distancia, peso, id)
	if err != nil {
		log.Printf("ERROR: [MySQLAdapter] Error al ejecutar UPDATE (ID: %d): %v", id, err)
		return fmt.Errorf("error al actualizar datos en MySQL (ID: %d): %w", id, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("INFO: [MySQLAdapter] Datos actualizados exitosamente (ID: %d).", id)
	} else if rowsAffected == 0 {
		log.Printf("ADVERTENCIA: [MySQLAdapter] UPDATE ejecutado pero no se encontró el registro (ID: %d) o los datos eran iguales.", id)
		// Podrías considerar devolver un error específico tipo "NotFound" aquí.
        // return fmt.Errorf("registro con ID %d no encontrado para actualizar", id)
	} else {
		log.Printf("ADVERTENCIA: [MySQLAdapter] UPDATE afectó un número inesperado de filas (%d) para ID %d.", rowsAffected, id)
	}

	return nil
}

// Delete implementa la eliminación de datos.
func (mysql *MySQLRutas) Delete(id int) error {
	query := "DELETE FROM rutas WHERE id = ?"
	result, err := mysql.conn.ExecutePreparedQuery(query, id)
	if err != nil {
		log.Printf("ERROR: [MySQLAdapter] Error al ejecutar DELETE (ID: %d): %v", id, err)
		return fmt.Errorf("error al eliminar datos en MySQL (ID: %d): %w", id, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("INFO: [MySQLAdapter] Datos eliminados exitosamente (ID: %d).", id)
	} else if rowsAffected == 0 {
		log.Printf("ADVERTENCIA: [MySQLAdapter] DELETE ejecutado pero no se encontró el registro (ID: %d).", id)
		// Podrías considerar devolver un error específico tipo "NotFound" aquí.
        // return fmt.Errorf("registro con ID %d no encontrado para eliminar", id)
	} else {
		log.Printf("ADVERTENCIA: [MySQLAdapter] DELETE afectó un número inesperado de filas (%d) para ID %d.", rowsAffected, id)
	}

	return nil
}


// FindByID (Implementación opcional)
/*
func (mysql *MySQLRutas) FindByID(id int) (*entities.Datos, error) {
	query := "SELECT id, temperatura, movimiento, distancia, peso FROM rutas WHERE id = ?"
	rows, err := mysql.conn.FetchRows(query, id)
	if err != nil {
		log.Printf("ERROR: [MySQLAdapter] Error al buscar por ID %d: %v", id, err)
		return nil, fmt.Errorf("error al buscar datos por ID en MySQL: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var dato entities.Datos
		if err := rows.Scan(&dato.ID, &dato.Temperatura, &dato.Movimiento, &dato.Distancia, &dato.Peso); err != nil {
			log.Printf("ERROR: [MySQLAdapter] Error al escanear fila (ID: %d): %v", id, err)
			return nil, fmt.Errorf("error al procesar fila de datos MySQL (ID: %d): %w", id, err)
		}
		// Verificar si hubo más filas (no debería si ID es único)
		if rows.Next() {
             log.Printf("ADVERTENCIA: [MySQLAdapter] Se encontró más de un registro para ID %d", id)
             // Podrías retornar error o solo el primero
        }
        if err := rows.Err(); err != nil { // Chequear error después de escanear
             log.Printf("ERROR: [MySQLAdapter] Error después de escanear fila (ID: %d): %v", id, err)
             return nil, fmt.Errorf("error al leer datos de MySQL (ID: %d): %w", id, err)
         }
		return &dato, nil
	}

	// Verificar error después de rows.Next() si no encontró filas
    if err := rows.Err(); err != nil {
        log.Printf("ERROR: [MySQLAdapter] Error al buscar ID %d (después de Next): %v", id, err)
        return nil, fmt.Errorf("error al buscar datos por ID en MySQL: %w", err)
    }

	// No encontró filas y no hubo error
	log.Printf("INFO: [MySQLAdapter] No se encontró registro con ID %d.", id)
	// Devolver nil, nil o un error específico NotFound
    return nil, sql.ErrNoRows // Devolver error estándar para "no encontrado"
}
*/