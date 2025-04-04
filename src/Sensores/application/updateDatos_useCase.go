//File: updateDatos_useCase.go:

package application

import (
	"API/src/Sensores/domain"
	// "API/src/Sensores/domain/entities" // Importar si necesitas notificar datos actualizados
	"log"
)

type UpdateDatos struct {
	db domain.DatosRepository
	// notifier DatosNotifier // Añadir si notificas actualizaciones
}

func NewUpdateDatos(db domain.DatosRepository /*, notifier DatosNotifier */) *UpdateDatos {
	if db == nil {
		log.Fatal("Error: UpdateDatos recibió dependencia db nula.")
	}
	return &UpdateDatos{
		db: db,
		// notifier: notifier,
	}
}

func (up *UpdateDatos) Execute(id int, temperatura string, movimiento string, distancia string, peso string, mac string) error {
	err := up.db.Update(id, 0, temperatura, movimiento, distancia, peso, mac) // Replace '0' with the appropriate int value
	if err != nil {
		log.Printf("ERROR: [UpdateDatos] Falló al actualizar datos (ID: %d): %v", id, err)
		return err
	}

	// if up.notifier != nil {
	// 	updatedData := entities.Datos{ ID: int32(id), Temperatura: temperatura, /* ... otros campos ... */ }
	// 	if errNotify := up.notifier.NotifyDataUpdated(updatedData); errNotify != nil {
	// 		log.Printf("ADVERTENCIA: [UpdateDatos] Falló la notificación de actualización (ID: %d): %v", id, errNotify)
	// 	}
	// }

	log.Printf("INFO: [UpdateDatos] Datos con ID %d actualizados exitosamente.", id)
	return nil
}