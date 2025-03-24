package domain

type DatosRepository interface {
	Save(temperatura string, movimiento string, distancia string, peso string) error
	GetAll() ([]map[string]interface{}, error)
	Update(id int, temperatura string, movimiento string, distancia string, peso string) error
	Delete(id int) error
}
