package v2

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/store/ip"
	dbIP "github.com/vmmgr/controller/pkg/api/store/ip/v0"
	dbVM "github.com/vmmgr/controller/pkg/api/store/vm/v0"
	"log"
	"strconv"
)

func (h *VMHandler) CreateVM() error {

	defer h.Conn.Close()

	log.Println("Create VM Process")

	// VNC Portが0の場合、自動生成を行う
	if h.VM.VNCPort == 0 {
		vnc, err := h.generateVNC()
		if err != nil {
			return err
		}
		h.VM.VNCPort = uint(vnc.VNCPort)
		h.VM.WebSocketPort = uint(vnc.WebSocketPort)
	}

	domCfg, err := h.xmlGenerate()
	if err != nil {
		log.Println(err)
		return err
	}

	xml, err := domCfg.Marshal()
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println(xml)
	log.Println("XML Apply start")
	dom, err := h.Conn.DomainDefineXML(xml)
	if err != nil {
		return err
	}
	log.Println("dom Create start")

	if h.GroupID != 0 {
		uuid, _ := dom.GetUUIDString()
		vm, err := dbVM.Create(&core.VM{
			NodeID:        h.Node.ID,
			GroupID:       h.GroupID,
			Name:          strconv.Itoa(int(h.GroupID)) + "-" + h.VM.Name,
			UUID:          uuid,
			VNCPort:       h.VM.VNCPort,
			WebSocketPort: h.VM.WebSocketPort,
		})
		if err != nil {
			log.Println(err)
		}
		err = dbIP.Update(ip.UpdateVMID, core.IP{Model: gorm.Model{ID: h.IPID}, VMID: vm.ID})
		if err != nil {
			log.Println(err)
		}
	}

	err = dom.Create()
	if err != nil {
		// node側でエラーを表示
		log.Println(err)
		return err
	}
	return nil
}
