package application

import (
	"API/src/Sensores/domain"
	"log"
)

type DeleteDatos struct {
	db domain.DatosRepository
	// notifier DatosNotifier // Podrías añadir notifier si quieres notificar eliminaciones
}

func NewDeleteDatos(db domain.DatosRepository /*, notifier DatosNotifier */) *DeleteDatos {
	if db == nil {
		log.Fatal("Error: DeleteDatos recibió dependencia db nula.")
	}
	return &DeleteDatos{
		db: db,
		// notifier: notifier,
	}
}

// Execute ahora devuelve error para indicar si la eliminación falló.
func (dp *DeleteDatos) Execute(id int) error {
	err := dp.db.Delete(id)
	if err != nil {
		log.Printf("ERROR: [DeleteDatos] Falló al eliminar datos con ID %d: %v", id, err)
		return err // Devuelve el error
	}

	// if dp.notifier != nil {
	// 	if errNotify := dp.notifier.NotifyDataDeleted(id); errNotify != nil {
	// 		log.Printf("ADVERTENCIA: [DeleteDatos] Falló la notificación de eliminación (ID: %d): %v", id, errNotify)
	// 	}
	// }

	log.Printf("INFO: [DeleteDatos] Datos con ID %d eliminados exitosamente.", id)
	return nil // Éxito
}