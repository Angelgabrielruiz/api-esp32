//Files/datos.go

package entities

type Datos struct {
	ID          int32   `json:"id"` 
	Temperatura string  `json:"temperatura"` // Podría ser float64
	Movimiento  string  `json:"movimiento"`  // Podría ser bool o string ("si", "no")
	Distancia   string  `json:"distancia"`   // Podría ser float64
	Peso        string  `json:"peso"`
	Mac			string	 `json:"peso"`        // Podría ser float64
}

func NewDatos(temperatura string, movimiento string, distancia string, peso string, mac string) *Datos {
	return &Datos{
        Temperatura: temperatura,
        Movimiento:  movimiento,
        Distancia:   distancia,
        Peso:        peso,
        Mac:		 mac,
    }
}
