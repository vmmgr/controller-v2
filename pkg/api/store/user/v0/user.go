package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/core/user"
	"github.com/vmmgr/controller/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(u *core.User) error {
	result := Get(user.Email, &core.User{Email: u.Email})

	if len(result.User) != 0 && result.Err == nil {
		log.Println(result.Err)
		log.Println("error: this email is already registered: " + u.Email)
		return fmt.Errorf("error: this email is already registered")
	}

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}

	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	resultDB := db.Create(&u)

	return resultDB.Error
}

func Delete(u *core.User) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	return db.Delete(u).Error
}

func Update(base int, u *core.User) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	var result *gorm.DB

	if user.UpdateVerifyMail == base {
		result = db.Model(&core.User{Model: gorm.Model{ID: u.ID}}).Updates(core.User{MailVerify: u.MailVerify})
	} else if user.UpdateInfo == base {
		result = db.Model(&core.User{Model: gorm.Model{ID: u.ID}}).Updates(core.User{
			Name:       u.Name,
			Email:      u.Email,
			Pass:       u.Pass,
			MailVerify: u.MailVerify,
			MailToken:  u.MailToken,
		})
	} else if user.UpdateGroupID == base {
		result = db.Model(&core.User{Model: gorm.Model{ID: u.ID}}).Updates(core.User{GroupID: u.GroupID})
	} else if user.UpdateLevel == base {
		result = db.Model(&core.User{Model: gorm.Model{ID: u.ID}}).Updates(core.User{Level: u.Level})
	} else if user.UpdateAll == base {
		result = db.Model(&core.User{Model: gorm.Model{ID: u.ID}}).Updates(core.User{
			GroupID:    u.GroupID,
			Name:       u.Name,
			Email:      u.Email,
			Pass:       u.Pass,
			MailVerify: u.MailVerify,
			MailToken:  u.MailToken,
		})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}

	return result.Error
}

// value of base can reference from api/core/user/interface.go
func Get(base int, u *core.User) user.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return user.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return user.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var userStruct []core.User

	if base == user.ID { //ID
		err = db.First(&userStruct, u.ID).Error
	} else if base == user.GID { //GroupID
		err = db.Where("group_id = ?", u.GroupID).Find(&userStruct).Error
	} else if base == user.Email { //Mail
		err = db.Where("email = ?", u.Email).Find(&userStruct).Error
	} else if base == user.MailToken { //Token
		err = db.Where("mail_token = ?", u.MailToken).Find(&userStruct).Error
	} else {
		log.Println("base select error")
		return user.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}

	return user.ResultDatabase{User: userStruct, Err: err}
}

func GetAll() user.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return user.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return user.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var users []core.User
	err = db.Find(&users).Error
	return user.ResultDatabase{User: users, Err: err}
}
