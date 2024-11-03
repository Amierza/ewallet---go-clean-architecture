package helpers

import "golang.org/x/crypto/bcrypt"

func HashPin(pin string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pin), 4)
	return string(bytes), err
}

func ChcekPin(hashPin string, plainPin []byte) (bool, error) {
	hashP := []byte(hashPin)
	if err := bcrypt.CompareHashAndPassword(hashP, plainPin); err != nil {
		return false, err
	}
	return true, nil
}
