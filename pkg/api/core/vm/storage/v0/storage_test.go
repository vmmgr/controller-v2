package v0

import (
	"github.com/vmmgr/controller/pkg/api/core/tool/ssh"
	"log"
	"testing"
)

func Test1(t *testing.T) {
	sh := ssh.Auth{
		IP:   "192.168.22.132",
		Port: 22,
		User: "yonedayuto",
	}
	//qemu-img create -f qcow2 file.qcow2 100M
	command := "/home/yonedayuto/imacon copy --uuid 2e684c62-d680-40f9-818b-2919ca02507e --url http://localhost:8081/api/v1/controller --src /home/yonedayuto/image/focal-server-cloudimg-amd64-disk-kvm.img --dst /home/yonedayuto/Documents/vmmgr/vm-image/test1.img --addr 192.168.22.1:22 --user yonedayuto --config /home/yonedayuto/config.json"
	log.Println(command)
	result, err := sh.SSHClientExecCmd(command)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)

}
