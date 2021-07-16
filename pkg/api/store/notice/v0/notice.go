package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/core/notice"
	"github.com/vmmgr/controller/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(notice *core.Notice) (*core.Notice, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return notice, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return notice, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Create(&notice).Error
	return notice, err
}

func Delete(notice *core.Notice) error {
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

	return db.Delete(notice).Error
}

func Update(base int, data core.Notice) error {
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

	if notice.UpdateAll == base {
		result = db.Model(&core.Notice{Model: gorm.Model{ID: data.ID}}).Updates(core.Notice{
			StartTime: data.StartTime,
			Important: data.Important,
			Fault:     data.Fault,
			Info:      data.Info,
			Title:     data.Title,
			Data:      data.Data,
		})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *core.Notice) notice.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return notice.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return notice.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var noticeStruct []core.Notice

	if base == notice.ID { //ID
		err = db.First(&noticeStruct, data.ID).Error
	} else if base == notice.Everyone { //Everyone
		err = db.Where("everyone = ?", data.Everyone).Find(&noticeStruct).Error
	} else if base == notice.Important { //Important
		err = db.Where("important = ?", data.Important).Find(&noticeStruct).Error
	} else if base == notice.Fault { //Fault
		err = db.Where("fault = ?", data.Fault).Find(&noticeStruct).Error
	} else if base == notice.Info { //Info
		err = db.Where("info = ?", data.Info).Find(&noticeStruct).Error
	} else {
		log.Println("base select error")
		return notice.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return notice.ResultDatabase{Notice: noticeStruct, Err: err}
}

func GetAll() notice.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return notice.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return notice.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var notices []core.Notice
	err = db.Find(&notices).Error
	return notice.ResultDatabase{Notice: notices, Err: err}
}
