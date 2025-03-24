package application

import "API/src/Sensores/domain"


type DeleteDatos struct {
	db domain.DatosRepository
}

func NewDeleteDatos(db domain.DatosRepository) *DeleteDatos {
	return &DeleteDatos{db: db}
}

func (dp *DeleteDatos) Execute(id int)  {
	dp.db.Delete(id)
}
