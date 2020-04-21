package data

import (
	"fmt"
	"google.golang.org/grpc/metadata"
	"strconv"
	"strings"
)

type AuthInData struct {
	User  string
	Pass  string
	Token string
}

func GenerateAddress(ip string, port int) string {
	return ip + ":" + strconv.Itoa(port)
}

func AuthDataExtraction(data metadata.MD) string {
	if Token, ok := data["authorization"]; ok {
		d := strings.Split(Token[0], " ")
		fmt.Println("RestAPI Token: " + d[1])
		return d[1]
	}
	return ""
}
