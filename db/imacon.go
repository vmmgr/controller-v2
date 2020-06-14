package db

//Add
func AddDBImaCon(data ImaCon) bool {
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
func DeleteDBImaCon(data ImaCon) bool {
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
func UpdateDBImaCon(data ImaCon) bool {
	db := InitDB()
	defer db.Close()
	db.Model(&data).Updates(ImaCon{HostName: data.HostName, IP: data.IP, Status: data.Status})

	if err := db.Error; err != nil {
		db.Rollback()
		return false
	} else {
		return true
	}
}

//Get
func GetAllDBImaCon() []ImaCon {
	db := InitDB()
	defer db.Close()

	var user []ImaCon
	db.Find(&user)
	return user
}

func SearchDBImaCon(data ImaCon) ImaCon {
	db := InitDB()
	defer db.Close()

	var result ImaCon
	//search ImaConName and ImaConID
	if data.ID != 0 { //初期値0であることが前提　確認の必要あり
		db.Where("id = ?", data.ID).First(&result)
	}

	return result
}
