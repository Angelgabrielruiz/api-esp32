package domain

import "API/src/Sensores/domain/entities" // Importar entidad si se usa en la interfaz

// DatosNotifier define el "puerto" secundario (driven port) para notificar
// sobre eventos relacionados con los datos.
type DatosNotifier interface {
	// NotifyNewData se llama cuando se crean nuevos datos.
	// Usar la entidad o un DTO específico es mejor que map[string]interface{}.
	NotifyNewData(data entities.Datos) error
}

// Podrías añadir más métodos si necesitas notificar otras acciones:
// type DatosNotifier interface {
//     NotifyNewData(data entities.Datos) error
//     NotifyDataUpdated(data entities.Datos) error
//     NotifyDataDeleted(id int) error
// }