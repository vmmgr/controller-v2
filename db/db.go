package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

const DBPath = "./controller.db"

type Node struct {
	ID        int
	HostName  string
	IP        string
	Path      string
	OnlyAdmin int
	MaxCPU    int
	MaxMem    int
	Status    int
}

type ImaCon struct {
	ID       int
	HostName string
	IP       string
	Status   int
}

type User struct {
	ID   int
	Name string
	Pass string
	Auth int
}

type Group struct {
	ID         int
	Name       string
	Admin      string
	User       string
	UUID       string
	MaxVM      int
	MaxCPU     int
	MaxMem     int
	MaxStorage int
	Net        string
}

type Token struct {
	ID        int
	Token     string
	Userid    int
	Groupid   string
	Begintime int
	Endtime   int
}

type Progress struct {
	ID        int
	VMName    string
	UUID      string
	StartTime int
}

func connectdb() *sql.DB {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		fmt.Println("SQL open error")
		fmt.Println(err)
	}

	return db
}

func createdb(database string) bool {
	db := *connectdb()
	defer db.Close()

	_, err := db.Exec(database)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func InitDB() bool {
	//Node data
	createdb(`CREATE TABLE IF NOT EXISTS "node" ("id" INTEGER PRIMARY KEY, "hostname" VARCHAR(255), "ip" VARCHAR(255), "path" VARCHAR(2000), "onlyadmin" INT,"maxcpu" INT ,"maxmem" INT, "status" INT)`)
	//imacon data
	createdb(`CREATE TABLE IF NOT EXISTS "imacon" ("id" INTEGER PRIMARY KEY, "hostname" VARCHAR(255), "ip" VARCHAR(255), "status" INT)`)
	//user data
	createdb(`CREATE TABLE IF NOT EXISTS "userdata" ("id" INTEGER PRIMARY KEY, "name" VARCHAR(255), "pass" VARCHAR(255))`)
	//group data
	createdb(`CREATE TABLE IF NOT EXISTS "groupdata" ("id" INTEGER PRIMARY KEY, "name" VARCHAR(255),"admin" VARCHAR(500),"user" VARCHAR(2000),"uuid" VARCHAR(20000),"maxvm" INT,"maxcpu" INT,"maxmem" INT,"maxstorage" INT,"net" VARCHAR(255))`)
	//token data
	createdb(`CREATE TABLE IF NOT EXISTS "tokendata" ("id" INTEGER PRIMARY KEY, "token" VARCHAR(1000), "userid" INT,"groupid" INT,"begintime" INT,"endtime" INT)`)
	//progress data
	createdb(`CREATE TABLE IF NOT EXISTS "progress" ("id" INTEGER PRIMARY KEY, "vmname" VARCHAR(255), "uuid" VARCHAR(255), "starttime" INT)`)
	return true
}
