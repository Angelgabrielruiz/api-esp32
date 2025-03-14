package core

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql" // Importar el driver para MySQL
)

type Conn_MySQL struct {
	DB  *sql.DB
	Err string
}

// GetDBPool establece la conexión con la base de datos y devuelve una instancia de Conn_MySQL
func GetDBPool() *Conn_MySQL {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Printf("Advertencia: No se pudo cargar el archivo .env: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbSchema := os.Getenv("DB_SCHEMA")

	// Construcción del DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPass, dbHost, dbSchema)

	// Abrir conexión con MySQL
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return &Conn_MySQL{DB: nil, Err: fmt.Sprintf("error al abrir la base de datos: %v", err)}
	}

	// Verificar la conexión
	if err := db.Ping(); err != nil {
		db.Close()
		return &Conn_MySQL{DB: nil, Err: fmt.Sprintf("error al verificar la conexión: %v", err)}
	}

	// Configurar el pool de conexiones
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	fmt.Println("✅ Conexión exitosa a MySQL")
	return &Conn_MySQL{DB: db, Err: ""}
}

// Ejecutar consultas preparadas (INSERT, UPDATE, DELETE)
func (conn *Conn_MySQL) ExecutePreparedQuery(query string, values ...interface{}) (sql.Result, error) {
	stmt, err := conn.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error al preparar la consulta: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(values...)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta preparada: %w", err)
	}

	return result, nil
}

// Ejecutar consultas SELECT y devolver filas
func (conn *Conn_MySQL) FetchRows(query string, values ...interface{}) (*sql.Rows, error) {
	rows, err := conn.DB.Query(query, values...)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta SELECT: %w", err)
	}
	return rows, nil
}
