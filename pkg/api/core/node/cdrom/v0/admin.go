package v0

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/pkg/api/core"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/common"
	cdrom "github.com/vmmgr/controller/pkg/api/core/node/cdrom"
	"github.com/vmmgr/controller/pkg/api/core/node/storage"
	"github.com/vmmgr/controller/pkg/api/core/tool/gen"
	image "github.com/vmmgr/controller/pkg/api/core/tool/image"
	"github.com/vmmgr/controller/pkg/api/core/tool/ssh"
	dbNIC "github.com/vmmgr/controller/pkg/api/store/node/nic/v0"
	dbStorage "github.com/vmmgr/controller/pkg/api/store/node/storage/v0"
	"log"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	var input cdrom.Post

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// 入力チェック
	err = inputCheck(input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if image.GetImageDownloadProcess() {
		c.JSON(http.StatusServiceUnavailable, common.Error{Error: "Another process is downloading"})
		return
	}

	resultStorage := dbStorage.Get(storage.ID, &core.Storage{Model: gorm.Model{ID: uint(id)}})
	if resultStorage.Err != nil {
		log.Println(resultStorage.Err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultStorage.Err.Error()})
		return
	}

	// ISO,Floppyイメージ格納ストレージであるか確認
	if *resultStorage.Storage[0].VMImage {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "This storage id is only VM image..."})
		return
	}

	// SSH接続
	conn, err := ssh.ConnectSSH(resultStorage.Storage[0].NodeID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	uuid := gen.GenerateUUID()

	imageHandler := image.ImageHandler{
		SSHClient: conn,
		DstPath:   resultStorage.Storage[0].Path + "/" + input.Name + ".iso",
		UUID:      uuid,
	}

	go imageHandler.ImageDownload(input.URL)

	c.JSON(http.StatusOK, common.Result{UUID: uuid})
}

func DeleteAdmin(c *gin.Context) {
	var input cdrom.Delete

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	resultStorage := dbStorage.Get(storage.ID, &core.Storage{Model: gorm.Model{ID: uint(id)}})
	if resultStorage.Err != nil {
		log.Println(resultStorage.Err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultStorage.Err.Error()})
		return
	}

	// SSH接続
	conn, err := ssh.ConnectSSH(resultStorage.Storage[0].NodeID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultStorage.Err.Error()})
		return
	}
	defer session.Close()
	//Check whoami
	var b bytes.Buffer
	session.Stdout = &b
	remoteCommand := "rm " + resultStorage.Storage[0].Path + "/" + input.Name + ".iso"
	if err = session.Run(remoteCommand); err != nil {
		log.Println("Failed to run: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultStorage.Err.Error()})
		return
	}
	log.Println(remoteCommand + ":" + b.String())

	c.JSON(http.StatusOK, common.Result{})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbNIC.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, cdrom.Result{NIC: result.NIC})
	}
}
