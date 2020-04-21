package etc

import (
	"fmt"
	"github.com/google/uuid"
)

func GenerateToken() string {
	u, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
	}
	uuid := u.String()
	fmt.Println("GenerateToken: " + uuid)
	return uuid
}
