package db

//Add
func AddDBNode(data Node) bool {
	db := InitDB()
	defer db.Close()
	db.Create(&data)

	if err := db.Error; err != nil {
		db.Rollback()
		return false
	} else {
		return true
	}
}

//Delete
func DeleteDBNode(data Node) bool {
	db := InitDB()
	defer db.Close()
	db.Delete(&data)

	if err := db.Error; err != nil {
		db.Rollback()
		return false
	} else {
		return true
	}
}

//Update
func UpdateDBNode(data Node) bool {
	db := InitDB()
	defer db.Close()
	db.Model(&data).Updates(Node{HostName: data.HostName, IP: data.IP, Path: data.Path,
		OnlyAdmin: data.OnlyAdmin, MaxCPU: data.MaxCPU, MaxMem: data.MaxMem, Active: data.Active})

	if err := db.Error; err != nil {
		db.Rollback()
		return false
	} else {
		return true
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

func SearchDBNode(data Node) Node {
	db := InitDB()
	defer db.Close()

	var result Node
	//search NodeName and NodeID
	if data.ID != 0 { //初期値0であることが前提　確認の必要あり
		db.Where("id = ?", data.ID).First(&result)
	}

	return result
}
