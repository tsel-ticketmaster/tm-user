package auth

type RegisterRequest struct {
	Name     string `json:"name" validate:"name"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"password"`
}
