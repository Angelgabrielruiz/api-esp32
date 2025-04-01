package application

import (
	"API/src/Sensores/domain"
	"API/src/Sensores/domain/entities" // Importar entidad
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

// Execute ahora devuelve []entities.Datos y error.
func (gp *GetDatos) Execute() ([]entities.Datos, error) {
	datos, err := gp.db.GetAll()
	if err != nil {
		log.Printf("ERROR: [GetDatos] Falló al obtener todos los datos: %v", err)
		// Devuelve un slice vacío y el error
		return []entities.Datos{}, err
	}
	log.Printf("INFO: [GetDatos] Se recuperaron %d registros.", len(datos))
	return datos, nil
}