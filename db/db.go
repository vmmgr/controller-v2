package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"time"
)

type Database struct {
	DB     *gorm.DB
	TestDB *gorm.DB
}

const DBPath = "./controller.db"

type Node struct {
	gorm.Model
	HostName  string
	IP        string
	Path      string
	OnlyAdmin int
	MaxCPU    int
	MaxMem    int
	Active    int
}

type ImaCon struct {
	gorm.Model
	HostName string
	IP       string
	Status   int
}

type User struct {
	ID         string `gorm:primary_key`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Name       string
	Pass       string
	Auth       int
	AdminGroup []Group `gorm:"many2many:users_admingroups"`
	UserGroup  []Group `gorm:"many2many:users_usergroups"`
}

type Group struct {
	ID           string `gorm:primary_key`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Name         string
	AdminUser    []User `gorm:"many2many:users_admingroups"`
	StandardUser []User `gorm:"many2many:users_usergroups"`
	Private      bool   `gorm:"default":true`
	//MaxVM      int    `gorm:"default":2`
	//MaxCPU     int    `gorm:"default":4`
	//MaxMem     int    `gorm:"default":4096`
	//MaxStorage int
	//Net        string
}

type Net struct {
	gorm.Model
	Name   string
	Bridge string
}

type Token struct {
	gorm.Model
	Token     string
	UserID    string
	User      int `gorm:"foreignkey:ID;association_foreignkey:UserID"`
	Begintime time.Time
	Endtime   time.Time
}

type Progress struct {
	gorm.Model
	VMName    string
	UUID      string
	StartTime int
}

func InitDB() *gorm.DB {
	return initSQLite3()
}

func initSQLite3() *gorm.DB {
	db, err := gorm.Open("sqlite3", DBPath)
	if err != nil {
		log.Println("SQL open error")
	}
	//db.LogMode(true)
	db.AutoMigrate(&User{}, &Group{}, &Token{})

	return db
}

func initMySQL() *gorm.DB {
	db, err := gorm.Open("mysql", DBPath)
	if err != nil {
		log.Println("SQL open error")
	}
	//db.LogMode(true)
	db.AutoMigrate(&User{}, &Group{}, &Token{})

	return db
}

//func InitDB() bool {
//	//Node data
//	createdb(`CREATE TABLE IF NOT EXISTS "node" ("id" INTEGER PRIMARY KEY, "hostname" VARCHAR(255), "ip" VARCHAR(255), "path" VARCHAR(2000), "onlyadmin" INT,"maxcpu" INT ,"maxmem" INT, "status" INT)`)
//	//imacon data
//	createdb(`CREATE TABLE IF NOT EXISTS "imacon" ("id" INTEGER PRIMARY KEY, "hostname" VARCHAR(255), "ip" VARCHAR(255), "status" INT)`)
//	//user data
//	createdb(`CREATE TABLE IF NOT EXISTS "userdata" ("id" INTEGER PRIMARY KEY, "name" VARCHAR(255), "pass" VARCHAR(255))`)
//	//group data
//	createdb(`CREATE TABLE IF NOT EXISTS "groupdata" ("id" INTEGER PRIMARY KEY, "name" VARCHAR(255),"admin" VARCHAR(500),"user" VARCHAR(2000),"uuid" VARCHAR(20000),"maxvm" INT,"maxcpu" INT,"maxmem" INT,"maxstorage" INT,"net" VARCHAR(255))`)
//	//token data
//	createdb(`CREATE TABLE IF NOT EXISTS "tokendata" ("id" INTEGER PRIMARY KEY, "token" VARCHAR(1000), "userid" INT,"groupid" INT,"begintime" INT,"endtime" INT)`)
//	//progress data
//	createdb(`CREATE TABLE IF NOT EXISTS "progress" ("id" INTEGER PRIMARY KEY, "vmname" VARCHAR(255), "uuid" VARCHAR(255), "starttime" INT)`)
//	return true
//}
