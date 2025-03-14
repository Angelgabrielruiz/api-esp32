// File: MySQL.go

package adapters

import (
	"API/src/core"
	"log"
)

type MySQLRutas struct {
    conn *core.Conn_MySQL
}

func NewMySQLRutas() *MySQLRutas {
    conn := core.GetDBPool()
    if conn.Err != "" {
        log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
    }
    return &MySQLRutas{conn: conn}
}

func (mysql *MySQLRutas) Save(temperatura string, movimiento string) error {
    query := "INSERT INTO rutas (temperatura, movimiento) VALUES (?, ?)"
    result, err := mysql.conn.ExecutePreparedQuery(query, temperatura, movimiento)
    if err != nil {
        return err
    }
    
    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 1 {
        log.Printf("[MySQL] - Ruta creada exitosamente: %d filas afectadas", rowsAffected)
    }
    return nil
}

func (mysql *MySQLRutas) GetAll() ([]map[string]interface{}, error) {
    query := "SELECT * FROM rutas"
    rows, err := mysql.conn.FetchRows(query)
    if err != nil {
    return nil, err
    }

    defer rows.Close()

    var rutas []map[string]interface{}
    for rows.Next() {
        var id int32
        var temperatura string
        var movimiento string
        if err := rows.Scan(&id, &temperatura, &movimiento); err != nil {
            return nil, err
        }
        
        ruta := map[string]interface{}{
            "id":          id,
            "temperatura": temperatura,
            "movimiento":  movimiento,
        }
        rutas = append(rutas, ruta)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return rutas, nil
}

func (mysql *MySQLRutas) Update(id int, temperatura string, movimiento string) error {
    query := "UPDATE rutas SET temperatura = ?, movimiento = ? WHERE id = ?"
    result, err := mysql.conn.ExecutePreparedQuery(query, temperatura, movimiento, id)
    if err != nil {
        return err
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 1 {
        log.Printf("[MySQL] - Ruta actualizada: %d", rowsAffected)
    }

    return nil
}

func (mysql *MySQLRutas) Delete(id int) error {
    query := "DELETE FROM rutas WHERE id = ?"
    result, err := mysql.conn.ExecutePreparedQuery(query, id)
    if err != nil {
        return err
    }
    
    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 1 {
        log.Printf("[MySQL] - Ruta eliminada: %d", rowsAffected)
    }
    
    return nil
}
