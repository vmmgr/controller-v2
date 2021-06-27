package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/libvirt/libvirt-go"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"github.com/vmmgr/controller/pkg/api/core/tool/ssh"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"github.com/vmmgr/controller/pkg/api/core/vm/storage"
	"github.com/vmmgr/node/pkg/api/meta/json"
	"log"
	"net/http"
	"strconv"
)

type StorageHandler struct {
	UUID      string
	Conn      *libvirt.Connect
	Input     storage.Storage
	VM        vm.VirtualMachine
	Address   *vm.Address
	Auth      *storage.Auth
	SrcImaCon core.ImaCon
	DstAuth   storage.Auth
	SrcPath   string
	DstPath   string
}

func NewStorageHandler(handler StorageHandler) *StorageHandler {
	return &StorageHandler{
		UUID:      handler.UUID,
		Conn:      handler.Conn,
		Input:     handler.Input,
		VM:        handler.VM,
		Address:   handler.Address,
		Auth:      handler.Auth,
		SrcImaCon: handler.SrcImaCon,
		DstAuth:   handler.DstAuth,
		SrcPath:   handler.SrcPath,
		DstPath:   handler.DstPath,
	}
}

func (h *StorageHandler) AddFromImage() error {
	// ImaConからイメージ取得(時間がかかるので、go funcにて処理)
	log.Println("From: " + h.SrcPath)
	log.Println("To: " + h.DstPath)

	sh := ssh.Auth{
		IP:   h.SrcImaCon.IP,
		Port: h.SrcImaCon.Port,
		User: h.SrcImaCon.User,
		Pass: h.SrcImaCon.Pass,
	}
	//qemu-img create -f qcow2 file.qcow2 100M
	command := h.SrcImaCon.AppPath + " copy --uuid " + h.UUID + " --url " + config.Conf.Controller.Admin.URL +
		" --src " + h.SrcPath + " --dst " + h.DstPath + " --addr " + h.DstAuth.IP + ":" +
		strconv.Itoa(int(h.DstAuth.Port)) + " --user " + h.DstAuth.User + " --config " + h.SrcImaCon.ConfigPath
	log.Println(command)
	result, err := sh.SSHClientExecCmd(command)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(result)

	log.Println("Done: Image Create")

	return nil
}

func (h *StorageHandler) Add(input storage.Storage) error {

	path := ""
	path = input.DstSrv.Path + "/" + input.VMName

	//for _, tmpConf := range config.Conf.Storage {
	//	if tmpConf.Type == input.PathType {
	//		if input.VMName == "" {
	//			path = tmpConf.Path + "/" + input.Path
	//		} else {
	//			if err := os.Mkdir(tmpConf.Path+"/"+input.VMName, 0775); err != nil {
	//				log.Println(err)
	//				return err
	//			}
	//			path = tmpConf.Path + "/" + input.VMName + "/" + input.Path
	//		}
	//	}
	//}

	log.Println(path)
	// Pathが見つからない場合
	if path == "" {
		return fmt.Errorf("Error: Not found... ")
	}

	if FileExistsCheck(path) {
		return fmt.Errorf("Error: file already exists... ")
	}

	// イメージの作成
	out, err := h.generateImage(storage.GetExtensionName(input.Type), input.Path, input.Capacity)
	if err != nil {
		log.Println(out)
		return err
	}
	log.Println("Done: Image Create")
	//node.SendServer(h.Input.Info, 1, 100, "Done: Image Create", nil)
	return nil
}

//func (h *StorageHandler) Add(c *gin.Context) {
//	var input storage.Storage
//
//	err := c.BindJSON(&input)
//	if err != nil {
//		json.ResponseError(c, http.StatusBadRequest, err)
//		return
//	}
//
//	path := ""
//
//	for _, tmpConf := range config.Conf.Storage {
//		if tmpConf.Type == input.PathType {
//			if input.VMName == "" {
//				path = tmpConf.Path + "/" + input.Path
//			} else {
//				if err := os.Mkdir(tmpConf.Path+"/"+input.VMName, 0775); err != nil {
//					log.Println(err)
//					json.ResponseError(c, http.StatusInternalServerError, err)
//					return
//				}
//				path = tmpConf.Path + "/" + input.VMName + "/" + input.Path
//			}
//		}
//	}
//
//	log.Println(path)
//	// Pathが見つからない場合
//	if path == "" {
//		json.ResponseError(c, http.StatusNotFound, fmt.Errorf("Error: Not found... "))
//		return
//	}
//
//	if FileExistsCheck(path) {
//		json.ResponseError(c, http.StatusNotFound, fmt.Errorf("Error: file already exists... "))
//		return
//	}
//
//	var out string
//
//	// イメージの作成
//	if input.Mode == 0 {
//		out, err = generateImage(storage.GetExtensionName(input.Type), input.Path, input.Capacity)
//		if err != nil {
//			json.ResponseError(c, http.StatusNotFound, err)
//			return
//		} else {
//			json.ResponseOK(c, out)
//		}
//	} else if input.Mode == 1 {
//		// ImaConからイメージ取得(時間がかかるので、go funcにて処理)
//		go func() {
//			log.Println("From: " + input.FromImaCon.Path)
//			log.Println("To: " + path)
//
//			//メソッドに各種情報の追加
//			h.Auth = &storage.Auth{
//				IP: input.FromImaCon.IP, User: config.Conf.ImaCon.User, Pass: config.Conf.ImaCon.Pass,
//			}
//			h.SrcPath = input.FromImaCon.Path
//			h.DstPath = path
//			h.Input = input
//
//			err := h.sftpRemoteToLocal()
//			log.Println(err)
//		}()
//
//		json.ResponseOK(c, out)
//	}
//}

func (h *StorageHandler) ConvertImage(c *gin.Context) {
	var input storage.Convert

	err := c.BindJSON(&input)
	if err != nil {
		json.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	// sourceファイルの確認
	if !FileExistsCheck(input.SrcFile) {
		json.ResponseError(c, http.StatusNotFound, fmt.Errorf("Error: file no exists... "))
		return
	}

	// Destinationファイルの確認
	if FileExistsCheck(input.DstFile) {
		json.ResponseError(c, http.StatusInternalServerError, fmt.Errorf("Error: file already exists... "))
		return
	}

	if err := h.convertImage(input); err != nil {
		json.ResponseError(c, http.StatusInternalServerError, err)
	} else {
		json.ResponseOK(c, nil)
	}
}

func (h *StorageHandler) InfoImage(c *gin.Context) {
	var input storage.Convert

	err := c.BindJSON(&input)
	if err != nil {
		json.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	// sourceファイルの確認
	if !FileExistsCheck(input.SrcFile) {
		json.ResponseError(c, http.StatusNotFound, fmt.Errorf("Error: file no exists... "))
		return
	}

	if data, err := infoImage(input.SrcFile); err != nil {
		json.ResponseError(c, http.StatusInternalServerError, err)
	} else {
		json.ResponseOK(c, data)
	}
}
