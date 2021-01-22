package v0

import (
	"encoding/json"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
	"github.com/vmmgr/controller/pkg/api/core/node/pci"
	"github.com/vmmgr/controller/pkg/api/core/tool/client"
	"log"
	"strconv"
)

func httpRequest(ip string, port uint) ([]libvirtxml.NodeDevice, error) {
	var res pci.Node

	response, err := client.Get("http://"+ip+":"+strconv.Itoa(int(port))+"/api/v1/pci", "")
	if err != nil {
		return []libvirtxml.NodeDevice{}, err
	}

	log.Println(response)

	if err = json.Unmarshal([]byte(response), &res); err != nil {
		return []libvirtxml.NodeDevice{}, err
	}

	return res.Data.PCI, nil
}

func httpRequest1(ip string, port uint) (string, error) {
	response, err := client.Get("http://"+ip+":"+strconv.Itoa(int(port))+"/api/v1/pci", "")
	if err != nil {
		return "", err
	}

	return response, nil
}
