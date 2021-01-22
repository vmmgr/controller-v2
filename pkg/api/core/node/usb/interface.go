package usb

import libvirtxml "github.com/libvirt/libvirt-go-xml"

type Node struct {
	Status uint `json:"status"`
	Data   struct {
		USB []libvirtxml.NodeDevice `json:"usb"`
	} `json:"data"`
}

type Result struct {
	Status bool                    `json:"status"`
	Error  string                  `json:"error"`
	USB    []libvirtxml.NodeDevice `json:"usb"`
}

type OneResult struct {
	Status bool                  `json:"status"`
	Error  string                `json:"error"`
	USB    libvirtxml.NodeDevice `json:"usb"`
}
