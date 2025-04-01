package entities

// Datos representa la entidad de dominio para los datos de los sensores.
// Es buena práctica usar tipos más específicos si es posible (ej: float64, bool).
type Datos struct {
	ID          int32   `json:"id"` // Añadir ID y exportar campos para marshalling JSON
	Temperatura string  `json:"temperatura"` // Podría ser float64
	Movimiento  string  `json:"movimiento"`  // Podría ser bool o string ("si", "no")
	Distancia   string  `json:"distancia"`   // Podría ser float64
	Peso        string  `json:"peso"`        // Podría ser float64
}

// NewDatos es un constructor para crear una instancia de Datos.
// Nota: Generalmente, el ID se asigna por la base de datos o al recuperar datos.
// Este constructor es más útil si se crea desde datos de entrada sin ID.
func NewDatos(temperatura string, movimiento string, distancia string, peso string) *Datos {
	return &Datos{
        // ID se deja en 0 o se omite aquí
        Temperatura: temperatura,
        Movimiento:  movimiento,
        Distancia:   distancia,
        Peso:        peso,
    }
}

// Métodos Getters (opcional, pero bueno para encapsulación si los campos no son exportados)
// func (d *Datos) GetID() int32 { return d.id }
// func (d *Datos) GetTemperatura() string { return d.temperatura }
// ... etc ...