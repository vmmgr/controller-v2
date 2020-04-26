package db

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
