package db

//Add
func AddDBUser(data User) bool {
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
func DeleteDBUser(data User) bool {
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
func UpdateDBUser(data User) bool {
	db := InitDB()
	defer db.Close()
	db.Model(&data).Updates(User{Name: data.Name, Pass: data.Pass})

	if err := db.Error; err != nil {
		db.Rollback()
		return false
	} else {
		return true
	}
}

//Get
func GetAllDBUser() []User {
	db := InitDB()
	defer db.Close()

	var user []User
	db.Find(&user)
	return user
}

func SearchDBUser(data User) User {
	db := InitDB()
	defer db.Close()

	var result User
	//search UserName and UserID
	if data.Name != "" {
		db.Where("name = ?", data.Name).First(&result)
	} else if data.ID != 0 { //初期値0であることが前提　確認の必要あり
		db.Where("id = ?", data.ID).First(&result)
	}

	return result
}
