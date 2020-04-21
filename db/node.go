package db

import (
	"database/sql"
	"fmt"
	"log"
)

//Node
func AddDBNode(data Node) bool {
	db := *connectdb()
	defer db.Close()

	if data.ID == 0 {
		addDb, err := db.Prepare(`INSERT INTO "node" ("hostname","ip","path","onlyadmin","maxcpu","maxmem","status") VALUES (?,?,?,?,?,?,?)`)
		if err != nil {
			fmt.Println("DBError!!")
			return false
		}

		if _, err := addDb.Exec(data.HostName, data.IP, data.Path, data.OnlyAdmin, data.MaxCPU, data.MaxMem, data.Status); err != nil {
			fmt.Println("Add Error!!")
			return false
		}
	} else {
		addDb, err := db.Prepare(`INSERT INTO "node" ("id","hostname","ip","path","onlyadmin","maxcpu","maxmem","status") VALUES (?,?,?,?,?,?,?,?)`)
		if err != nil {
			fmt.Println("DBError!!")
			return false
		}

		if _, err := addDb.Exec(data.ID, data.HostName, data.IP, data.Path, data.OnlyAdmin, data.MaxCPU, data.MaxMem, data.Status); err != nil {
			fmt.Println("Add Error!!")
			return false
		}
	}
	return true
}

func RemoveDBNode(id int) bool {
	db := *connectdb()
	defer db.Close()

	deleteDb := "DELETE FROM node WHERE id = ?"
	_, err := db.Exec(deleteDb, id)
	if err != nil {
		fmt.Println("Delete Failed!!")
		return false
	}
	return true
}

func NodeDBStatusUpdate(id, status int) bool {
	db := *connectdb()
	defer db.Close()

	cmd := "UPDATE node SET status = ? WHERE id = ?"
	_, err := db.Exec(cmd, status, id)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func GetDBNodeID(id int) (Node, bool) {
	db := *connectdb()
	defer db.Close()

	rows := db.QueryRow("SELECT * FROM node WHERE id = ?", id)

	var b Node
	err := rows.Scan(&b.ID, &b.HostName, &b.IP, &b.Path, &b.OnlyAdmin, &b.MaxCPU, &b.MaxMem, &b.Status)

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

func GetDBAllNode() []Node {
	db := *connectdb()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM node")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var bg []Node
	for rows.Next() {
		var b Node
		err := rows.Scan(&b.ID, &b.HostName, &b.IP, &b.Path, &b.OnlyAdmin, &b.MaxCPU, &b.MaxMem, &b.Status)
		if err != nil {
			log.Println(err)
		}
		bg = append(bg, b)
	}
	return bg
}
