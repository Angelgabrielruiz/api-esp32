package entities

type Datos struct {
	id int32
	temperatura string
	movimiento  string
	distancia   string
	peso 	    string
}

func NewDatos(temperatura string, movimiento string, distancia string, peso string) *Datos {
	return &Datos{temperatura: temperatura, movimiento: movimiento, distancia: distancia, peso: peso}
}
