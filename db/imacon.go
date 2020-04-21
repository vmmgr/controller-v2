package db

import (
	"database/sql"
	"fmt"
	"log"
)

func AddDBImaCon(data ImaCon) bool {
	db := *connectdb()
	defer db.Close()

	if data.ID == 0 {
		addDb, err := db.Prepare(`INSERT INTO "imacon" ("hostname","ip","status") VALUES (?,?,?)`)
		if err != nil {
			fmt.Println("DBError!!")
			return false
		}

		if _, err := addDb.Exec(data.HostName, data.IP, data.Status); err != nil {
			fmt.Println("Add Error!!")
			return false
		}
	} else {
		addDb, err := db.Prepare(`INSERT INTO "imacon" ("id","hostname","ip","status") VALUES (?,?,?,?)`)
		if err != nil {
			fmt.Println("DBError!!")
			return false
		}

		if _, err := addDb.Exec(data.ID, data.HostName, data.IP, data.Status); err != nil {
			fmt.Println("Add Error!!")
			return false
		}
	}
	return true
}

func RemoveDBImaCon(id int) bool {
	db := *connectdb()
	defer db.Close()

	deleteDb := "DELETE FROM imacon WHERE id = ?"
	_, err := db.Exec(deleteDb, id)
	if err != nil {
		fmt.Println("Delete Failed!!")
		return false
	}
	return true
}

func GetDBImaCon(id int) (ImaCon, bool) {
	db := *connectdb()
	defer db.Close()

	rows := db.QueryRow("SELECT * FROM imacon WHERE id = ?", id)

	var b ImaCon
	err := rows.Scan(&b.ID, &b.HostName, &b.IP, &b.Status)

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

func GetDBAllImaCon() []ImaCon {
	db := *connectdb()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM imacon")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var bg []ImaCon
	for rows.Next() {
		var b ImaCon
		err := rows.Scan(&b.ID, &b.HostName, &b.IP, &b.Status)
		if err != nil {
			log.Println(err)
		}
		bg = append(bg, b)
	}
	return bg
}
