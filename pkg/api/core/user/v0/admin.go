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
	dbUser "github.com/vmmgr/controller/pkg/api/store/user/v0"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func AddAdmin(c *gin.Context) {
	var input, data core.User

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusForbidden, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	c.BindJSON(&input)

	if !strings.Contains(input.Email, "@") {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("wrong email address")})
		return
	}
	if input.Name == "" {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("wrong name")})
		return
	}

	if err := check(input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	mailToken, _ := toolToken.Generate(4)

	pass := ""

	// 新規ユーザ
	if input.GroupID == nil { //new user
		if input.Pass == "" {
			c.JSON(http.StatusInternalServerError, common.Error{Error: fmt.Sprintf("wrong pass")})
			return
		}
		data = core.User{
			GroupID:    nil,
			Name:       input.Name,
			Email:      input.Email,
			Pass:       input.Pass,
			Level:      1,
			MailVerify: &[]bool{false}[0],
			MailToken:  mailToken,
		}

		// グループ所属ユーザの登録
	} else {
		pass = gen.GenerateUUID()
		log.Println("Email: " + input.Email)
		log.Println("tmp_Pass: " + pass)

		data = core.User{
			GroupID:    input.GroupID,
			Name:       input.Name,
			Email:      input.Email,
			Level:      input.Level,
			Pass:       strings.ToLower(hash.Generate(pass)),
			MailVerify: &[]bool{false}[0],
			MailToken:  mailToken,
		}
	}

	//check exist for database
	if err := dbUser.Create(&data); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	} else {
		attachment := slack.Attachment{}
		attachment.AddField(slack.Field{Title: "E-Mail", Value: input.Email}).
			AddField(slack.Field{Title: "GroupID", Value: strconv.Itoa(int(*input.GroupID))}).
			AddField(slack.Field{Title: "Name", Value: input.Name})

		notification.SendSlack(notification.Slack{Attachment: attachment, Channel: "user"})

		if pass == "" {
			mail.SendMail(mail.Mail{
				ToMail:  data.Email,
				Subject: "本人確認のメールにつきまして",
				Content: " " + input.Name + "様\n\n" + "以下のリンクから本人確認を完了してください。\n" +
					config.Conf.Controller.User.URL + "/api/v1/user/verify/" + mailToken + "\n" +
					"本人確認が完了次第、ログイン可能になります。\n",
			})
		} else {
			mail.SendMail(mail.Mail{
				ToMail:  data.Email,
				Subject: "本人確認メールにつきまして",
				Content: " " + input.Name + "様\n\n" + "以下のリンクから本人確認を完了してください。\n" +
					config.Conf.Controller.User.URL + "/api/v1/user/verify/" + mailToken + "\n" +
					"本人確認が完了次第、ログイン可能になります。\n" + "仮パスワード: " + pass,
			})
		}

		c.JSON(http.StatusOK, user.Result{})
	}
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
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, user.Result{User: result.User})
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
		c.JSON(http.StatusOK, user.Result{User: result.User})
	}
}
