package v0

import (
	"encoding/json"
	"github.com/vmmgr/controller/pkg/api/core/tool/client"
	nodeNIC "github.com/vmmgr/node/pkg/api/core/nic"
	"github.com/vmmgr/node/pkg/api/core/storage"
	"github.com/vmmgr/node/pkg/api/core/tool/cloudinit"
	nodeVM "github.com/vmmgr/node/pkg/api/core/vm"
	"testing"
)

func Test(t *testing.T) {
	body, _ := json.Marshal(nodeVM.VirtualMachine{
		Name:    "test101",
		Memory:  512,
		CPUMode: 0,
		VCPU:    1,
		NIC: []nodeNIC.NIC{
			{
				Type:   0,
				Driver: 0,
				Mode:   0,
				MAC:    "",
				Device: "br100",
			},
		},
		VNCPort: 0, //VNCポートをNode側で自動生成
		Storage: []storage.VMStorage{
			{
				Type:     10, // BootDisk(virtIO)
				FileType: 0,  //qcow2
				Path:     "/home/yonedayuto/vm/1/10fef5d71-bd37-4a31-95c5-2ac7734a0b23-1",
				ReadOnly: false,
				Boot:     0,
			},
		},
		CloudInit: cloudinit.CloudInit{
			MetaData: cloudinit.MetaData{LocalHostName: "test"},
			UserData: cloudinit.UserData{
				Password:  "test123",
				ChPasswd:  "{ expire: False }",
				SshPwAuth: true,
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
								Address: "172.40.0.124",
								Netmask: "255.255.252.0",
								Gateway: "172.40.0.1",
								DNS:     []string{"1.1.1.1"},
							},
						},
					},
				},
			},
		},
		CloudInitApply: true,
	})

	resultVMCreateProcess, err := client.Post(
		"http://127.0.0.1:8080/api/v1/vm",
		body)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(resultVMCreateProcess)
}
