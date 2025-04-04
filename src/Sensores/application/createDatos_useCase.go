// File: src/Sensores/application/createDatos_useCase.go

package application

import (
	sensorDomain "API/src/Sensores/domain" // Alias para claridad
	"API/src/Sensores/domain/entities"
	userDomain "API/src/Sensores/domain" // Importar dominio de usuarios
	"database/sql"                    // Para sql.ErrNoRows
	"fmt"
	"log"
)

type CreateDatos struct {
	datosRepo sensorDomain.DatosRepository // Puerto hacia persistencia de sensores
	userRepo  userDomain.UserRepository   // NUEVO: Puerto hacia persistencia de usuarios
	notifier  sensorDomain.DatosNotifier  // Puerto hacia la notificación
}

// Ahora recibe UserRepository también
func NewCreateDatos(datosRepo sensorDomain.DatosRepository, userRepo userDomain.UserRepository, notifier sensorDomain.DatosNotifier) *CreateDatos {
	if datosRepo == nil || notifier == nil || userRepo == nil {
		log.Fatal("Error: CreateDatos recibió dependencias nulas (datosRepo, userRepo o notifier).")
	}
	return &CreateDatos{
		datosRepo: datosRepo,
		userRepo:  userRepo,
		notifier:  notifier,
	}
}

// Execute ya NO necesita id, recibe los datos tal cual llegan
// Cambiar los parámetros para que coincidan con tu struct SensorData o el JSON
func (cr *CreateDatos) Execute(temperatura string, movimiento string, distancia string, peso string, mac string) error {
	// 1. Validar MAC (opcional pero recomendado)
	if mac == "" {
		log.Println("ERROR: [CreateDatos] Se recibió un mensaje sin dirección MAC.")
		// Puedes decidir devolver un error específico aquí si la MAC es obligatoria
		return fmt.Errorf("dirección MAC es requerida")
	}

	// 2. Buscar el UserID asociado a la MAC
	userID, err := cr.userRepo.FindUserIDByMAC(mac)
	if err != nil {
		if err == sql.ErrNoRows {
			// MAC no asignada a ningún usuario
			log.Printf("ADVERTENCIA: [CreateDatos] MAC '%s' recibida pero no está asignada a ningún usuario. Descartando datos.", mac)
			// DECISIÓN IMPORTANTE: ¿Qué hacer aquí?
			// Opción 1: Devolver un error específico para que el consumidor haga ACK (no reintentar)
			return fmt.Errorf("mac_no_asignada: %s", mac) // Error específico
			// Opción 2: Simplemente retornar nil (ignorar silenciosamente)
			// return nil
		}
		// Otro error al buscar el usuario
		log.Printf("ERROR: [CreateDatos] Falló la búsqueda de usuario por MAC '%s': %v", mac, err)
		return fmt.Errorf("error interno al buscar usuario: %w", err) // Error genérico
	}

	// 3. Guardar en la base de datos usando el repositorio, AHORA con UserID
	err = cr.datosRepo.Save(userID, temperatura, movimiento, distancia, peso, mac)
	if err != nil {
		log.Printf("ERROR: [CreateDatos] Falló al guardar datos para UserID %d (MAC %s): %v", userID, mac, err)
		return err // Retornar el error de guardado
	}

	log.Printf("INFO: [CreateDatos] Datos guardados exitosamente para UserID %d (MAC %s).", userID, mac)

	// 4. Notificar (si usas WebSockets dirigidos, necesitarás el userID)
	newData := entities.Datos{
		// ID se genera en la BD, no lo tenemos aquí a menos que Save lo devuelva
		Temperatura: temperatura,
		Movimiento:  movimiento,
		Distancia:   distancia,
		Peso:        peso,
		Mac:         mac,
		// UserID:     int32(userID), // Añade UserID a tu entidad si lo necesitas en el frontend
	}


	// ASUMIENDO que NotifyNewData ahora necesita el userID para dirigir el mensaje
	// Cambia la firma de NotifyNewData en la interfaz y la implementación
	// if errNotify := cr.notifier.NotifyNewData(userID, newData); errNotify != nil {
	// Sin userID por ahora, asumiendo broadcast o notificación genérica
   if errNotify := cr.notifier.NotifyNewData(newData); errNotify != nil {
		log.Printf("ADVERTENCIA: [CreateDatos] Falló la notificación de nuevos datos para UserID %d (pero fueron guardados): %v", userID, errNotify)
	} else {
		log.Printf("INFO: [CreateDatos] Notificación de nuevos datos iniciada exitosamente para UserID %d.", userID)
	}

	return nil
}