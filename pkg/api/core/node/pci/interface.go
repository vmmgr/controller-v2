package pci

import libvirtxml "github.com/libvirt/libvirt-go-xml"

type Node struct {
	Status uint `json:"status"`
	Data   struct {
		PCI []libvirtxml.NodeDevice `json:"pci"`
	} `json:"data"`
}

type Result struct {
	Status bool                    `json:"status"`
	Error  string                  `json:"error"`
	PCI    []libvirtxml.NodeDevice `json:"pci"`
}

type OneResult struct {
	Status bool                  `json:"status"`
	Error  string                `json:"error"`
	PCI    libvirtxml.NodeDevice `json:"pci"`
}
