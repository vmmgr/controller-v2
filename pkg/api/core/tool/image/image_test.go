package v0

import (
	"github.com/pkg/sftp"
	"github.com/schollz/progressbar"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"github.com/vmmgr/controller/pkg/api/core/tool/ssh"
	"io"
	"log"
	"net/http"
	"testing"
	"time"
)

//var url = "https://releases.ubuntu.com/20.04.2.0/ubuntu-20.04.2.0-desktop-amd64.iso"
var url = "http://tinycorelinux.net/12.x/x86/release/TinyCore-current.iso"
var dstPath = "/home/ubuntu/test.iso"
var nodeID uint = 2

func TestConfigApply(t *testing.T) {
	if err := config.GetConfig("../../../../../../cmd/backend/con.json"); err != nil {
		t.Fatal(err)
	}
}

func TestImageDownload(t *testing.T) {

	conn, err := ssh.ConnectSSHNodeID(nodeID)
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

	// dstFileの作成(ssh先)
	dstFile, err := client.Create(dstPath)
	if err != nil {
		log.Println(err)
	}
	defer dstFile.Close()

	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Log(resp.ContentLength)

	p := Progress{total: resp.ContentLength}

	count := 100
	count64 := int64(count)
	bar := progressbar.Default(count64)

	// Node側の表示
	go func() {
		for {
			if p.size != p.total {
				<-time.NewTimer(200 * time.Microsecond).C
				bar.Set(int(float64(p.size) / float64(p.total) * 100))
			} else {
				return
			}
		}
	}()
	// Write the body to file

	_, err = io.Copy(dstFile, io.TeeReader(resp.Body, &p))
	if err != nil {
		t.Fatal(err)
	}
}
