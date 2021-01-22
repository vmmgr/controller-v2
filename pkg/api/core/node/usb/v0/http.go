package v0

import (
	"encoding/json"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
	"github.com/vmmgr/controller/pkg/api/core/node/usb"
	"github.com/vmmgr/controller/pkg/api/core/tool/client"
	"strconv"
)

func httpRequest(ip string, port uint) ([]libvirtxml.NodeDevice, error) {
	var res usb.Node

	response, err := client.Get("http://"+ip+":"+strconv.Itoa(int(port))+"/api/v1/usb", "")
	if err != nil {
		return []libvirtxml.NodeDevice{}, err
	}

	if err = json.Unmarshal([]byte(response), &res); err != nil {
		return []libvirtxml.NodeDevice{}, err
	}

	return res.Data.USB, nil
}
