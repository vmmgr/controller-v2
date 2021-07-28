package v2

import (
	"encoding/xml"
	"github.com/libvirt/libvirt-go"
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"log"
)

func Startup(nodeID uint, uuid string) (*vm.Detail, error) {
	_, conn, err := connectLibvirt(nodeID)
	if err != nil {
		return nil, err
	}

	dom, err := conn.LookupDomainByUUIDString(uuid)
	if err != nil {
		return nil, err
	}

	stat, _, err := dom.GetState()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if stat != libvirt.DOMAIN_RUNNING {
		if err = dom.Create(); err != nil {
			log.Println(err)
			return nil, err
		}
	}

	t := libVirtXml.Domain{}
	stat, _, _ = dom.GetState()
	xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	xml.Unmarshal([]byte(xmlString), &t)

	err = dom.Free()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &vm.Detail{VM: t, Stat: uint(stat)}, nil
}

func Shutdown(nodeID uint, uuid string, force bool) (*vm.Detail, error) {
	_, conn, err := connectLibvirt(nodeID)
	if err != nil {
		return nil, err
	}

	dom, err := conn.LookupDomainByUUIDString(uuid)
	if err != nil {
		return nil, err
	}

	stat, _, err := dom.GetState()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if stat != libvirt.DOMAIN_SHUTOFF {
		// Forceがtrueである場合、強制終了
		if force {
			if err = dom.Destroy(); err != nil {
				log.Println(err)
				return nil, err
			}
		} else {
			if err = dom.Shutdown(); err != nil {
				log.Println(err)
				return nil, err
			}
		}
	}

	t := libVirtXml.Domain{}
	stat, _, _ = dom.GetState()
	xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	xml.Unmarshal([]byte(xmlString), &t)

	err = dom.Free()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &vm.Detail{VM: t, Stat: uint(stat)}, nil
}

func Reset(nodeID uint, uuid string) (*vm.Detail, error) {
	_, conn, err := connectLibvirt(nodeID)
	if err != nil {
		return nil, err
	}

	dom, err := conn.LookupDomainByUUIDString(uuid)
	if err != nil {
		return nil, err
	}

	if err = dom.Reset(0); err != nil {
		return nil, err
	}

	t := libVirtXml.Domain{}
	stat, _, _ := dom.GetState()
	xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	xml.Unmarshal([]byte(xmlString), &t)

	err = dom.Free()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &vm.Detail{VM: t, Stat: uint(stat)}, nil
}
