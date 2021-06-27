package v2

import (
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	"github.com/vmmgr/controller/pkg/api/core/tool/gen"
	"github.com/vmmgr/controller/pkg/api/core/tool/ssh"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"github.com/vmmgr/controller/pkg/api/core/vm/cloudinit"
	"github.com/vmmgr/controller/pkg/api/core/vm/cloudinit/v0"
	nic "github.com/vmmgr/controller/pkg/api/core/vm/nic/v0"
	storageInterface "github.com/vmmgr/controller/pkg/api/core/vm/storage"
	storage "github.com/vmmgr/controller/pkg/api/core/vm/storage/v0"
	//"github.com/vmmgr/node/pkg/api/core/tool/file"
	"log"
	"strconv"
)

func (h *VMHandler) xmlGenerate() (*libVirtXml.Domain, error) {

	uuid := gen.GenerateUUID()

	// nic xmlの生成
	hNIC := nic.NewNICHandler(nic.NICHandler{
		Conn:    h.Conn,
		VM:      h.VM,
		BaseMAC: h.Node.Mac,
		Address: &vm.Address{PCICount: 0, DiskCount: 0},
	})
	nics, err := hNIC.XmlGenerate()
	if err != nil {
		return nil, err
	}

	// CloudInit周りの処理
	if h.VM.CloudInitApply {
		directory := h.Node.Storage[0].Path + "/" + h.VM.Name

		//if !file.ExistsCheck(directory) {
		//	if err = os.Mkdir(directory, 0755); err != nil {
		//		log.Println(err)
		//		return nil, err
		//	}
		//}

		for i, a := range nics {
			h.VM.CloudInit.NetworkConfig.Config[i].MacAddress = a.MAC.Address
			h.VM.CloudInit.NetworkConfig.Config[i].Name = "eth" + strconv.Itoa(i)
		}

		log.Println(h.VM.CloudInit.NetworkConfig)
		log.Println(h.VM.CloudInit.UserData)

		hCloudInit := v0.NewCloudInitHandler(v0.CloudInit{
			DirPath: directory,
			Auth: ssh.Auth{
				IP:   h.Node.IP,
				Port: h.Node.Port,
				User: h.Node.User,
				Pass: h.Node.Pass,
			},
			MetaData:      cloudinit.MetaData{InstanceID: h.Node.Name, LocalHostName: h.VM.Name},
			UserData:      h.VM.CloudInit.UserData,
			NetworkConfig: h.VM.CloudInit.NetworkConfig,
		})

		err = hCloudInit.Generate()
		if err != nil {
			log.Println(err)
			return nil, err
		}

		h.VM.Storage = append(h.VM.Storage, storageInterface.VMStorage{
			Type:     1,
			Path:     directory + "/cloudinit.img",
			ReadOnly: true,
		})
	}

	log.Println(hNIC.Address)

	// storage xmlの生成
	hStorage := storage.NewStorageHandler(storage.StorageHandler{
		Conn:    h.Conn,
		VM:      h.VM,
		Address: hNIC.Address,
	})
	disks, err := hStorage.XmlGenerate()
	if err != nil {
		return nil, err
	}

	domCfg := &libVirtXml.Domain{
		Type: "kvm",
		Memory: &libVirtXml.DomainMemory{
			Value:    h.VM.Memory,
			Unit:     "MB",
			DumpCore: "on",
		},
		VCPU:        &libVirtXml.DomainVCPU{Value: h.VM.VCPU},
		CPU:         &libVirtXml.DomainCPU{Mode: getCPUMode(h.VM.CPUMode)},
		UUID:        uuid,
		Name:        h.VM.Name,
		Title:       h.VM.Name,
		Description: h.VM.Name,
		Features: &libVirtXml.DomainFeatureList{
			ACPI: &libVirtXml.DomainFeature{},
			APIC: &libVirtXml.DomainFeatureAPIC{},
		},
		OS: &libVirtXml.DomainOS{
			BootDevices: []libVirtXml.DomainBootDevice{{Dev: "hd"}},
			Kernel:      "",
			//Initrd:  "/home/markus/workspace/worker-management/centos/kvm-centos.ks",
			//Cmdline: "ks=file:/home/markus/workspace/worker-management/centos/kvm-centos.ks method=http://repo02.agfa.be/CentOS/7/os/x86_64/",
			Type: &libVirtXml.DomainOSType{
				Arch:    getArchConvert(h.VM.OS.Arch),
				Machine: h.Node.Machine,
				Type:    "hvm",
			},
		},
		Devices: &libVirtXml.DomainDeviceList{
			Emulator: h.Node.Emulator,
			Inputs: []libVirtXml.DomainInput{
				{Type: "mouse", Bus: "ps2"},
				{Type: "keyboard", Bus: "ps2"},
			},
			Graphics: []libVirtXml.DomainGraphic{
				{
					VNC: &libVirtXml.DomainGraphicVNC{
						Port:      int(h.VM.VNCPort),
						WebSocket: int(h.VM.WebSocketPort),
						Keymap:    h.VM.KeyMap,
						Listen:    "0.0.0.0",
					},
				},
			},
			Videos: []libVirtXml.DomainVideo{
				{
					Model: libVirtXml.DomainVideoModel{
						Type:    "qxl",
						Heads:   1,
						Ram:     65536,
						VRam:    65536,
						VGAMem:  16384,
						Primary: "yes",
					},
					Alias: &libVirtXml.DomainAlias{Name: "video0"},
				},
			},
			Disks:      disks,
			Interfaces: nics,
		},
	}

	return domCfg, nil
}
