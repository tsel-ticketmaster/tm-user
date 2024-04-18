package customer

type ChangeEmailEvent struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	ExistingEmail   string `json:"existing_email"`
	NewEmail        string `json:"new_email"`
	VerficationLink string `json:"verification_link"`
}
