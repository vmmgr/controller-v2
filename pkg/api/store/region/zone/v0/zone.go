package v0

import (
	"fmt"
	"github.com/jinzhu/gorm"
	zone "github.com/vmmgr/controller/pkg/api/core/region/zone"
	"github.com/vmmgr/controller/pkg/api/store"
	"log"
	"time"
)

func Create(zone *zone.Zone) (*zone.Zone, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return zone, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&zone).Error
	return zone, err
}

func Delete(zone *zone.Zone) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(zone).Error
}

func Update(base int, data zone.Zone) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB
	if zone.UpdateAll == base {
		result = db.Model(&zone.Zone{Model: gorm.Model{ID: data.ID}}).Update(zone.Zone{
			Name: data.Name, Comment: data.Comment, Lock: data.Lock})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *zone.Zone) zone.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return zone.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var zoneStruct []zone.Zone

	if base == zone.ID { //ID
		err = db.First(&zoneStruct, data.ID).Error
	} else if base == zone.RegionID {
		err = db.Where("region_id = ?", data.RegionID).Find(&zoneStruct).Error
	} else if base == zone.Name {
		err = db.Where("name = ?", data.Name).Find(&zoneStruct).Error
	} else {
		log.Println("base select error")
		return zone.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return zone.ResultDatabase{Zone: zoneStruct, Err: err}
}

func GetAll() zone.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return zone.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var zones []zone.Zone
	err = db.Find(&zones).Error
	return zone.ResultDatabase{Zone: zones, Err: err}
}
