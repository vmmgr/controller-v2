package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/core/group"
	"github.com/vmmgr/controller/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(g *core.Group) (*core.Group, error) {
	result := Get(group.Org, &core.Group{Org: g.Org})
	if result.Err != nil {
		return &core.Group{}, result.Err
	}
	if len(result.Group) != 0 {
		log.Println("error: this Org Name is already registered: " + g.Org)
		return &core.Group{}, fmt.Errorf("error: this org name is already registered")
	}

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

	err = db.Create(&g).Error
	return g, err
}

func Delete(group *core.Group) error {
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

	return db.Delete(group).Error
}

func Update(base int, g core.Group) error {
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

	if group.UpdateOrg == base {
		result = db.Model(&core.Group{Model: gorm.Model{ID: g.ID}}).Updates(core.Group{Org: g.Org})
	} else if group.UpdateStatus == base {
		result = db.Model(&core.Group{Model: gorm.Model{ID: g.ID}}).Updates(core.Group{Status: g.Status})
	} else if group.UpdateInfo == base {
		result = db.Model(&core.Group{Model: gorm.Model{ID: g.ID}}).Updates(core.Group{Org: g.Org})
	} else if group.UpdateAll == base {
		result = db.Model(&core.Group{Model: gorm.Model{ID: g.ID}}).Updates(core.Group{
			Org:     g.Org,
			Status:  g.Status,
			Comment: g.Comment,
			Lock:    g.Lock,
		})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *core.Group) group.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return group.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return group.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var groupStruct []core.Group

	if base == group.ID { //ID
		err = db.Preload("VMs").Preload("VMs.Node").First(&groupStruct, data.ID).Error
	} else if base == group.Org { //Org
		err = db.Where("org = ?", data.Org).Find(&groupStruct).Error
	} else {
		log.Println("base select error")
		return group.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return group.ResultDatabase{Group: groupStruct, Err: err}
}

func GetAll() group.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return group.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return group.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var groups []core.Group
	err = db.Find(&groups).Error
	return group.ResultDatabase{Group: groups, Err: err}
}
