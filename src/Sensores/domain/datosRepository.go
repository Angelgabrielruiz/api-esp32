package domain

import "API/src/Sensores/domain/entities" // Importar la entidad

type DatosRepository interface {
	Save(temperatura string, movimiento string, distancia string, peso string, mac string) error
	GetAll() ([]entities.Datos, error) // Cambiado para devolver slice de la entidad
	Update(id int, temperatura string, movimiento string, distancia string, peso string, mac string) error
	Delete(id int) error

}