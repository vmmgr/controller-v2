package db

import "fmt"

//Add
func AddDBNode(data Node) error {
	db := InitDB()
	defer db.Close()
	db.Create(&data)

	if err := db.Error; err != nil {
		db.Rollback()
		return fmt.Errorf("DB Error ")
	} else {
		return nil
	}
}

//Delete
func DeleteDBNode(data Node) error {
	db := InitDB()
	defer db.Close()
	db.Delete(&data)

	if err := db.Error; err != nil {
		db.Rollback()
		return fmt.Errorf("DB Error ")
	} else {
		return nil
	}
}

//Update
func UpdateDBNode(data Node) error {
	db := InitDB()
	defer db.Close()
	db.Model(&data).Updates(Node{HostName: data.HostName, IP: data.IP, Path: data.Path,
		OnlyAdmin: data.OnlyAdmin, MaxCPU: data.MaxCPU, MaxMem: data.MaxMem, Active: data.Active})

	if err := db.Error; err != nil {
		db.Rollback()
		return fmt.Errorf("DB Error ")
	} else {
		return nil
	}
}

//Get
func GetAllDBNode() []Node {
	db := InitDB()
	defer db.Close()

	var user []Node
	db.Find(&user)
	return user
}

func SearchDBNode(id int) Node {
	db := InitDB()
	defer db.Close()

	var result Node
	//search NodeName and NodeID
	if id != 0 { //初期値0であることが前提　確認の必要あり
		db.Where("id = ?", id).First(&result)
	}

	return result
}
