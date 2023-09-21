package users

import "testing"

func TestGenerateJWT(t *testing.T) {
	email := "jasonhorks@gmail.com"
	_, err := GenerateJWT(email)
	if err != nil {
		t.Errorf("failure")
	}
}

func TestValidateJWT(t *testing.T) {
	realEmail := "jasonhorks@gmail.com"
	jwt, err := GenerateJWT(realEmail)
	if err != nil {
		t.Errorf("failed to generate jwt")
	}

	email, err := ValidateJWT(jwt)

	if err != nil {
		t.Errorf("error occured when validating jwt: %v", err)
	}

	if email != realEmail {
		t.Errorf("failed to validate jwt. Email should be %v, came out %v", realEmail, email)
	}

}
