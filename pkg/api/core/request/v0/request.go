package v0

import (
	"fmt"
	request "github.com/vmmgr/controller/pkg/api/core/request"
	"log"
	"time"
)

var req = make(map[string]request.Request, 1000)

func Add(input request.Request) error {
	// ローカル変数のReqにinputを代入
	req[input.UUID] = input

	return nil
}

func Delete(uuid string) error {
	// ローカル変数のReqにinputを代入
	delete(req, uuid)

	return nil
}

func Get(uuid string) (request.Request, error) {
	// ローカル変数のReqにinputを代入
	data, result := req[uuid]
	if result {
		return data, nil
	} else {
		return request.Request{}, fmt.Errorf("not found... ")
	}
}

func AutoDelete() {
	timeNow := time.Now()
	for key, value := range req {
		log.Println(key, value)
		if value.ExpirationDate.Unix() < timeNow.Unix() {
			delete(req, key)
		}
	}
}
