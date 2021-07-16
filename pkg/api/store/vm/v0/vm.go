package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"github.com/vmmgr/controller/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(vm *core.VM) (*core.VM, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return nil, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return nil, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Create(&vm).Error
	return vm, err
}

func Delete(vm *core.VM) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	return db.Delete(vm).Error
}

func Update(base int, data core.VM) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	var result *gorm.DB
	if vm.UpdateName == base {
		result = db.Model(&core.VM{Model: gorm.Model{ID: data.ID}}).Updates(core.VM{Name: data.Name})
	} else if vm.UpdateNode == base {
		result = db.Model(&core.VM{Model: gorm.Model{ID: data.ID}}).Updates(core.VM{NodeID: data.NodeID})
	} else if vm.UpdateGroup == base {
		result = db.Model(&core.VM{Model: gorm.Model{ID: data.ID}}).Updates(core.VM{GroupID: data.GroupID})
	} else if vm.UpdateInfo == base {
		result = db.Model(&core.VM{Model: gorm.Model{ID: data.ID}}).Updates(core.VM{
			Name:    data.Name,
			UUID:    data.UUID,
			VNCPort: data.VNCPort,
		})
	} else if vm.UpdateAll == base {
		result = db.Model(&core.VM{Model: gorm.Model{ID: data.ID}}).Updates(core.VM{
			NodeID: data.NodeID, GroupID: data.GroupID, Name: data.Name, UUID: data.UUID, VNCPort: data.VNCPort, Lock: data.Lock,
		})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *core.VM) vm.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return vm.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return vm.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var vmStruct []core.VM

	if base == vm.ID { //ID
		err = db.First(&vmStruct, data.ID).Error
	} else if base == vm.NodeID { //Node内の全VM検索
		err = db.Where("node_id = ?", data.NodeID).Find(&vmStruct).Error
	} else if base == vm.GroupID { //GroupID
		err = db.Where("group_id = ?", data.GroupID).Find(&vmStruct).Error
	} else if base == vm.UUID { //UUID
		err = db.Where("uuid = ?", data.UUID).Find(&vmStruct).Error
	} else if base == vm.NodeAndVNCPort { //VNCPortの空きポートの検索
		err = db.Where("node_id = ? AND vnc_port = ?", data.NodeID, data.VNCPort).Find(&vmStruct).Error
	} else if base == vm.Lock { //VM Lock
		err = db.Where("lock = ?", data.Lock).Find(&vmStruct).Error
	} else {
		log.Println("base select error")
		return vm.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return vm.ResultDatabase{VMs: vmStruct, Err: err}
}

func GetAll() vm.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return vm.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return vm.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var vms []core.VM
	err = db.Find(&vms).Error
	return vm.ResultDatabase{VMs: vms, Err: err}
}
