package db

import (
	"database/sql"
	"fmt"
	"github.com/vmmgr/controller/etc"
)

//groupdata
func AddDBGroup(data Group) bool {
	db := connectdb()
	defer db.Close()

	addDb, err := db.Prepare(`INSERT INTO "groupdata" ("name","admin","user","uuid","maxvm","maxcpu","maxmem","maxstorage","net") VALUES (?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if _, err := addDb.Exec(data.Name, data.Admin, data.User, etc.GenerateToken(), data.MaxVM, data.MaxCPU, data.MaxMem, data.MaxStorage, data.Net); err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func RemoveDBGroup(id int) bool {
	db := connectdb()
	defer db.Close()

	deletedb := "DELETE FROM groupdata WHERE id = ?"
	_, err := db.Exec(deletedb, id)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func GetDBAllGroup() []Group {
	db := *connectdb()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM groupdata")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var bg []Group
	for rows.Next() {
		var b Group
		err := rows.Scan(&b.ID, &b.Name, &b.Admin, &b.User, &b.UUID, &b.MaxVM, &b.MaxCPU, &b.MaxMem, &b.MaxStorage, &b.Net)
		if err != nil {
			fmt.Println(err)
		}
		bg = append(bg, b)
	}
	return bg
}

func GetDBGroupToken(uuid string) (Group, bool) {
	db := connectdb()
	defer db.Close()

	rows := db.QueryRow("SELECT * FROM groupdata WHERE uuid = ?", uuid)

	var b Group
	err := rows.Scan(&b.ID, &b.Name, &b.Admin, &b.User, &b.UUID, &b.MaxVM, &b.MaxCPU, &b.MaxMem, &b.MaxStorage, &b.Net)

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

func GetDBGroup(id int) (Group, bool) {
	db := connectdb()
	defer db.Close()

	rows := db.QueryRow("SELECT * FROM groupdata WHERE id = ?", id)

	var b Group
	err := rows.Scan(&b.ID, &b.Name, &b.Admin, &b.User, &b.UUID, &b.MaxVM, &b.MaxCPU, &b.MaxMem, &b.MaxStorage, &b.Net)

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

func GetDBGroupID(name string) (int, bool) {
	db := connectdb()
	defer db.Close()

	var id int
	if err := db.QueryRow("SELECT id FROM groupdata WHERE name = ?", name).Scan(&id); err != nil {
		fmt.Println(err)
		return -1, false
	}

	return id, true

}

func ChangeDBGroupName(id int, data string) bool {
	db := connectdb()
	defer db.Close()

	dbdata := "UPDATE groupdata SET name = ? WHERE id = ?"
	_, err := db.Exec(dbdata, data, id)

	if err != nil {
		fmt.Println("Error: DBUpdate Error (Group Name)")
		return false
	}

	return true
}

func ChangeDBGroupAdmin(id int, data string) bool {
	db := connectdb()
	defer db.Close()

	dbdata := "UPDATE groupdata SET admin = ? WHERE id = ?"
	_, err := db.Exec(dbdata, data, id)

	if err != nil {
		fmt.Println("Error: DBUpdate Error (Group Admin)")
		return false
	}

	return true
}

func ChangeDBGroupUser(id int, data string) bool {
	db := connectdb()
	defer db.Close()

	dbdata := "UPDATE groupdata SET user = ? WHERE id = ?"
	_, err := db.Exec(dbdata, data, id)

	if err != nil {
		fmt.Println("Error: DBUpdate Error (Group User)")
		return false
	}

	return true
}
