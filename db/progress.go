package db

import (
	"database/sql"
	"fmt"
	"log"
)

func AddDBProgress(data Progress) bool {
	db := *connectdb()
	defer db.Close()

	addDb, err := db.Prepare(`INSERT INTO "progress" ("vmname","uuid","starttime") VALUES (?,?,?)`)
	if err != nil {
		fmt.Println("DBError!!")
		return false
	}

	if _, err := addDb.Exec(data.VMName, data.UUID, data.StartTime); err != nil {
		fmt.Println("Add Error!!")
		return false
	}

	return true
}

func RemoveDBProgress(id int) bool {
	db := *connectdb()
	defer db.Close()

	deleteDb := "DELETE FROM progress WHERE id = ?"
	_, err := db.Exec(deleteDb, id)
	if err != nil {
		fmt.Println("Delete Failed!!")
		return false
	}
	return true
}

func GetDBProgressVMName(vmname string) (Progress, bool) {
	db := *connectdb()
	defer db.Close()

	rows := db.QueryRow("SELECT * FROM progress WHERE vmname = ?", vmname)

	var b Progress
	err := rows.Scan(&b.ID, &b.VMName, &b.UUID, &b.StartTime)

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

func GetDBAllProgress() []Progress {
	db := *connectdb()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM progress")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var bg []Progress
	for rows.Next() {
		var b Progress
		err := rows.Scan(&b.ID, &b.VMName, &b.UUID, &b.StartTime)
		if err != nil {
			log.Println(err)
		}
		bg = append(bg, b)
	}
	return bg
}
