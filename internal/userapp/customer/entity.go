package customer

import "github.com/google/uuid"

type Customer struct {
	ID           uuid.UUID
	Name         string
	Email        string
	Password     string
	PasswordSalt string
}
