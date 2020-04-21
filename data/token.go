package data

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/vmmgr/controller/db"
	"strconv"
	"time"
)

const (
	holdtime = 21600 //secound
)

func NewToken(user string) (string, bool) {
	userid, result := db.GetDBUserID(user)
	if result == false {
		return "ID SeachFailed", false
	}
	bt, et := TimeCalc()
	token := GenerateToken()

	db.AddDBToken(db.Token{
		Token:     token,
		Userid:    userid,
		Begintime: bt,
		Endtime:   et,
	})

	return token, true
}

func TimeCalc() (int, int) {
	begin := time.Now().Unix()
	return int(begin), int(begin + holdtime)
}

func VerifyTime(t int) bool {
	if int(time.Now().Unix()) < t {
		return true
	}
	return false
}

func DeleteExpiredToken() {
	data := db.GetDBAllToken()
	for _, a := range data {
		if a.Endtime < int(time.Now().Unix()) {
			info, result := db.RemoveDBToken(a.ID)
			if result {
				fmt.Println("---------DeleteExpiredToken---------")
				fmt.Println("Token: " + a.Token + "|UserID: " + strconv.Itoa(a.Userid) + "|BeginTime: " + strconv.Itoa(a.Begintime) + "|EndTime: " + strconv.Itoa(a.Endtime))
				fmt.Println("------------------------------------")
			} else {
				fmt.Println("Token Delete Error!!")
				fmt.Println("ErrorInfo: " + info)
			}
		}
	}
}

func GenerateToken() string {
	u, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
	}
	uuid := u.String()
	fmt.Println("GenerateToken: " + uuid)
	return uuid
}
