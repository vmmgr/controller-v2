package etc

import (
	"github.com/google/uuid"
	"log"
)

func GenerateUUID() string {
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Println(err)
	}

	return uuid.String()
}

func GenerateToken() string {
	u1, err := uuid.NewRandom()
	if err != nil {
		log.Println(err)
	}
	u2, err := uuid.NewRandom()
	if err != nil {
		log.Println(err)
	}

	uuid := u1.String() + u2.String()
	log.Println("GenerateToken:  " + uuid)
	return uuid
}
