package v0

import (
	"github.com/vmmgr/controller/pkg/api/core/tool/ssh"
	"github.com/vmmgr/controller/pkg/api/core/vm/storage"
	"log"
	"os/exec"
	"strconv"
)

func (h *StorageHandler) convertImage(d storage.Convert) error {
	sh := ssh.Auth{
		IP:   h.Auth.IP,
		Port: h.Auth.Port,
		User: h.Auth.User,
		Pass: h.Auth.Pass,
	}
	//qemu-img convert -f raw -O qcow2 image.img image.qcow2
	command := "qemu-img" + " convert" + " -f " + d.SrcType + " -O " + d.DstType + " " + d.SrcFile + " " + d.DstFile
	result, err := sh.SSHClientExecCmd(command)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(result)
	//out, err := exec.Command("qemu-img", "convert", "-f", d.SrcType, "-O", d.DstType, d.SrcFile, d.DstFile).Output()
	//if err != nil {
	//	return err
	//}
	//log.Println(string(out))
	return nil
}

func (h *StorageHandler) generateImage(fileType, filePath string, fileSize uint) (string, error) {
	sh := ssh.Auth{
		IP:   h.Auth.IP,
		Port: h.Auth.Port,
		User: h.Auth.User,
		Pass: h.Auth.Pass,
	}
	//qemu-img create -f qcow2 file.qcow2 100M
	command := "qemu-img " + "create " + "-f " + fileType + " " + filePath + " " + strconv.Itoa(int(fileSize))
	result, err := sh.SSHClientExecCmd(command)
	if err != nil {
		log.Println(err)
		return "", err
	}

	log.Println(result)

	return result, nil

	//size := strconv.Itoa(int(fileSize)) + "M"
	//out, err := exec.Command("qemu-img", "create", "-f", fileType, filePath, size).Output()
	//if err != nil {
	//	return "", err
	//}
	//log.Println(string(out))
	//return string(out), nil
}

func infoImage(filePath string) (string, error) {
	//qemu-img info file.qcow2
	out, err := exec.Command("qemu-img", "info", filePath).Output()
	if err != nil {
		return "", err
	}
	log.Println(string(out))
	return string(out), nil
}

func capacityExpansion(filePath string, size uint) (string, error) {
	sizeString := strconv.Itoa(int(size)) + "M"
	out, err := exec.Command("qemu-img", "resize", filePath, sizeString).Output()
	if err != nil {
		return "", err
	}
	log.Println(string(out))
	return string(out), nil
}
