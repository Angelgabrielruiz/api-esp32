package application

import "API/src/Sensores/domain"

type CreateDatos struct {
	db domain.DatosRepository
}

func NewCreateDatos(db domain.DatosRepository) *CreateDatos {
	return &CreateDatos{db: db}
}

func (cr *CreateDatos) Execute(temperatura string, movimiento string, distancia string, peso string) error {
	return cr.db.Save(temperatura, movimiento, distancia, peso)
}