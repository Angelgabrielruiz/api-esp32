package application

import (
	//"API/src/Sensores/domain"
	"API/src/Sensores/domain"
	"API/src/Sensores/domain/entities"
	//"API/src/Sensores/domain/entities" // Importar entidad
	"log"
)

type CreateDatos struct {
	db       domain.DatosRepository // Puerto hacia la persistencia
	notifier domain.DatosNotifier        // Puerto hacia la notificación
}


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
func (cr *CreateDatos) Execute(temperatura string, movimiento string, distancia string, peso string, mac string) error {
	// 1. Guardar en la base de datos usando el repositorio
	err := cr.db.Save(temperatura, movimiento, distancia, peso, mac)
	if err != nil {
		log.Printf("ERROR: [CreateDatos] Falló al guardar datos: %v", err)
		return err // Retornar el error de guardado
	}

	newData := entities.Datos{
		Temperatura: temperatura,
		Movimiento:  movimiento,
		Distancia:   distancia,
		Peso:        peso,
		Mac:		 mac,
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