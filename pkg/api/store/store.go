package store

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/vmmgr/controller/pkg/api/core/group"
	"github.com/vmmgr/controller/pkg/api/core/notice"
	"github.com/vmmgr/controller/pkg/api/core/support/chat"
	"github.com/vmmgr/controller/pkg/api/core/support/ticket"
	"github.com/vmmgr/controller/pkg/api/core/token"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"github.com/vmmgr/controller/pkg/api/core/user"
	"log"
	"strconv"
)

func ConnectDB() (*gorm.DB, error) {
	user := config.Conf.DB.User
	pass := config.Conf.DB.Pass
	protocol := "tcp(" + config.Conf.DB.IP + ":" + strconv.Itoa(config.Conf.DB.Port) + ")"
	dbName := config.Conf.DB.DBName

	db, err := gorm.Open("mysql", user+":"+pass+"@"+protocol+"/"+dbName+"?parseTime=true")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitDB() {
	db, _ := ConnectDB()
	result := db.AutoMigrate(&user.User{}, &group.Group{}, &token.Token{}, &notice.Notice{},
		&ticket.Ticket{}, &chat.Chat{})
	log.Println(result.Error)
	//return nil
}
