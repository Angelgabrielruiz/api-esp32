package domain // O el paquete de dominio que uses para usuarios

// Opcional: Define una entidad User si no la tienes
type User struct {
    ID       int
    Username string
    // ... otros campos que necesites
}

type UserRepository interface {

    FindUserIDByMAC(macAddress string) (int, error)

}