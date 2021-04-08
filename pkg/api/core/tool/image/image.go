package v0

import (
	"github.com/pkg/sftp"
	"github.com/schollz/progressbar"
	"github.com/vmmgr/controller/pkg/api/core/node/storage"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net/http"
	"time"
)

var ImageDownloadProcess = false

type Progress struct {
	total int64
	size  int64
}

type ImageHandler struct {
	SSHClient *ssh.Client
	DstPath   string
	UUID      string
}

func (p *Progress) Write(data []byte) (int, error) {
	n := len(data)
	p.size += int64(n)

	return n, nil
}

func (h ImageHandler) ImageDownload(url string) error {
	ImageDownloadProcess = true

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
		ImageDownloadProcess = false
		return err
	}
	defer client.Close()

	// dstFileの作成(ssh先)
	dstFile, err := client.Create(h.DstPath)
	if err != nil {
		ImageDownloadProcess = false
		log.Println(err)
	}

	defer dstFile.Close()

	resp, err := http.Get(url)
	if err != nil {
		ImageDownloadProcess = false
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

	// Client側に通知
	go func() {
		for {
			if p.size != p.total {
				<-time.NewTimer(1 * time.Second).C
				storage.ClientBroadcast <- storage.WebSocketResult{
					NodeID:    0,
					Name:      "[Image Download]+URL: " + url,
					Err:       "",
					CreatedAt: time.Time{},
					Status:    0,
					Code:      0,
					FilePath:  "",
					Admin:     false,
					Message:   "",
					Progress:  uint(float64(p.size) / float64(p.total) * 100),
					UUID:      h.UUID,
				}
			} else {
				storage.ClientBroadcast <- storage.WebSocketResult{
					Name:      "[Image Download]+URL: " + url,
					CreatedAt: time.Time{},
					FilePath:  "",
					Admin:     false,
					Message:   "Finish!!",
					Progress:  100,
					UUID:      h.UUID,
				}
				return
			}
		}
	}()

	// Write the body to file
	_, err = io.Copy(dstFile, io.TeeReader(resp.Body, &p))
	if err != nil {
		ImageDownloadProcess = false
		return err
	}
	ImageDownloadProcess = false

	return nil
}

func GetImageDownloadProcess() bool {
	return ImageDownloadProcess
}
