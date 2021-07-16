package v2

import (
	"github.com/libvirt/libvirt-go"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	cloudinitInt "github.com/vmmgr/controller/pkg/api/core/vm/cloudinit"
	cloudinit "github.com/vmmgr/controller/pkg/api/core/vm/cloudinit/v0"
	nodeNIC "github.com/vmmgr/controller/pkg/api/core/vm/nic"
	storageInt "github.com/vmmgr/controller/pkg/api/core/vm/storage"
	storage "github.com/vmmgr/controller/pkg/api/core/vm/storage/v0"
	dbTemplatePlan "github.com/vmmgr/controller/pkg/api/store/imacon/template_plan/v0"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
)

// key: imacon_id value: {time}
var LoadCopy = make(map[string]VMHandler)

type VMTemplateHandler struct {
	uuid     string
	input    vm.VirtualMachine
	template vm.Template
	node     core.Node
	storage  core.Storage
	conn     *libvirt.Connect
	admin    bool
	groupID  uint
	ipID     uint
	ctrlType uint // 1:Admin 2:User
}

func NewVMAdminTemplateHandler(input VMTemplateHandler) *VMTemplateHandler {
	return &VMTemplateHandler{
		uuid:     input.uuid,
		template: input.template,
		node:     input.node,
		storage:  input.storage,
		admin:    true,
		conn:     input.conn,
		ctrlType: 1,
	}
}

func NewVMUserTemplateHandler(input VMTemplateHandler) *VMTemplateHandler {
	return &VMTemplateHandler{
		uuid:     input.uuid,
		template: input.template,
		node:     input.node,
		storage:  input.storage,
		admin:    false,
		conn:     input.conn,
		groupID:  input.groupID,
		ipID:     input.ipID,
		ctrlType: 2,
	}
}

func (t *VMTemplateHandler) templateApply() error {
	log.Println(t.template.TemplatePlanID)
	resultTemplatePlan, err := dbTemplatePlan.Get(core.TemplatePlan{Model: gorm.Model{ID: t.template.TemplatePlanID}})
	if err != nil {
		log.Println(err)
		return err
	}

	var path, nic, name string

	// 管理者側
	if t.admin {
		name = strconv.Itoa(0) + "-" + t.template.Name
		path = t.storage.Path + "/" + name + "/" + "1.img"
		nic = t.template.NICType
	} else {
		// 管理者以外
		name = strconv.Itoa(int(t.groupID)) + "-" + t.template.Name
		path = t.storage.Path + "/" + name + "/" + "1.img"
		nic = t.node.PrimaryNIC
	}

	storageh := storage.NewStorageHandler(storage.StorageHandler{
		UUID:      t.uuid,
		Conn:      t.conn,
		Input:     storageInt.Storage{},
		VM:        vm.VirtualMachine{},
		Address:   nil,
		Auth:      nil,
		SrcImaCon: *resultTemplatePlan[0].Template.Image.ImaCon,
		DstAuth: storageInt.Auth{
			User: t.node.User,
			IP:   t.node.IP,
			Port: t.node.Port,
		},
		SrcPath:  resultTemplatePlan[0].Template.Image.Path,
		DstPath:  path,
		CtrlType: t.ctrlType,
	})

	//VM作成用のデータ
	virtualMachineTemplate := vm.VirtualMachine{
		Name:    name,
		Memory:  resultTemplatePlan[0].Mem,
		CPUMode: 1, //host-model
		VCPU:    resultTemplatePlan[0].CPU,
		NIC: []nodeNIC.NIC{
			{
				Type:   0,
				Driver: 0,
				Mode:   0,
				MAC:    "",
				Device: nic,
			},
		},
		VNCPort: 0, //VNCポートをNode側で自動生成
		Storage: []storageInt.VMStorage{
			{
				Type:     uint(0),
				FileType: 0, //qcow2
				Path:     path,
				ReadOnly: false,
				Boot:     0,
			},
		},
		CloudInit: cloudinit.CloudInit{
			MetaData: cloudinitInt.MetaData{LocalHostName: t.template.Name},
			UserData: cloudinitInt.UserData{
				Password: t.template.Password,
				//ChPasswd:  "{ expire: False }",
				SshPwAuth: false,
			},
			NetworkConfig: cloudinitInt.NetworkCon{
				Config: []cloudinitInt.NetworkConfig{
					{
						Type:       "physical",
						Name:       "",
						MacAddress: "",
						Subnets: []cloudinitInt.NetworkConfigSubnet{
							{
								Type:    "static",
								Address: t.template.IP,
								Netmask: t.template.NetMask,
								Gateway: t.template.Gateway,
								DNS:     strings.Split(t.template.DNS, ","),
							},
						},
					},
				},
			},
		},
		CloudInitApply: *resultTemplatePlan[0].Template.Image.CloudInit,
		Template: vm.TemplateVM{
			Apply: true,
			Storage: storageInt.Storage{
				Mode:     1,
				Type:     10, // BootDisk(virtIO)
				FileType: 0,  // qcow2
				PathType: t.template.StoragePathType,
				Capacity: t.template.StorageCapacity,
				ReadOnly: false,
				Path:     path,
			},
		},
	}

	LoadCopy[t.uuid] = VMHandler{
		Conn:    t.conn,
		VM:      virtualMachineTemplate,
		Node:    t.node,
		GroupID: t.groupID,
		IPID:    t.ipID,
	}

	// Image Copy

	go func() {
		log.Println(resultTemplatePlan[0].Template.Image.ImaCon.IP)
		log.Println(t.node.IP)

		//log.Println(t.template.PCI)
		//
		//if t.admin {
		//	if len(t.template.USB) > 0 {
		//		virtualMachineTemplate.USB = t.template.USB
		//	}
		//	if len(t.template.PCI) > 0 {
		//		virtualMachineTemplate.PCI = t.template.PCI
		//	}
		//}

		err = storageh.AddFromImage()
		if err != nil {
			log.Println(err)
		}

	}()
	return nil
}
