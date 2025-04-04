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


func (n *WebSocketNotifier) NotifyNewData(data entities.Datos) error {

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("ERROR: [WebSocketNotifier] Error al codificar datos a JSON: %v. Data: %+v", err, data)
		// Retornar el error de codificación, ya que no podemos enviar nada útil.
		return fmt.Errorf("error al codificar datos para websocket: %w", err)
	}

	// Usar el manager para transmitir el mensaje JSON a todos los clientes conectados
	log.Printf("INFO: [WebSocketNotifier] Transmitiendo datos vía WebSocket: %s", string(jsonData))
	n.wsManager.BroadcastMessage(jsonData) // El manager se encarga del envío

	return nil
}

