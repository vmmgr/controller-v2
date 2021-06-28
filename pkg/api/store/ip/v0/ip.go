package v0

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/store"
	ip2 "github.com/vmmgr/controller/pkg/api/store/ip"
	"log"
	"time"
)

func Create(nic *core.IP) (*core.IP, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return nic, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&nic).Error
	return nic, err
}

func Delete(nic *core.IP) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(nic).Error
}

func Update(base int, data core.IP) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB
	if ip2.UpdateVMID == base {
		result = db.Model(&core.IP{Model: gorm.Model{ID: data.ID}}).Update(&core.IP{VMID: data.VMID})
	} else if ip2.UpdateReserved == base {
		result = db.Model(&core.IP{Model: gorm.Model{ID: data.ID}}).Update(&core.IP{Reserved: data.Reserved})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *core.IP) ([]core.IP, error) {
	var ips []core.IP

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return ips, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	if base == ip2.GetID { //ID
		err = db.First(&ips, data.ID).Error
	} else if base == ip2.GetUnused { //Zone内の全VM検索
		err = db.Where("vm_id = ? AND reserved = ?", 0, false).Find(&ips).Error
	} else {
		log.Println("base select error")
		return ips, fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return ips, err
}

func GetAll() ([]core.IP, error) {
	var ips []core.IP

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return ips, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Find(&ips).Error
	return ips, err
}
