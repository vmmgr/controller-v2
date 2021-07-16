package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core"
	node "github.com/vmmgr/controller/pkg/api/core/node"
	"github.com/vmmgr/controller/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(node *core.Node) (*core.Node, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return node, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return node, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Create(&node).Error
	return node, err
}

func Delete(node *core.Node) error {
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

	return db.Delete(node).Error
}

func Update(base int, data core.Node) error {
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
	if node.UpdateAll == base {
		result = db.Model(&core.Node{Model: gorm.Model{ID: data.ID}}).Updates(core.Node{
			ZoneID:    data.ZoneID,
			GroupID:   data.GroupID,
			AdminOnly: data.AdminOnly,
			Name:      data.Name,
			IP:        data.IP,
			Port:      data.Port,
			WsPort:    data.WsPort,
			ManageNet: data.ManageNet,
			Comment:   data.Comment,
		})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *core.Node) node.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return node.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return node.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var nodeStruct []core.Node

	if base == node.ID { //ID
		err = db.Preload("Storage").
			Preload("NIC").
			First(&nodeStruct, data.ID).Error
	} else if base == node.ZoneID { //Zone内の全VM検索
		err = db.Where("zone_id = ?", data.ZoneID).Find(&nodeStruct).Error
	} else if base == node.GroupID { //GroupID
		err = db.Where("group_id = ?", data.GroupID).Find(&nodeStruct).Error
	} else if base == node.AdminOnly { //UUID
		err = db.Where("admin_only = ?", data.AdminOnly).Find(&nodeStruct).Error
	} else if base == node.Name { //VNCPortの空きポートの検索
		err = db.Where("name = ? ", data.Name).Find(&nodeStruct).Error
	} else {
		log.Println("base select error")
		return node.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return node.ResultDatabase{Node: nodeStruct, Err: err}
}

func GetAll() node.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return node.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return node.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var nodes []core.Node
	err = db.Preload("Storage").
		Preload("NIC").
		Find(&nodes).Error
	return node.ResultDatabase{Node: nodes, Err: err}
}
