package application

import "API/src/Sensores/domain"

type UpdateDatos struct {
	db domain.DatosRepository
}

func NewUpdateDatos(db domain.DatosRepository) *UpdateDatos {
	return &UpdateDatos{db: db}
}

func (up *UpdateDatos) Execute(id int, temperatura string, movimiento string, distancia string, peso string) error {
	return up.db.Update(id, temperatura, movimiento, distancia, peso)
}