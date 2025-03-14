package application

import "API/src/Sensores/domain"

type GetDatos struct {
	db domain.DatosRepository
}

func NewGetDatos(db domain.DatosRepository) *GetDatos {
	return &GetDatos{db: db}
}

func (gp *GetDatos) Execute() ([]map[string]interface{}, error) {
	return gp.db.GetAll()
}
