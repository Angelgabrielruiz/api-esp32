// File: src/Sensores/domain/datosRepository.go

package domain

import "API/src/Sensores/domain/entities"

type DatosRepository interface {
    // Save ahora requiere el user_id asociado
    Save(userID int, temperatura string, movimiento string, distancia string, peso string, mac string) error

    // Para obtener TODOS los datos (quizás para un admin)
    GetAll() ([]entities.Datos, error)

    // NUEVO: Para obtener datos solo de un usuario específico
    GetByUserID(userID int) ([]entities.Datos, error)

    // Update y Delete probablemente también deberían verificar el user_id si la lógica lo requiere
    Update(id int, userID int, temperatura string, movimiento string, distancia string, peso string, mac string) error // userID añadido para posible validación
    Delete(id int, userID int) error // userID añadido para posible validación
}