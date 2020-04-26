package db

import "log"

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
	db.Model(&data).Updates(User{Name: data.Name, Pass: data.Pass, AdminGroup: data.AdminGroup, UserGroup: data.UserGroup})

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
	log.Println(db.Find(&user))
	return user
}

func SearchDBUser(data User) User {
	db := InitDB()
	defer db.Close()

	var result User
	//search UserName and UserID
	if data.Name != "" {
		db.Where("name = ?", data.Name).First(&result)
	} else if data.ID != "" {
		db.Where("id = ?", data.ID).First(&result)
	}

	return result
}
