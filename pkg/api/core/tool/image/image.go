package v0

import (
	"github.com/pkg/sftp"
	"github.com/schollz/progressbar"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net/http"
	"time"
)

type Progress struct {
	total int64
	size  int64
}

type ImageHandler struct {
	SSHClient *ssh.Client
	DstPath   string
}

func (p *Progress) Write(data []byte) (int, error) {
	n := len(data)
	p.size += int64(n)

	return n, nil
}

func (h ImageHandler) ImageDownload(url string) error {

	log.Println(h)
	log.Println(url)

	defer h.SSHClient.Close()

	//conn, err := ssh.ConnectSSH(nodeID)
	//if err != nil {
	//	t.Fatal(err)conn
	//}
	//
	//defer conn.Close()

	// SFTP Client
	client, err := sftp.NewClient(h.SSHClient)
	if err != nil {
		return err
	}
	defer client.Close()

	// dstFileの作成(ssh先)
	dstFile, err := client.Create(h.DstPath)
	if err != nil {
		log.Println(err)
	}

	defer dstFile.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println(resp.ContentLength)

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
				//log.Println(p.size)
			} else {
				return
			}
		}
	}()

	// Write the body to file
	_, err = io.Copy(dstFile, io.TeeReader(resp.Body, &p))
	if err != nil {
		return err
	}
	return nil
}
