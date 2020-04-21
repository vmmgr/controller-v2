package etc

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func Hashgenerate(data string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

func Verifyhashdata(data string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(data)) == nil
}
