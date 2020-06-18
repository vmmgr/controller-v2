package db

//Add
func AddDBGroup(data Group) bool {
	db := InitDB()
	defer db.Close()
	//db.Table("group").CreateTable(&data)
	db.Create(&data)

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
	db.Model(&data).Updates(Group{Name: data.Name, Private: data.Private})

	if err := db.Error; err != nil {
		db.Rollback()
		return false
	} else {
		return true
	}
}

//Add DBGroupUser
func AddDBGroupUser(data Group, userID string, admin bool) bool {
	db := InitDB()
	defer db.Close()

	if admin {
		db.Model(&data).Association("AdminUser").Append(&User{ID: userID})
	} else {
		db.Model(&data).Association("StandardUser").Append(&User{ID: userID})
	}

	if err := db.Error; err != nil {
		db.Rollback()
		return false
	} else {
		return true
	}
}

//Delete DBGroupUser
func DeleteDBGroupUser(data Group, userID string, admin bool) bool {
	db := InitDB()
	defer db.Close()

	if admin {
		db.Model(&data).Association("AdminUser").Delete(&User{ID: userID})
	} else {
		db.Model(&data).Association("StandardUser").Delete(&User{ID: userID})
	}
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
	//db.Table("group").Find(&group)
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
