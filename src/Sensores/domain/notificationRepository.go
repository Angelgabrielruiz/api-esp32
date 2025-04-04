//File: notificationRepository.go

package domain

import "API/src/Sensores/domain/entities" // Importar entidad si se usa en la interfaz

type DatosNotifier interface {

	NotifyNewData(data entities.Datos) error
}

