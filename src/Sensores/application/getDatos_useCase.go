// File: getDatos_useCase.go

package application

import (
	"API/src/Sensores/domain"
	"API/src/Sensores/domain/entities"
	"log"
)

type GetDatos struct {
	db domain.DatosRepository
}

func NewGetDatos(db domain.DatosRepository) *GetDatos {
	if db == nil {
		log.Fatal("Error: GetDatos recibió dependencia db nula.")
	}
	return &GetDatos{db: db}
}

// Execute AHORA recibe el userID del usuario que hace la petición
func (gp *GetDatos) Execute(userID int) ([]entities.Datos, error) {
	// Llamar al nuevo método del repositorio
	datos, err := gp.db.GetByUserID(userID)
	if err != nil {
		log.Printf("ERROR: [GetDatos] Falló al obtener datos para UserID %d: %v", userID, err)
		return []entities.Datos{}, err // Devuelve slice vacío y error
	}
	log.Printf("INFO: [GetDatos] Se recuperaron %d registros para UserID %d.", len(datos), userID)
	return datos, nil
}

// Si necesitas una función para obtener TODOS (admin)
func (gp *GetDatos) ExecuteAll() ([]entities.Datos, error) {
    datos, err := gp.db.GetAll()
    if err != nil {
        log.Printf("ERROR: [GetDatos] Falló al obtener todos los datos (admin): %v", err)
        return []entities.Datos{}, err
    }
    log.Printf("INFO: [GetDatos] Se recuperaron %d registros en total (admin).", len(datos))
    return datos, nil
}