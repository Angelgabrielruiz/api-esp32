package application

import (
	//"API/src/Sensores/domain"
	"API/src/Sensores/domain"
	"API/src/Sensores/domain/entities"
	//"API/src/Sensores/domain/entities" // Importar entidad
	"log"
)

// CreateDatos es el caso de uso para crear nuevos datos de sensores.
// Orquesta la interacción entre el repositorio y el notificador.
type CreateDatos struct {
	db       domain.DatosRepository // Puerto hacia la persistencia
	notifier domain.DatosNotifier        // Puerto hacia la notificación
}

// NewCreateDatos constructor inyecta las dependencias (interfaces).
func NewCreateDatos(db domain.DatosRepository, notifier domain.DatosNotifier) *CreateDatos {
	if db == nil || notifier == nil {
		// Es crucial validar las dependencias inyectadas
		log.Fatal("Error: CreateDatos recibió dependencias nulas (db o notifier).")
	}
	return &CreateDatos{
		db:       db,
		notifier: notifier,
	}
}

// Execute contiene la lógica de negocio principal para la creación.
func (cr *CreateDatos) Execute(temperatura string, movimiento string, distancia string, peso string) error {
	// 1. Guardar en la base de datos usando el repositorio
	err := cr.db.Save(temperatura, movimiento, distancia, peso)
	if err != nil {
		log.Printf("ERROR: [CreateDatos] Falló al guardar datos: %v", err)
		return err // Retornar el error de guardado
	}

	// 2. Si el guardado fue exitoso, notificar usando el puerto de notificación.
	//    Crear la entidad/DTO con los datos que se guardaron para notificar.
	//    Nota: Aquí no tenemos el ID generado por la DB a menos que Save lo devuelva.
	newData := entities.Datos{
		Temperatura: temperatura,
		Movimiento:  movimiento,
		Distancia:   distancia,
		Peso:        peso,
	}

	// Llamar al método de la interfaz del notificador
	if errNotify := cr.notifier.NotifyNewData(newData); errNotify != nil {
		// No hacer fallar la operación principal por un fallo de notificación.
		// Solo registrar la advertencia.
		log.Printf("ADVERTENCIA: [CreateDatos] Falló la notificación de nuevos datos (pero fueron guardados): %v", errNotify)
	} else {
		log.Printf("INFO: [CreateDatos] Notificación de nuevos datos iniciada exitosamente.")
	}


	// La operación principal (guardado) fue exitosa.
	return nil
}