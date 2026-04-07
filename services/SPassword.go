package services

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password),14)
	if err!=nil {
		return "",fmt.Errorf("error to create password_hash: %w",err)
	}
	return string(bytes),nil
}

func ChechPassword(password,hash string)bool{
	err:=bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	return err==nil
}