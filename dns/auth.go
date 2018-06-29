package dns

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func IsUser(user string) bool {
	// Get pwd file if there is one
	// Hash pw and compare to pw from dnsimple
	rec, err := net.LookupTXT(fmt.Sprintf("%s.%s", user, DomainName))
	if err != nil {
		// TODO make password record
		return true
	}
	pw := os.Getenv("PASSWORD")
	return CheckPasswordHash(pw, rec[0])

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
