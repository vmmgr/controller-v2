package v0

import (
	"fmt"
	"github.com/jinzhu/gorm"
	node "github.com/vmmgr/controller/pkg/api/core/node"
	"github.com/vmmgr/controller/pkg/api/store"
	"log"
	"time"
)

func Create(node *node.Node) (*node.Node, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return node, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&node).Error
	return node, err
}

func Delete(node *node.Node) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(node).Error
}

func Update(base int, data node.Node) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB
	if node.UpdateAll == base {
		result = db.Model(&node.Node{Model: gorm.Model{ID: data.ID}}).Update(node.Node{
			ZoneID: data.ZoneID, GroupID: data.GroupID, AdminOnly: data.AdminOnly, Name: data.Name, IP: data.IP,
			Port: data.Port, WsPort: data.WsPort, ManageNet: data.ManageNet, Comment: data.Comment})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *node.Node) node.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return node.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var nodeStruct []node.Node

	if base == node.ID { //ID
		err = db.First(&nodeStruct, data.ID).Error
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
	defer db.Close()

	var nodes []node.Node
	err = db.Find(&nodes).Error
	return node.ResultDatabase{Node: nodes, Err: err}
}
