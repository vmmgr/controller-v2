package v0

//
//func Create(storage *core.Storage) (*core.Storage, error) {
//	db, err := store.ConnectDB()
//	if err != nil {
//		log.Println("database connection error")
//		return storage, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
//	}
//	defer db.Close()
//
//	err = db.Create(&storage).Error
//	return storage, err
//}
//
//func Delete(storage *core.Storage) error {
//	db, err := store.ConnectDB()
//	if err != nil {
//		log.Println("database connection error")
//		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
//	}
//	defer db.Close()
//
//	return db.Delete(storage).Error
//}
//
//func Update(base int, data core.Storage) error {
//	db, err := store.ConnectDB()
//	if err != nil {
//		log.Println("database connection error")
//		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
//	}
//	defer db.Close()
//
//	var result *gorm.DB
//	if storage.UpdateName == base {
//		result = db.Model(&core.Storage{Model: gorm.Model{ID: data.ID}}).Update(core.Storage{Name: data.Name})
//	} else if storage.UpdateNodeS == base {
//		result = db.Model(&core.Storage{Model: gorm.Model{ID: data.ID}}).Update(core.Storage{NodeStorageID: data.NodeStorageID})
//	} else if storage.UpdateGroup == base {
//		result = db.Model(&core.Storage{Model: gorm.Model{ID: data.ID}}).Update(core.Storage{GroupID: data.GroupID})
//	} else if storage.UpdateAll == base {
//		result = db.Model(&core.Storage{Model: gorm.Model{ID: data.ID}}).Update(core.Storage{
//			VMID: data.VMID, NodeStorageID: data.NodeStorageID, GroupID: data.GroupID,
//			Name: data.Name, Type: data.Type, FileType: data.FileType, MaxCapacity: data.MaxCapacity,
//			ReadOnly: data.ReadOnly, Comment: data.Comment, Lock: data.Lock})
//	} else {
//		log.Println("base select error")
//		return fmt.Errorf("(%s)error: base select\n", time.Now())
//	}
//	return result.Error
//}
//
//func Get(base int, data *core.Storage) storage.ResultDatabase {
//	db, err := store.ConnectDB()
//	if err != nil {
//		log.Println("database connection error")
//		return storage.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
//	}
//	defer db.Close()
//
//	var storageStruct []core.Storage
//
//	if base == storage.ID { //ID
//		err = db.First(&storageStruct, data.ID).Error
//	} else if base == storage.NodeStorageID { //NodeStorage内の全Storage検索
//		err = db.Where("node_storage_id = ?", data.NodeStorageID).Find(&storageStruct).Error
//	} else if base == storage.GroupID { //GroupID無いの全Storage検索
//		err = db.Where("group_id = ?", data.GroupID).Find(&storageStruct).Error
//	} else if base == storage.NodeSAndVMID { //Node StorageID とVMIDから検索
//		err = db.Where("node_storage_id = ? AND vm_id = ?", data.NodeStorageID, data.VMID).Find(&storageStruct).Error
//	} else if base == storage.Lock { //VM Lock
//		err = db.Where("lock = ?", data.Lock).Find(&storageStruct).Error
//	} else {
//		log.Println("base select error")
//		return storage.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
//	}
//	return storage.ResultDatabase{Storage: storageStruct, Err: err}
//}
//
//func GetAll() storage.ResultDatabase {
//	db, err := store.ConnectDB()
//	if err != nil {
//		log.Println("database connection error")
//		return storage.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
//	}
//	defer db.Close()
//
//	var storages []core.Storage
//	err = db.Find(&storages).Error
//	return storage.ResultDatabase{Storage: storages, Err: err}
//}
