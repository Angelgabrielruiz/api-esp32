package domain

type DatosRepository interface {
	Save(temperatura string, movimiento string) error
	GetAll() ([]map[string]interface{}, error)
	Update(id int, temperatura string, movimiento string) error
	Delete(id int) error
}
