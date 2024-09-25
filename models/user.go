package models

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

type ConfirmUser struct {
    Username         string `json:"username"`
    ConfirmationCode string `json:"confirmation_code"`
}

type Credentials struct {
	Username    string `json:"username"`
	Password string `json:"password"`
}
