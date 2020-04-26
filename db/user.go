package db

import "log"

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
