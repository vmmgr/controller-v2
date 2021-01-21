package v0

import (
	"encoding/json"
	"github.com/vmmgr/controller/pkg/api/core/tool/client"
	"strconv"
)

func Get(ip string, port uint, uuid string) (nodeOneVMResponse, error) {
	var res nodeOneVMResponse

	response, err := client.Get("http://"+ip+":"+strconv.Itoa(int(port))+"/api/v1/vm/"+uuid, "")
	if err != nil {
		return res, err
	}

	if json.Unmarshal([]byte(response), &res) != nil {
		return res, err
	}
	return res, nil
}
