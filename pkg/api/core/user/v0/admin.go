package v0

import (
	"fmt"
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/gin-gonic/gin"
	"github.com/vmmgr/controller/pkg/api/core"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/common"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"github.com/vmmgr/controller/pkg/api/core/tool/gen"
	"github.com/vmmgr/controller/pkg/api/core/tool/hash"
	"github.com/vmmgr/controller/pkg/api/core/tool/mail"
	"github.com/vmmgr/controller/pkg/api/core/tool/notification"
	toolToken "github.com/vmmgr/controller/pkg/api/core/tool/token"
	"github.com/vmmgr/controller/pkg/api/core/user"
	dbGroup "github.com/vmmgr/controller/pkg/api/store/group/v0"
	dbUser "github.com/vmmgr/controller/pkg/api/store/user/v0"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func AddAdmin(c *gin.Context) {
	var data core.User
	var input user.CreateAdmin

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusForbidden, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	c.BindJSON(&input)
	log.Println(input)

	if err := checkAdmin(input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	mailToken, _ := toolToken.Generate(4)

	pass := ""

	// 新規ユーザ
	if input.GroupID == 0 { //new user
		if input.Pass == "" {
			c.JSON(http.StatusInternalServerError, common.Error{Error: fmt.Sprintf("wrong pass")})
			return
		}
		data = core.User{
			GroupID:       nil,
			Name:          input.Name,
			NameEn:        input.NameEn,
			Email:         input.Mail,
			Pass:          input.Pass,
			ExpiredStatus: &[]uint{0}[0],
			Level:         1,
			MailVerify:    &[]bool{false}[0],
			MailToken:     mailToken,
		}

		// グループ所属ユーザの登録
	} else {
		pass = gen.GenerateUUID()
		log.Println("Email: " + input.Mail)
		log.Println("tmp_Pass: " + pass)

		data = core.User{
			GroupID:       &[]uint{input.GroupID}[0],
			Name:          input.Name,
			NameEn:        input.NameEn,
			Email:         input.Mail,
			Level:         input.Level,
			Pass:          strings.ToLower(hash.Generate(pass)),
			ExpiredStatus: &[]uint{0}[0],
			MailVerify:    &[]bool{false}[0],
			MailToken:     mailToken,
		}
	}

	//check exist for database
	if err := dbUser.Create(&data); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	// added auto group
	if input.GroupID == 0 {
		groupData := core.Group{
			Org:       input.Name,
			Status:    0,
			Comment:   "",
			Vlan:      0,
			Enable:    &[]bool{true}[0],
			MaxVM:     1,
			MaxCPU:    2,
			MaxMemory: 2048,
		}
		_, err := dbGroup.Create(&groupData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
			return
		}

		data.GroupID = &[]uint{groupData.ID}[0]

		err = dbUser.Update(user.UpdateAll, &data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
			return
		}
	}

	go func() {
		attachment := slack.Attachment{}
		attachment.AddField(slack.Field{Title: "E-Mail", Value: input.Mail}).
			AddField(slack.Field{Title: "GroupID", Value: strconv.Itoa(int(input.GroupID))}).
			AddField(slack.Field{Title: "Name", Value: input.Name})

		notification.SendSlack(notification.Slack{Attachment: attachment, Channel: "user"})
	}()

	if !input.MailVerify {
		if pass == "" {
			go func() {
				mail.SendMail(mail.Mail{
					ToMail:  data.Email,
					Subject: "本人確認のメールにつきまして",
					Content: " " + input.Name + "様\n\n" + "以下のリンクから本人確認を完了してください。\n" +
						config.Conf.Controller.User.URL + "/api/v1/user/verify/" + mailToken + "\n" +
						"本人確認が完了次第、ログイン可能になります。\n",
				})
			}()
		} else {
			go func() {
				mail.SendMail(mail.Mail{
					ToMail:  data.Email,
					Subject: "本人確認メールにつきまして",
					Content: " " + input.Name + "様\n\n" + "以下のリンクから本人確認を完了してください。\n" +
						config.Conf.Controller.User.URL + "/api/v1/user/verify/" + mailToken + "\n" +
						"本人確認が完了次第、ログイン可能になります。\n" + "仮パスワード: " + pass,
				})
			}()
		}
	}
	log.Println(data)
	c.JSON(http.StatusOK, user.Result{})
}

func GenerateAdmin(c *gin.Context) {
	var input core.User

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusForbidden, common.Error{Error: resultAdmin.Err.Error()})
		return
	}
	c.BindJSON(&input)

	if err := dbUser.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, user.Result{})
}

func DeleteAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusForbidden, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	if err := dbUser.Delete(&core.User{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, user.Result{})
}

func UpdateAdmin(c *gin.Context) {
	var input core.User

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusForbidden, common.Error{Error: resultAdmin.Err.Error()})
		return
	}
	c.BindJSON(&input)

	tmp := dbUser.Get(user.ID, &core.User{Model: gorm.Model{ID: input.ID}})
	if tmp.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: tmp.Err.Error()})
		return
	}

	replace, err := updateAdminUser(input, tmp.User[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "error: this email is already registered"})
		return
	}

	if err = dbUser.Update(user.UpdateAll, &replace); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, user.Result{})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusForbidden, common.Error{Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	result := dbUser.Get(user.ID, &core.User{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, result.User)
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusForbidden, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbUser.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, result.User)
	}
}
