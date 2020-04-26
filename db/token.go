package db

//Add
func AddToken(token Token) bool {
	db := InitDB()
	defer db.Close()
	db.Create(&token)

	if err := db.Error; err != nil {
		db.Rollback()
		return false
	} else {
		return true
	}
}

//Delete
func DeleteToken(token Token) bool {
	db := InitDB()
	defer db.Close()
	db.Delete(&token)

	if err := db.Error; err != nil {
		db.Rollback()
		return false
	} else {
		return true
	}
}

//Update
func UpdateDBToken(data Token) bool {
	db := InitDB()
	defer db.Close()
	db.Model(&data).Updates(Token{Token: data.Token})

	if err := db.Error; err != nil {
		db.Rollback()
		return false
	} else {
		return true
	}
}

//Get
func GetAllDBToken() []Token {
	db := InitDB()
	defer db.Close()

	var token []Token
	db.Find(&token)
	return token
}

func SearchToken(token Token) Token {
	db := InitDB()
	defer db.Close()
	var result Token
	//search UserName and UserID
	if token.Token != "" {
		db.Where("token = ?", token.Token).First(&result)
	}
	return result
}
