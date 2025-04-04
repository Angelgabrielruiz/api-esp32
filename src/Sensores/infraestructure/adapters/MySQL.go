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
func (mysql *MySQLRutas) Save(temperatura string, movimiento string, distancia string, peso string, mac string) error {
	query := "INSERT INTO rutas (temperatura, movimiento, distancia, peso, mac) VALUES (?, ?, ?, ?, ?)"
	result, err := mysql.conn.ExecutePreparedQuery(query, temperatura, movimiento, distancia, peso, mac)
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
	query := "SELECT id, temperatura, movimiento, distancia, peso, mac FROM rutas ORDER BY id DESC" // Añadir un orden
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
		if err := rows.Scan(&dato.ID, &dato.Temperatura, &dato.Movimiento, &dato.Distancia, &dato.Peso, &dato.Mac); err != nil {
			log.Printf("ERROR: [MySQLAdapter] Error al escanear fila: %v", err)
			
			return nil, fmt.Errorf("error al procesar fila de datos MySQL: %w", err)
		}
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
func (mysql *MySQLRutas) Update(id int, temperatura string, movimiento string, distancia string, peso string, mac string) error {
	query := "UPDATE rutas SET temperatura = ?, movimiento = ?, distancia = ?, peso = ?, mac = ? WHERE id = ?"
	result, err := mysql.conn.ExecutePreparedQuery(query, temperatura, movimiento, distancia, peso, mac, id)
	if err != nil {
		log.Printf("ERROR: [MySQLAdapter] Error al ejecutar UPDATE (ID: %d): %v", id, err)
		return fmt.Errorf("error al actualizar datos en MySQL (ID: %d): %w", id, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("INFO: [MySQLAdapter] Datos actualizados exitosamente (ID: %d).", id)
	} else if rowsAffected == 0 {
		log.Printf("ADVERTENCIA: [MySQLAdapter] UPDATE ejecutado pero no se encontró el registro (ID: %d) o los datos eran iguales.", id)
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
	} else {
		log.Printf("ADVERTENCIA: [MySQLAdapter] DELETE afectó un número inesperado de filas (%d) para ID %d.", rowsAffected, id)
	}

	return nil
}
