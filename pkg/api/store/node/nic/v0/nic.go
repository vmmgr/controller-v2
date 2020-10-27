package v0

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/pkg/api/core/node/nic"
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
	if nic.UpdateAll == base {
		result = db.Model(&nic.NIC{Model: gorm.Model{ID: data.ID}}).Update(nic.NIC{
			NodeID: data.NodeID, GroupID: data.GroupID, AdminOnly: data.AdminOnly, Name: data.Name, Enable: data.Enable,
			Virtual: data.Virtual, Type: data.Type, Vlan: data.Vlan, Speed: data.Speed, MAC: data.MAC,
			Comment: data.Comment})
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
	} else if base == nic.NodeID { //Zone内の全VM検索
		err = db.Where("zone_id = ?", data.NodeID).Find(&nicStruct).Error
	} else if base == nic.GroupID { //GroupID
		err = db.Where("group_id = ?", data.GroupID).Find(&nicStruct).Error
	} else if base == nic.Name { //UUID
		err = db.Where("name = ?", data.Name).Find(&nicStruct).Error
	} else if base == nic.AdminOnly {
		err = db.Where("admin_only = ? ", data.AdminOnly).Find(&nicStruct).Error
	} else if base == nic.Enable {
		err = db.Where("enable = ?", data.Enable).Find(&nicStruct).Error
	} else if base == nic.Virtual {
		err = db.Where("virtual = ?", data.Virtual).Find(&nicStruct).Error
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
