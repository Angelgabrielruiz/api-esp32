package core

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time" // Importar para timeouts

	_ "github.com/go-sql-driver/mysql" // Driver MySQL
	//"github.com/joho/godotenv"
)

type Conn_MySQL struct {
	DB  *sql.DB
	Err string // Cambiado para que sea un string simple para el error inicial
}

// GetDBPool establece la conexión con la base de datos y devuelve una instancia de Conn_MySQL
func GetDBPool() *Conn_MySQL {
	// Cargar .env solo si no está ya cargado (opcional, main lo hace)
	// godotenv.Load()

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbSchema := os.Getenv("DB_SCHEMA")
	dbPort := os.Getenv("DB_PORT") // Opcional, por defecto 3306

	if dbPort == "" {
		dbPort = "3306" // Puerto por defecto de MySQL
	}

	if dbHost == "" || dbUser == "" || dbSchema == "" {
		log.Println("Advertencia: Variables de entorno DB_HOST, DB_USER, DB_SCHEMA no configuradas completamente.")
		// Podrías retornar error aquí si son obligatorias
	}


	// Construcción del DSN (Data Source Name)
	// Añadir parámetros útiles como parseTime y timeouts
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s&readTimeout=5s&writeTimeout=5s",
		dbUser, dbPass, dbHost, dbPort, dbSchema)

	// Abrir conexión con MySQL
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		// Retornar la instancia con el error almacenado
		return &Conn_MySQL{DB: nil, Err: fmt.Sprintf("error al abrir la conexión con la base de datos: %v", err)}
	}

	// Configurar el pool de conexiones
	db.SetMaxOpenConns(10) // Número máximo de conexiones abiertas
	db.SetMaxIdleConns(5)  // Número máximo de conexiones inactivas
	db.SetConnMaxLifetime(time.Minute * 5) // Tiempo máximo que una conexión puede ser reutilizada

	// Verificar la conexión inicial (Ping)
	// Es importante usar un contexto con timeout aquí
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := db.PingContext(ctx); err != nil {
    // Usar Ping simple por ahora si no manejas contextos en todos lados
	if err := db.Ping(); err != nil {
		db.Close() // Cerrar la conexión si el ping falla
		return &Conn_MySQL{DB: nil, Err: fmt.Sprintf("error al verificar la conexión (ping): %v", err)}
	}

	fmt.Println("✅ Conexión exitosa a MySQL establecida.")
	return &Conn_MySQL{DB: db, Err: ""} // Sin error inicial
}

// Close cierra la conexión a la base de datos. Debe ser llamado al final.
func (conn *Conn_MySQL) Close() error {
	if conn.DB != nil {
		fmt.Println("⚪ Cerrando conexión a MySQL...")
		return conn.DB.Close()
	}
	return nil
}


// ExecutePreparedQuery ejecuta consultas INSERT, UPDATE, DELETE.
// Es mejor usar transacciones para operaciones que modifican datos.
func (conn *Conn_MySQL) ExecutePreparedQuery(query string, values ...interface{}) (sql.Result, error) {
	if conn.Err != "" {
		return nil, fmt.Errorf("conexión inicial fallida: %s", conn.Err)
	}
	if conn.DB == nil {
        return nil, fmt.Errorf("la instancia de base de datos es nula")
    }

	// Considerar usar contexto con timeout
	stmt, err := conn.DB.Prepare(query) // Preparar fuera si se reutiliza mucho
	if err != nil {
		return nil, fmt.Errorf("error al preparar la consulta [%s]: %w", query, err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(values...)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta preparada [%s]: %w", query, err)
	}

	return result, nil
}

// FetchRows ejecuta consultas SELECT.
func (conn *Conn_MySQL) FetchRows(query string, values ...interface{}) (*sql.Rows, error) {
	if conn.Err != "" {
		return nil, fmt.Errorf("conexión inicial fallida: %s", conn.Err)
	}
    if conn.DB == nil {
        return nil, fmt.Errorf("la instancia de base de datos es nula")
    }
	// Considerar usar contexto con timeout
	rows, err := conn.DB.Query(query, values...)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta SELECT [%s]: %w", query, err)
	}
	// El вызывающий código es responsable de llamar a rows.Close()
	return rows, nil
}