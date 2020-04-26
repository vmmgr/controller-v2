package db

import "log"

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
