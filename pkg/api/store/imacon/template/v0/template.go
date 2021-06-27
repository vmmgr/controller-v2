package v0

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/store"
	"log"
	"time"
)

func Create(node *core.Template) (*core.Template, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return node, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&node).Error
	return node, err
}

func Delete(node *core.Template) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(node).Error
}

func Update(data core.Template) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB
	result = db.Model(&core.Template{Model: gorm.Model{ID: data.ID}}).Update(data)

	return result.Error
}

func Get(data core.Template) ([]core.Template, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return nil, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var templates []core.Template

	err = db.First(&templates, data.ID).Error
	return templates, nil
}

func GetAll() ([]core.Template, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return nil, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var templates []core.Template
	err = db.Preload("Image").
		Preload("Image.ImaCon").
		Preload("TemplatePlan").
		Find(&templates).Error
	return templates, nil
}
