package v0

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/store"
	"log"
	"time"
)

func Create(node *core.TemplatePlan) (*core.TemplatePlan, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return node, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&node).Error
	return node, err
}

func Delete(node *core.TemplatePlan) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(node).Error
}

func Update(data core.TemplatePlan) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB
	result = db.Model(&core.TemplatePlan{Model: gorm.Model{ID: data.ID}}).Update(data)

	return result.Error
}

func Get(data core.TemplatePlan) (*[]core.TemplatePlan, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return nil, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var plans []core.TemplatePlan

	err = db.First(&plans, data.ID).Error
	return &plans, nil
}

func GetAll() (*[]core.TemplatePlan, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return nil, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var plans []core.TemplatePlan
	err = db.Find(&plans).Error
	return &plans, nil
}
