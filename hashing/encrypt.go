package hashing


import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

 func HashPassword(password string) (string, error) {
        bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
        return string(bytes), err

 }

 func checkPasswordHash(password, hash string) bool {
    err :=bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
 }
func Check() {

    password :="secret"
    hash, _ :=HashPassword(password)

    fmt.Println("Paswword", password)
    fmt.Println("hash ", hash)

    match := checkPasswordHash(password, hash)

    fmt.Println("Match:  ", match)
    
}


