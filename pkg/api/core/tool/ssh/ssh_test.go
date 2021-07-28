package ssh

import (
	"bytes"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/core/node"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
	"strconv"
	"testing"
)

var nodeID uint = 2

func TestConfigApply(t *testing.T) {
	if err := config.GetConfig("../../../../../cmd/backend/con.json"); err != nil {
		t.Fatal(err)
	}
}

func TestSSHClient(t *testing.T) {
	resultNode := dbNode.Get(node.ID, &core.Node{Model: gorm.Model{ID: nodeID}})
	if resultNode.Err != nil {
		t.Fatal(resultNode.Err)
	}

	//t.Log(resultNode.Node[0])

	var config *ssh.ClientConfig

	ip := resultNode.Node[0].IP
	port := strconv.Itoa(int(resultNode.Node[0].Port))
	user := resultNode.Node[0].User
	pass := resultNode.Node[0].Pass

	// 鍵認証
	if *resultNode.Node[0].UseKey {
		t.Log("UseKey")
		key, err := ssh.ParsePrivateKey([]byte(resultNode.Node[0].PublicKey))
		if err != nil {
			panic(err)
		}

		config = &ssh.ClientConfig{
			User:            user,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // https://github.com/golang/go/issues/19767
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(key),
			},
		}
	} else {
		t.Log("UsePassword")
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
		t.Fatal(err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	//Check whoami
	var b bytes.Buffer
	session.Stdout = &b
	remoteCommand := "whoami"
	if err = session.Run(remoteCommand); err != nil {
		t.Fatal("Failed to run: " + err.Error())
	}
	t.Log(remoteCommand + ":" + b.String())
}
