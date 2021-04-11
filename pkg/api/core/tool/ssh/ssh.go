package ssh

import (
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/core/node"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	"golang.org/x/crypto/ssh"
	"log"
	"strconv"
)

func ConnectSSHNodeID(nodeID uint) (*ssh.Client, error) {
	resultNode := dbNode.Get(node.ID, &core.Node{Model: gorm.Model{ID: nodeID}})
	if resultNode.Err != nil {
		log.Println(resultNode.Err)
		return nil, resultNode.Err
	}

	var config *ssh.ClientConfig

	ip := resultNode.Node[0].IP
	port := strconv.Itoa(int(resultNode.Node[0].Port))
	user := resultNode.Node[0].UserName
	pass := resultNode.Node[0].Password

	// 鍵認証
	if *resultNode.Node[0].UseKey {
		key, err := ssh.ParsePrivateKey([]byte(resultNode.Node[0].PublicKey))
		if err != nil {
			log.Println(err)
			return nil, err
		}

		config = &ssh.ClientConfig{
			User:            user,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // https://github.com/golang/go/issues/19767
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(key),
			},
		}
	} else {
		//Passフレーズ認証
		config = &ssh.ClientConfig{
			User:            user,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // https://github.com/golang/go/issues/19767
			Auth: []ssh.AuthMethod{
				ssh.Password(pass),
			},
		}

	}

	conn, err := ssh.Dial("tcp", ip+":"+port, config)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}

func ConnectSSH(node core.Node) (*ssh.Client, error) {
	var config *ssh.ClientConfig

	ip := node.IP
	port := strconv.Itoa(int(node.Port))
	user := node.UserName
	pass := node.Password

	// 鍵認証
	if *node.UseKey {
		key, err := ssh.ParsePrivateKey([]byte(node.PublicKey))
		if err != nil {
			log.Println(err)
			return nil, err
		}

		config = &ssh.ClientConfig{
			User:            user,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // https://github.com/golang/go/issues/19767
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(key),
			},
		}
	} else {
		//Passフレーズ認証
		config = &ssh.ClientConfig{
			User:            user,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // https://github.com/golang/go/issues/19767
			Auth: []ssh.AuthMethod{
				ssh.Password(pass),
			},
		}

	}

	conn, err := ssh.Dial("tcp", ip+":"+port, config)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}
