package v0

import (
	"encoding/json"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/core/auth"
	requestInt "github.com/vmmgr/controller/pkg/api/core/request"
	request "github.com/vmmgr/controller/pkg/api/core/request/v0"
	"github.com/vmmgr/controller/pkg/api/core/tool/client"
	"github.com/vmmgr/controller/pkg/api/core/tool/gen"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	template "github.com/vmmgr/controller/pkg/api/core/vm/template/v0"
	"github.com/vmmgr/node/pkg/api/core/gateway"
	nodeNIC "github.com/vmmgr/node/pkg/api/core/nic"
	"github.com/vmmgr/node/pkg/api/core/storage"
	"github.com/vmmgr/node/pkg/api/core/tool/cloudinit"
	nodeVM "github.com/vmmgr/node/pkg/api/core/vm"
	"log"
	"strconv"
	"strings"
	"time"
)

type VMTemplateHandler struct {
	input    nodeVM.VirtualMachine
	template vm.Template
	node     core.Node
	authUser auth.GroupResult
	admin    bool
}

func NewVMAdminTemplateHandler(input VMTemplateHandler) *VMTemplateHandler {
	return &VMTemplateHandler{template: input.template, node: input.node, admin: true}
}

func NewVMUserTemplateHandler(input VMTemplateHandler) *VMTemplateHandler {
	return &VMTemplateHandler{template: input.template, node: input.node, authUser: input.authUser, admin: false}
}

func (t *VMTemplateHandler) templateApply() error {
	vmTemplate, vmTemplatePlan, err := template.GetTemplate(t.template.TemplateID, t.template.TemplatePlanID)
	if err != nil {
		return err
	}

	log.Println(vmTemplate, vmTemplatePlan)

	go func() {
		imaConResult, imageResult, err := extractImaCon(vmTemplate)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(imaConResult)
		uuid := gen.GenerateUUID()
		var path, nic, name string

		// 管理者側
		if t.admin {
			name = strconv.Itoa(0) + "-" + t.template.Name
			path = name + ".img"
			nic = t.template.NICType
		} else {
			// 管理者以外
			name = strconv.Itoa(int(t.authUser.Group.ID)) + "-" + gen.GenerateUUID()
			path = name + "-1.img"
			nic = "br190"
		}
		gid := uint(0)

		storageType := 0
		if !imageResult.Data.VirtIO {
			storageType = 11
		}

		//VM作成用のデータ
		virtualMachineTemplate := nodeVM.VirtualMachine{
			Info: gateway.Info{
				GroupID: gid,
				UUID:    uuid,
			},
			Name:    name,
			Memory:  vmTemplatePlan.Mem,
			CPUMode: 1, //host-model
			VCPU:    vmTemplatePlan.CPU,
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
			Storage: []storage.VMStorage{
				{
					Type:     uint(storageType),
					FileType: 0, //qcow2
					Path:     path,
					ReadOnly: false,
					Boot:     0,
				},
			},
			CloudInit: cloudinit.CloudInit{
				MetaData: cloudinit.MetaData{LocalHostName: t.template.Name},
				UserData: cloudinit.UserData{
					Password: t.template.Password,
					//ChPasswd:  "{ expire: False }",
					SshPwAuth: false,
				},
				NetworkConfig: cloudinit.NetworkCon{
					Config: []cloudinit.NetworkConfig{
						{
							Type:       "physical",
							Name:       "",
							MacAddress: "",
							Subnets: []cloudinit.NetworkConfigSubnet{
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
			CloudInitApply: imageResult.Data.CloudInit,
			Template: nodeVM.Template{
				Apply: true,
				Storage: storage.Storage{
					Mode: 1,
					FromImaCon: storage.ImaCon{
						IP:   imaConResult.IP,
						Path: imageResult.Data.Path,
					},
					Type:     10, // BootDisk(virtIO)
					FileType: 0,  // qcow2
					PathType: t.template.StoragePathType,
					Capacity: t.template.StorageCapacity,
					ReadOnly: false,
					Path:     path,
				},
			},
		}

		log.Println(t.template.PCI)

		if t.admin {
			if len(t.template.USB) > 0 {
				virtualMachineTemplate.USB = t.template.USB
			}
			if len(t.template.PCI) > 0 {
				virtualMachineTemplate.PCI = t.template.PCI
			}
		}

		body, _ := json.Marshal(virtualMachineTemplate)

		log.Println(string(body))

		resultVMCreateProcess, err := client.Post(
			"http://"+t.node.IP+":"+strconv.Itoa(int(t.node.Port))+"/api/v1/vm", body)
		if err != nil {
			log.Println(err)
			return
		}

		//timeout時間(分)
		request.Add(requestInt.Request{ExpirationDate: time.Now().Add(20 * time.Minute), GroupID: gid, UUID: uuid})

		log.Println(resultVMCreateProcess)

		//DB追加
		//if t.admin {
		//	dbVM.Create(&core.VM{NodeID: t.input.NodeID, GroupID: 0, Name: t.input.Name, UUID: uuid})
		//} else {
		//	dbVM.Create(&core.VM{NodeID: t.input.NodeID, GroupID: t.authUser.Group.ID, Name: t.input.Name, UUID: uuid})
		//}
	}()
	return nil
}
