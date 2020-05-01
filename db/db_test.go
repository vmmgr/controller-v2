package db

import (
	"fmt"
	"testing"
)

func TestInitDatabase(t *testing.T) {
	InitCreateDB()
}

func TestCreateDatabase(t *testing.T) {
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

	AddDBToken(Token{
		Token: "sdfafsafdsfsafdskjgldsfak",
	})

	result := AddDBGroup(Group{
		ID:   "TestGroup1",
		Name: "Group1",
		//AdminUser: []User{{ID: "TestUser1"}},
		Private: false,
	})

	AddDBUser(User{
		ID:   "TestUser3",
		Name: "User3",
		Pass: "password",
		//AdminGroup: []Group{{ID: "TestGroup1"}},
		Auth: 0,
	})

	fmt.Println(result)
	fmt.Println("=====User=====")
	fmt.Println(GetAllDBUser())
	fmt.Println("=====Group=====")
	fmt.Println(GetAllDBGroup())
	fmt.Println("=====Token=====")
	fmt.Println(GetAllDBToken())
	fmt.Println("==========")
}

func TestJoinGroupDatabase(t *testing.T) {
	db := InitDB()
	defer db.Close()

	//var user []User
	var group Group
	group = Group{ID: "TestGroup1"}
	//user = append(user, User{ID: "TestUser2"})

	//db.Model(&group).Association("AdminUser").Append(user)

	//UpdateDBGroupUser(Group{
	//	ID:        "TestGroup1",
	//	AdminUser: []User{{ID: "TestUser1"}},
	//})
	//user1 := User{ID: "TestUser1"}
	//db.Model(&user1).Association("AdminGroup").Append(Group{Name: "Group1"})
	fmt.Println("===========")
	//user2 := User{ID: "TestUser1"}
	//var group Group
	//db.Model(&user2).Association("UserGroup").Append(Group{ID: "TestGroup1"})
	//db.Model(&user2).Association("User").Find(&group)
	fmt.Println("=====User=====")
	fmt.Println(GetAllDBUser())
	fmt.Println("=====Group=====")
	fmt.Println(GetAllDBGroup())
	fmt.Println("==========")
	db.Model(&group).Association("StandardUser").Append(&User{ID: "TestUser1"})

	//UpdateDBGroupUser(Group{
	//	ID:           "TestGroup1",
	//	StandardUser: user,
	//})
	fmt.Println("=====User=====")
	fmt.Println(GetAllDBUser())
	fmt.Println("=====Group=====")
	fmt.Println(GetAllDBGroup())
	fmt.Println("==========")

	//db.Model(&group).Related(&user)

}
