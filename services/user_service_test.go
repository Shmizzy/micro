package services

import (
	"service/models"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	user := models.User{
		Username:    "izzylppz",
		Email:       "izzylppz@gmail.com",
		Password:    "Deeboo49",
		PhoneNumber: "+19568623733",
	}

	err := RegisterUser(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

/*
	 func TestRegisterUser_InvalidEmail(t *testing.T) {
		user := models.User{
			Username: "testuser",
			Email:    "invalid-email",
			Password: "Test12345",
		}

		err := RegisterUser(user)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	}
*/
/* func TestLoginUser(t *testing.T) {
	credentials := models.Credentials{
		Username: "testuser",
		Password: "Test12345",
	}

	token, err := LoginUser(credentials)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Fatalf("expected token, got empty string")
	}
}
 */
/* func TestLoginUser_InvalidCredentials(t *testing.T) {
	credentials := models.Credentials{
		Email:    "test@example.com",
		Password: "WrongPassword",
	}

	_, err := LoginUser(credentials)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
*/
