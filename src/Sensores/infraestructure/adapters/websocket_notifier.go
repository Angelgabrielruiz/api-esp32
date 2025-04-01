package adapters

import (
	"API/src/Sensores/domain/entities"           // Importar la entidad
	infraWS "API/src/Sensores/infraestructure/websocket" // Importar el manager desde su ubicación
	"encoding/json"
	"fmt"
	"log"
)

// WebSocketNotifier es el "adaptador" que cumple con la interfaz application.DatosNotifier
// utilizando el WebSocket Manager como mecanismo de envío.
type WebSocketNotifier struct {
	wsManager *infraWS.Manager // Dependencia del manager concreto
}

// NewWebSocketNotifier crea una nueva instancia del adaptador notificador.
func NewWebSocketNotifier(wsManager *infraWS.Manager) *WebSocketNotifier {
	if wsManager == nil {
		log.Fatal("CRÍTICO: Se intentó crear WebSocketNotifier con un wsManager nulo.")
	}
	log.Println("INFO: Adaptador WebSocketNotifier creado.")
	return &WebSocketNotifier{wsManager: wsManager}
}

// NotifyNewData implementa el método de la interfaz application.DatosNotifier.
// Toma la entidad de datos, la convierte a JSON y la envía a través del WebSocket Manager.
func (n *WebSocketNotifier) NotifyNewData(data entities.Datos) error {
	// Convertir la entidad de datos a JSON
	// Los campos de entities.Datos deben estar exportados (mayúscula inicial) y tener `json:"tag"`
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("ERROR: [WebSocketNotifier] Error al codificar datos a JSON: %v. Data: %+v", err, data)
		// Retornar el error de codificación, ya que no podemos enviar nada útil.
		return fmt.Errorf("error al codificar datos para websocket: %w", err)
	}

	// Usar el manager para transmitir el mensaje JSON a todos los clientes conectados
	log.Printf("INFO: [WebSocketNotifier] Transmitiendo datos vía WebSocket: %s", string(jsonData))
	n.wsManager.BroadcastMessage(jsonData) // El manager se encarga del envío

	// Asumimos que BroadcastMessage es asíncrono o maneja errores internamente.
	// Si necesitáramos saber si el broadcast falló para algún cliente, el manager
	// tendría que exponer esa información, lo cual complica el diseño.
	// Para este caso, retornamos nil indicando que la notificación fue *iniciada*.
	return nil
}

// Podrías implementar otros métodos de notificación aquí si la interfaz los tuviera
// func (n *WebSocketNotifier) NotifyDataUpdated(data entities.Datos) error { ... }
// func (n *WebSocketNotifier) NotifyDataDeleted(id int) error { ... }