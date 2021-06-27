package ssh

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"strconv"
)

type Auth struct {
	IP   string
	Port uint
	User string
	Pass string
}

func (h *Auth) SSHClientExecCmd(command string) (string, error) {
	config := &ssh.ClientConfig{
		User:            h.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//Auth:            []ssh.AuthMethod{ssh.Password(h.Pass)},
		Auth: []ssh.AuthMethod{PublicKeyFile("/home/yonedayuto/.ssh/id_rsa")},
	}

	conn, err := ssh.Dial("tcp", h.IP+":"+strconv.Itoa(int(h.Port)), config)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err = session.Run(command); err != nil {
		log.Fatal("Failed to run: " + err.Error())
		return "", err
	}
	log.Println(command + ":" + b.String())

	return b.String(), nil
}

func (h *Auth) SSHClient() (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User:            h.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//Auth:            []ssh.AuthMethod{ssh.Password(h.Pass)},
		Auth: []ssh.AuthMethod{PublicKeyFile("/home/yonedayuto/.ssh/id_rsa")},
	}

	conn, err := ssh.Dial("tcp", h.IP+":"+strconv.Itoa(int(h.Port)), config)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		log.Println(err)
		return nil
	}

	return ssh.PublicKeys(key)
}
