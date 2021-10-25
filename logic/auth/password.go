package auth

import "github.com/aldy505/phc-crypto/argon2"

func GeneratePassword(password string) (string, error) {
	hash, err := argon2.Hash(password, argon2.Config{})
	if err != nil {
		return "", err
	}
	return hash, nil
}

func VerifyPassword(plain, hashed string) (bool, error) {
	verify, err := argon2.Verify(hashed, plain)
	if err != nil {
		return false, err
	}
	return verify, nil
}
