package v0

import (
	"github.com/pkg/sftp"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	ssh "github.com/vmmgr/controller/pkg/api/core/tool/ssh"
	"regexp"
	"strconv"
	"testing"
)

var url = "http://tinycorelinux.net/12.x/x86/release/TinyCore-current.iso"
var dstPath1 = "/home/ubuntu/test.iso"
var dstPath2 = "/home/ubuntu/test.iso"
var searchPath = "/home/ubuntu/iso"
var nodeID uint = 2

func TestConfigApply(t *testing.T) {
	if err := config.GetConfig("../../../../../../../cmd/backend/con.json"); err != nil {
		t.Fatal(err)
	}
}

func TestListCDROM(t *testing.T) {
	conn, err := ssh.ConnectSSHNodeID(2)
	if err != nil {
		t.Fatal(err)
	}

	defer conn.Close()

	// SFTP Client
	client, err := sftp.NewClient(conn)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	t.Log("BasePath: " + searchPath)
	//client.
	w := client.Walk(searchPath)
	for w.Step() {
		if w.Err() != nil {
			continue
		}
		match, _ := regexp.MatchString("^"+searchPath+"$", w.Path())
		if !match {
			t.Log("Name: " + w.Path()[len(searchPath)+1:] + ", " + strconv.Itoa(int(w.Stat().Size())) +
				" Byte, Time: " + w.Stat().ModTime().String())
		}
	}
}

//func TestGet1CDROM(t *testing.T) {
//	conn, err := ssh.ConnectSSHNodeID(2)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	defer conn.Close()
//
//	imageHandler := image.ImageHandler{
//		SSHClient: conn,
//		DstPath:   dstPath1,
//	}
//	imageHandler.ImageDownload(url)
//
//	imageHandler = image.ImageHandler{
//		SSHClient: conn,
//		DstPath:   dstPath2,
//	}
//	imageHandler.ImageDownload(url)
//
//}
