package hash

import "golang.org/x/crypto/bcrypt"

const bcryptCost = 10

func BCryptHash(s string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(s), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

func BCryptCompare(plainStr, hashedStr string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(plainStr))
	return err == nil
}
