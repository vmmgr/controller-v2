package v0

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/pkg/api/core/node/storage"
	"github.com/vmmgr/controller/pkg/api/store"
	"log"
	"time"
)

func Create(storage *storage.Storage) (*storage.Storage, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return storage, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&storage).Error
	return storage, err
}

func Delete(storage *storage.Storage) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(storage).Error
}

func Update(base int, data storage.Storage) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB
	if storage.UpdateAll == base {
		result = db.Model(&storage.Storage{Model: gorm.Model{ID: data.ID}}).Update(storage.Storage{
			NodeID: data.NodeID, AdminOnly: data.AdminOnly, Type: data.Type, Path: data.Path,
			MaxCapacity: data.MaxCapacity, Comment: data.Comment})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *storage.Storage) storage.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return storage.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var storageStruct []storage.Storage

	if base == storage.ID { //ID
		err = db.First(&storageStruct, data.ID).Error
	} else if base == storage.NodeID { //Node内の全Storage検索
		err = db.Where("node_id = ?", data.NodeID).Find(&storageStruct).Error
	} else if base == storage.AdminOnly { //Node StorageID とVMIDから検索
		err = db.Where("admin_only = ?", data.AdminOnly).Find(&storageStruct).Error
	} else if base == storage.Name { //Name
		err = db.Where("lock = ?", data.Name).Find(&storageStruct).Error
	} else {
		log.Println("base select error")
		return storage.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return storage.ResultDatabase{Storage: storageStruct, Err: err}
}

func GetAll() storage.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return storage.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var storages []storage.Storage
	err = db.Find(&storages).Error
	return storage.ResultDatabase{Storage: storages, Err: err}
}
