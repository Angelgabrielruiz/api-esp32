package domain

import "API/src/Sensores/domain/entities" // Importar la entidad

// DatosRepository define el "puerto" hacia la persistencia de datos.
// La capa de aplicación dependerá de esta interfaz, no de la implementación concreta.
type DatosRepository interface {
	// Save guarda nuevos datos de sensores. Podría devolver la entidad guardada con ID.
	Save(temperatura string, movimiento string, distancia string, peso string) error
	// GetAll recupera todos los registros de datos. Devolver []entities.Datos sería más type-safe.
	GetAll() ([]entities.Datos, error) // Cambiado para devolver slice de la entidad
	// Update actualiza un registro existente por ID.
	Update(id int, temperatura string, movimiento string, distancia string, peso string) error
	// Delete elimina un registro por ID.
	Delete(id int) error
	// FindByID (Opcional pero útil) recupera un registro por ID.
    // FindByID(id int) (*entities.Datos, error)
}