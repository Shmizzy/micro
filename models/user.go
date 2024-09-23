package models

type User struct {
	UID 	  string `json:"uid"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
