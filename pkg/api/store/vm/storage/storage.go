package v0

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/pkg/api/core/vm/storage"
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
	if storage.UpdateName == base {
		result = db.Model(&storage.Storage{Model: gorm.Model{ID: data.ID}}).Update(storage.Storage{Name: data.Name})
	} else if storage.UpdateNodeS == base {
		result = db.Model(&storage.Storage{Model: gorm.Model{ID: data.ID}}).Update(storage.Storage{NodeStorageID: data.NodeStorageID})
	} else if storage.UpdateGroup == base {
		result = db.Model(&storage.Storage{Model: gorm.Model{ID: data.ID}}).Update(storage.Storage{GroupID: data.GroupID})
	} else if storage.UpdateAll == base {
		result = db.Model(&storage.Storage{Model: gorm.Model{ID: data.ID}}).Update(storage.Storage{
			VMID: data.VMID, NodeStorageID: data.NodeStorageID, GroupID: data.GroupID,
			Name: data.Name, Type: data.Type, FileType: data.FileType, MaxCapacity: data.MaxCapacity,
			Path: data.Path, ReadOnly: data.ReadOnly, Comment: data.Comment, Lock: data.Lock})
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
	} else if base == storage.NodeStorageID { //NodeStorage内の全Storage検索
		err = db.Where("node_storage_id = ?", data.NodeStorageID).Find(&storageStruct).Error
	} else if base == storage.GroupID { //GroupID無いの全Storage検索
		err = db.Where("group_id = ?", data.GroupID).Find(&storageStruct).Error
	} else if base == storage.NodeSAndVMID { //Node StorageID とVMIDから検索
		err = db.Where("node_storage_id = ? AND vm_id = ?", data.NodeStorageID, data.VMID).Find(&storageStruct).Error
	} else if base == storage.Lock { //VM Lock
		err = db.Where("lock = ?", data.Lock).Find(&storageStruct).Error
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
