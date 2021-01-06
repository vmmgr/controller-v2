package v0

import (
	"encoding/json"
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core/auth"
	"github.com/vmmgr/controller/pkg/api/core/node"
	"github.com/vmmgr/controller/pkg/api/core/tool/client"
	"github.com/vmmgr/controller/pkg/api/core/tool/gen"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	template "github.com/vmmgr/controller/pkg/api/core/vm/template/v0"
	dbVM "github.com/vmmgr/controller/pkg/api/store/vm/v0"
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
	input    vm.Template
	node     node.Node
	authUser auth.GroupResult
	admin    bool
}

func NewVMAdminTemplateHandler(input VMTemplateHandler) *VMTemplateHandler {
	return &VMTemplateHandler{input: input.input, node: input.node, admin: true}
}

func NewVMUserTemplateHandler(input VMTemplateHandler) *VMTemplateHandler {
	return &VMTemplateHandler{input: input.input, node: input.node, authUser: input.authUser, admin: false}
}

func (t *VMTemplateHandler) templateApply() error {
	vmTemplate, vmTemplatePlan, err := template.GetTemplate(t.input.TemplateID, t.input.TemplatePlanID)
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
		var path string
		var nic string
		var name string

		if t.admin {
			name = strconv.Itoa(0) + "-" + t.input.Name
			path = name + "-1.img"
			nic = t.input.NICType
		} else {
			name = strconv.Itoa(int(t.authUser.Group.ID)) + "-" + gen.GenerateUUID()
			path = name + "-1.img"
			nic = "br190"
		}
		gid := uint(0)
		// Storage作成用にbodyを作成する
		createStorageBody, _ := json.Marshal(storage.Storage{
			Mode: 1,
			FromImaCon: storage.ImaCon{
				IP:   imaConResult.IP,
				Path: imageResult.Data.Path,
			},
			Type:     10, // BootDisk(virtIO)
			FileType: 0,  // qcow2
			GroupID:  gid,
			UUID:     uuid,
			PathType: t.input.StoragePathType,
			Capacity: t.input.StorageCapacity,
			ReadOnly: false,
			Path:     path,
		})

		log.Println(string(createStorageBody))

		resultStorageProcess, err := client.Post(
			"http://"+t.node.IP+":"+strconv.Itoa(int(t.node.Port))+"/api/v1/storage",
			createStorageBody)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(resultStorageProcess)

		timer := time.NewTimer(20 * time.Minute)
		defer timer.Stop()

		//Todo 取りこぼす可能性があるので、要調査
	L:
		for {
			select {
			//20分以上かかる場合はタイムアウトさせる
			case <-timer.C:
				log.Println("Error: timeout")
				err = fmt.Errorf("Error: timeout ")
				break L
				//UUIDとGroupIDがMatchし、Progressが100の場合、storage転送処理が終了
			case msg := <-vm.Broadcast:
				if msg.UUID == uuid && msg.GroupID == gid && msg.Progress == 100 {
					//path変数にnode側のストレージをフルパスで代入する
					path = msg.FilePath
					err = nil
					break L
				}
			}
		}

		// Errorが発生した場合
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("End: copy storage")

		//VM作成用のデータ
		body, _ := json.Marshal(nodeVM.VirtualMachine{
			Name:    name,
			Memory:  vmTemplatePlan.Mem,
			CPUMode: 0,
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
					Type:     10, // BootDisk(virtIO)
					FileType: 0,  //qcow2
					Path:     path,
					ReadOnly: false,
					Boot:     0,
				},
			},
			CloudInit: cloudinit.CloudInit{
				MetaData: cloudinit.MetaData{LocalHostName: t.input.Name},
				UserData: cloudinit.UserData{
					Password: t.input.Password,
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
									Address: t.input.IP,
									Netmask: t.input.NetMask,
									Gateway: t.input.Gateway,
									DNS:     strings.Split(t.input.DNS, ","),
								},
							},
						},
					},
				},
			},
			CloudInitApply: imageResult.Data.CloudInit,
		})

		resultVMCreateProcess, err := client.Post(
			"http://"+t.node.IP+":"+strconv.Itoa(int(t.node.Port))+"/api/v1/vm",
			body)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(resultVMCreateProcess)

		//DB追加
		if t.admin {
			dbVM.Create(&vm.VM{NodeID: t.input.NodeID, GroupID: 0, Name: t.input.Name, UUID: uuid})
		} else {
			dbVM.Create(&vm.VM{NodeID: t.input.NodeID, GroupID: t.authUser.Group.ID, Name: t.input.Name, UUID: uuid})
		}
	}()
	return nil
}
