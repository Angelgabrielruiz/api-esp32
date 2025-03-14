package entities

type Datos struct {
	id int32
	temperatura string
	movimiento  string
}

func NewDatos(temperatura string, movimiento string) *Datos {
	return &Datos{temperatura: temperatura, movimiento: movimiento}
}
