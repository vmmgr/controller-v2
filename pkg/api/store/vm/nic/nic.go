package v0

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/pkg/api/core/vm/nic"
	"github.com/vmmgr/controller/pkg/api/store"
	"log"
	"time"
)

func Create(nic *nic.NIC) (*nic.NIC, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return nic, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&nic).Error
	return nic, err
}

func Delete(nic *nic.NIC) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(nic).Error
}

func Update(base int, data nic.NIC) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB
	if nic.UpdateName == base {
		result = db.Model(&nic.NIC{Model: gorm.Model{ID: data.ID}}).Update(nic.NIC{Name: data.Name})
	} else if nic.UpdateNodeN == base {
		result = db.Model(&nic.NIC{Model: gorm.Model{ID: data.ID}}).Update(nic.NIC{NodeNICID: data.NodeNICID})
	} else if nic.UpdateGroup == base {
		result = db.Model(&nic.NIC{Model: gorm.Model{ID: data.ID}}).Update(nic.NIC{GroupID: data.GroupID})
	} else if nic.UpdateMac == base {
		result = db.Model(&nic.NIC{Model: gorm.Model{ID: data.ID}}).Update(nic.NIC{Mac: data.Mac})
	} else if nic.UpdateAll == base {
		result = db.Model(&nic.NIC{Model: gorm.Model{ID: data.ID}}).Update(nic.NIC{VMID: data.VMID,
			NodeNICID: data.NodeNICID, GroupID: data.GroupID, Name: data.Name, Type: data.Type, Driver: data.Driver,
			Mode: data.Mode, Mac: data.Mac, Vlan: data.Vlan, Comment: data.Comment, Lock: data.Lock})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *nic.NIC) nic.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return nic.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var nicStruct []nic.NIC

	if base == nic.ID { //ID
		err = db.First(&nicStruct, data.ID).Error
	} else if base == nic.NodeNICID { //NodeNICID内の全VM検索
		err = db.Where("node_nic_id = ?", data.NodeNICID).Find(&nicStruct).Error
	} else if base == nic.GroupID { //GroupID
		err = db.Where("group_id = ?", data.GroupID).Find(&nicStruct).Error
	} else if base == nic.Name { //UUID
		err = db.Where("name = ?", data.Name).Find(&nicStruct).Error
	} else if base == nic.Type {
		err = db.Where("type = ? ", data.Type).Find(&nicStruct).Error
	} else if base == nic.Vlan {
		err = db.Where("vlan = ?", data.Lock).Find(&nicStruct).Error
	} else {
		log.Println("base select error")
		return nic.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return nic.ResultDatabase{NIC: nicStruct, Err: err}
}

func GetAll() nic.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return nic.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var nics []nic.NIC
	err = db.Find(&nics).Error
	return nic.ResultDatabase{NIC: nics, Err: err}
}
