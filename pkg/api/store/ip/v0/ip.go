package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/store"
	ip2 "github.com/vmmgr/controller/pkg/api/store/ip"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(ip *core.IP) (*core.IP, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return ip, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return ip, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Create(&ip).Error
	return ip, err
}

func Delete(ip *core.IP) error {
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

	return db.Delete(ip).Error
}

func Update(base int, data core.IP) error {
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
	if ip2.UpdateVMID == base {
		result = db.Model(&core.IP{Model: gorm.Model{ID: data.ID}}).Updates(&core.IP{VMID: data.VMID})
	} else if ip2.UpdateReserved == base {
		result = db.Model(&core.IP{Model: gorm.Model{ID: data.ID}}).Updates(&core.IP{Reserved: data.Reserved})
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
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return ips, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	if base == ip2.GetID { //ID
		err = db.First(&ips, data.ID).Error
	} else if base == ip2.GetUnused { //Zone内の全VM検索
		err = db.Where("reserved = ?", false).Find(&ips).Error
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
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return ips, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Find(&ips).Error
	return ips, err
}
