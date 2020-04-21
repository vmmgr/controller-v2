package db

import (
	"database/sql"
	"fmt"
)

func AddDBToken(data Token) (string, bool) {
	db := connectdb()
	defer db.Close()

	addDb, err := db.Prepare(`INSERT INTO "tokendata" ("token","userid","groupid","begintime","endtime") VALUES (?,?,?,?,?)`)
	if err != nil {
		fmt.Println(err)
		return "Token Add Error", false
	}

	if _, err := addDb.Exec(data.Token, data.Userid, data.Groupid, data.Begintime, data.Endtime); err != nil {
		fmt.Println(err)
		return "Token Add Error", false
	}
	return "ok", true
}

func RemoveDBToken(id int) (string, bool) {
	db := connectdb()
	defer db.Close()

	deleteDb := "DELETE FROM tokendata WHERE id = ?"
	_, err := db.Exec(deleteDb, id)
	if err != nil {
		fmt.Println(err)
		return "Token Delete Error!!", false
	}
	return "ok", true
}

func GetDBAllToken() []Token {
	db := *connectdb()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM tokendata")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var bg []Token
	for rows.Next() {
		var b Token
		err := rows.Scan(&b.ID, &b.Token, &b.Userid, &b.Groupid, &b.Begintime, &b.Endtime)
		if err != nil {
			fmt.Println(err)
		}
		bg = append(bg, b)
	}
	return bg
}

func GetDBToken(token string) (Token, bool) {
	db := connectdb()
	defer db.Close()

	rows := db.QueryRow("SELECT * FROM tokendata WHERE token = ?", token)

	var b Token
	err := rows.Scan(&b.ID, &b.Token, &b.Userid, &b.Groupid, &b.Begintime, &b.Endtime)

	switch {
	case err == sql.ErrNoRows:
		fmt.Printf("Not found")
		return b, false
	case err != nil:
		fmt.Println(err)
		fmt.Println("Error: DBError")
		return b, false
	default:
		return b, true
	}
}
