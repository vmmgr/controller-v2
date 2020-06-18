package data

import (
	"fmt"
	"github.com/vmmgr/controller/db"
)

type userData struct {
	id     int
	aGroup []*db.Group
	uGroup []*db.Group
	err    error
}

func UserCertification(token string) int {
	data := db.SearchToken(db.Token{Token: token})
	return data.UserID
}

func getUserData(token string) userData {
	tokenData := db.GetAllDBToken()
	for _, t := range tokenData {
		if t.Token == token {
			user := db.SearchDBUser(db.User{ID: t.UserID})
			return userData{
				id:     t.UserID,
				aGroup: user.AdminGroup,
				uGroup: user.UserGroup,
				err:    nil,
			}
		}
	}
	return userData{err: fmt.Errorf("Error: data none ")}
}

// 0:AdministratorGroup 1: AdminGroup 5:UserGroup 10:None 20:Error
func VerifyGroup(token string) uint {
	user := getUserData(token)
	if user.err != nil {
		return 20
	}
	for _, g := range user.aGroup {
		if g.ID == 0 {
			return 0
		}
	}
	return 10
}

// 0:AdministratorGroup 1: AdminGroup 5:UserGroup 10:None 20:Error
func VerifySameGroup(token string, groupID int) uint {
	user := getUserData(token)
	if user.err != nil {
		return 20
	}
	for _, g := range user.aGroup {
		if groupID == g.ID && g.ID == 1 {
			return 0
		}
		if groupID == g.ID {
			return 1
		}
	}
	for _, g := range user.uGroup {
		if groupID == g.ID {
			return 5
		}
	}
	return 10
}
