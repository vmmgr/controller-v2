package v0

import (
	"fmt"
	"github.com/pkg/sftp"
	"github.com/vmmgr/controller/pkg/api/core/tool/ssh"
	"github.com/vmmgr/controller/pkg/api/core/vm/cloudinit"
	"gopkg.in/yaml.v2"
	"log"
	"path/filepath"
)

type CloudInit struct {
	DirPath       string
	Auth          ssh.Auth
	MetaData      cloudinit.MetaData   `json:"meta"`
	UserData      cloudinit.UserData   `json:"user"`
	NetworkConfig cloudinit.NetworkCon `json:"network"`
}

func NewCloudInitHandler(input CloudInit) *CloudInit {
	return &CloudInit{
		DirPath:       input.DirPath,
		Auth:          input.Auth,
		MetaData:      input.MetaData,
		UserData:      input.UserData,
		NetworkConfig: input.NetworkConfig,
	}
}

// humstackのcloud-init部分を引用
// https://github.com/ophum/humstack/tree/master/pkg/utils/cloudinit

func (c *CloudInit) Generate() error {
	// SFTP
	conn, err := c.Auth.SSHClient()
	if err != nil {
		log.Println(err)
		return err
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Println(err)
		return err
	}
	defer client.Close()

	// make directory
	client.MkdirAll(c.DirPath)
	// MetaData
	file, err := client.Create(c.DirPath + "/meta-data")
	if err != nil {
		log.Println(err)
		return err
	}
	metaDataYAML, err := yaml.Marshal(c.MetaData)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = file.Write(metaDataYAML)
	if err != nil {
		log.Println(err)
		return err
	}

	// UserData
	file, err = client.Create(c.DirPath + "/user-data")
	if err != nil {
		log.Println(err)
		return err
	}
	userDataYAML, err := yaml.Marshal(c.UserData)
	if err != nil {
		return err
	}
	userDataYAML = []byte(fmt.Sprintf("#cloud-config\n%s", userDataYAML))
	_, err = file.Write(userDataYAML)
	if err != nil {
		return err
	}

	// NetworkConfig
	// NetworkConfig Version指定
	file, err = client.Create(c.DirPath + "/network-config")
	if err != nil {
		return err
	}

	c.NetworkConfig.Version = 1
	networkConfigYAML, err := yaml.Marshal(c.NetworkConfig)
	if err != nil {
		return err
	}
	_, err = file.Write(networkConfigYAML)
	if err != nil {
		return err
	}

	//ファイルがすでに存在する場合は削除する
	client.Remove(filepath.Join(c.DirPath, "cloudinit.img"))

	ah := c.Auth
	ah.SSHClientExecCmd("cloud-localds -N " + filepath.Join(c.DirPath, "network-config") + " " +
		filepath.Join(c.DirPath, "cloudinit.img") + " " +
		filepath.Join(c.DirPath, "user-data") + " " +
		filepath.Join(c.DirPath, "meta-data"))

	return nil
}
