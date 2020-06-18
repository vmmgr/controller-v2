package data

import "github.com/vmmgr/controller/db"

func ExistUser(name string) bool {
	//true: exists username  false: not exists username
	result := db.SearchDBUser(db.User{Name: name})
	if name == result.Name {
		return true
	} else {
		return false
	}
}
