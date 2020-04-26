package db

import (
	"fmt"
	"log"
	"testing"
)

func testUserCreate() {
	db := InitDB()

	db.Create(
		&User{
			ID:   "test01",
			Name: "TestUser",
			Pass: "Test",
			Auth: 0,
			//AdminGroup: []Group{
			//	{ID: "testgroup1", Name: "Group1"},
			//},
		},
	)

	if err := db.Error; err != nil {
		db.Rollback()
	}
}

func TestInitDatabase(t *testing.T) {
	fmt.Println("fasf")
	//fmt.Println(d)
	//t.Failed()

	fmt.Println("tsat")
}

func TestCreateDatabase(t *testing.T) {
	//testUserCreate()
	AddDBUser(User{
		ID:   "TestUser1",
		Name: "User1",
		Pass: "password",
		Auth: 0,
	})

	AddDBUser(User{
		ID:   "TestUser2",
		Name: "User2",
		Pass: "password",
		Auth: 0,
	})
	log.Println(GetAllDBUser())
}
