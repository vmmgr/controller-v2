package v0

import (
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	image "github.com/vmmgr/controller/pkg/api/core/tool/image"
	ssh "github.com/vmmgr/controller/pkg/api/core/tool/ssh"
	"testing"
)

var url = "http://tinycorelinux.net/12.x/x86/release/TinyCore-current.iso"
var dstPath1 = "/home/ubuntu/test.iso"
var dstPath2 = "/home/ubuntu/test.iso"

var nodeID uint = 2

func TestConfigApply(t *testing.T) {
	if err := config.GetConfig("../../../../../../cmd/backend/con.json"); err != nil {
		t.Fatal(err)
	}
}

func TestGet1CDROM(t *testing.T) {
	conn, err := ssh.ConnectSSH(2)
	if err != nil {
		t.Fatal(err)
	}

	defer conn.Close()

	imageHandler := image.ImageHandler{
		SSHClient: conn,
		DstPath:   dstPath1,
	}
	imageHandler.ImageDownload(url)

	imageHandler = image.ImageHandler{
		SSHClient: conn,
		DstPath:   dstPath2,
	}
	imageHandler.ImageDownload(url)

}
