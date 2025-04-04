//File: manager.go

package websocket // El paquete es 'websocket' dentro de infraestructure

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket" 
)

// Manager maneja las conexiones WebSocket activas y el broadcasting.
type Manager struct {
	clients    map[*websocket.Conn]bool 
	broadcast  chan []byte             
	register   chan *websocket.Conn    
	unregister chan *websocket.Conn    
	mutex      sync.Mutex              
}

// NewManager crea e inicializa un nuevo WebSocket manager.
func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}


func (m *Manager) Run() {
	log.Println("INFO: WebSocket Manager iniciado y escuchando eventos...")
	for {
		select {
		case conn := <-m.register:
			// Registrar nuevo cliente
			m.mutex.Lock()
			m.clients[conn] = true
			m.mutex.Unlock()
			log.Printf("INFO: Cliente WebSocket conectado: %s. Clientes totales: %d", conn.RemoteAddr(), len(m.clients))

		case conn := <-m.unregister:
			// Desregistrar cliente
			m.mutex.Lock()
			// Verificar si el cliente aún existe antes de intentar borrar y cerrar
			if _, ok := m.clients[conn]; ok {
				delete(m.clients, conn)
				conn.Close() // Cerrar la conexión WebSocket
				log.Printf("INFO: Cliente WebSocket desconectado: %s. Clientes restantes: %d", conn.RemoteAddr(), len(m.clients))
			}
			m.mutex.Unlock()

		case message := <-m.broadcast:
			// Enviar mensaje a todos los clientes conectados
			m.mutex.Lock() // Bloquear mientras iteramos sobre los clientes
			if len(m.clients) > 0 {
				log.Printf("INFO: Transmitiendo mensaje a %d cliente(s) WebSocket...", len(m.clients))
			}
			for conn := range m.clients {
				// Enviar en una goroutine para no bloquear el broadcast si un cliente es lento
				go func(c *websocket.Conn, msg []byte) {
					
					if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
						log.Printf("ERROR: Error al escribir en WebSocket para %s: %v. Desregistrando cliente.", c.RemoteAddr(), err)
						m.unregister <- c
					}
				}(conn, message)
			}
			m.mutex.Unlock() // Desbloquear después de lanzar las goroutines de envío
		}
	}
}

// BroadcastMessage envía un mensaje al canal de broadcast para ser distribuido.
func (m *Manager) BroadcastMessage(message []byte) {
	// Simplemente envía el mensaje al canal. El bucle Run se encargará del resto.
	if len(message) > 0 {
		m.broadcast <- message
	} else {
		log.Println("ADVERTENCIA: Intento de transmitir mensaje WebSocket vacío.")
	}
}

// upgrader configura los parámetros para actualizar una conexión HTTP a WebSocket.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // Tamaño del buffer de lectura
	WriteBufferSize: 1024, // Tamaño del buffer de escritura
	
	CheckOrigin: func(r *http.Request) bool {
		// Ejemplo permisivo para desarrollo:
		log.Printf("DEBUG: Verificando origen de WebSocket: %s", r.Header.Get("Origin"))
		return true 

	},
}


func (m *Manager) HandleConnections(w http.ResponseWriter, r *http.Request) {
	
	conn, err := upgrader.Upgrade(w, r, nil) // w, r, y cabeceras adicionales (nil aquí)
	if err != nil {
		
		log.Printf("ERROR: Falló la actualización a WebSocket: %v", err)
		
		return
	}

	// Registrar la nueva conexión exitosa.
	m.register <- conn

	
	go m.readLoop(conn)
}

// readLoop lee mensajes del cliente WebSocket. Necesario para detectar cierres.
func (m *Manager) readLoop(conn *websocket.Conn) {
	
	defer func() {
		m.unregister <- conn
	}()

	

	for {
		
		_, _, err := conn.ReadMessage()
		if err != nil {
			// Verificar si el error es un cierre esperado de la conexión.
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				// Error inesperado
				log.Printf("ERROR: Error inesperado de lectura WebSocket para %s: %v", conn.RemoteAddr(), err)
			} else {
				
				log.Printf("INFO: Conexión WebSocket cerrada por el cliente %s: %v", conn.RemoteAddr(), err)
			}
			
			break
		}
		
	}
}