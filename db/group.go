package db

import "log"

//Add
func AddDBGroup(data Group) bool {
	db := InitDB()
	defer db.Close()
	db.Create(&data)
	log.Println(db.Create(&data))
	if err := db.Error; err != nil {
		db.Rollback()
		return false
	} else {
		return true
	}
}

//Delete
func DeleteDBGroup(data Group) bool {
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
func UpdateDBGroup(data Group) bool {
	db := InitDB()
	defer db.Close()
	db.Model(&data).Updates(Group{Name: data.Name, Private: data.Private, AdminUser: data.AdminUser, StandardUser: data.StandardUser})

	if err := db.Error; err != nil {
		db.Rollback()
		return false
	} else {
		return true
	}
}

//Get
func GetAllDBGroup() []Group {
	db := InitDB()
	defer db.Close()

	var group []Group
	db.Find(&group)
	return group
}

func SearchDBGroup(data Group) Group {
	db := InitDB()
	defer db.Close()

	var result Group
	db.Where("name = ?", data.Name).First(&result)

	return result
}
