package v0

import (
	"fmt"
	toolToken "github.com/vmmgr/controller/pkg/api/core/tool/token"
	"github.com/vmmgr/controller/pkg/api/core/user"
	dbUser "github.com/vmmgr/controller/pkg/api/store/user/v0"
	"log"
	"strings"
)

func replaceUser(serverData, input, replace user.User) (user.User, error) {
	updateInfo := 0
	//Name
	if input.Name == "" {
		replace.Name = serverData.Name
	} else {
		replace.Name = input.Name
	}

	//E-Mail
	if input.Email == "" {
		replace.Email = serverData.Email
		replace.MailToken = serverData.MailToken
		replace.MailVerify = serverData.MailVerify
	} else {
		if !strings.Contains(input.Email, "@") {
			return user.User{}, fmt.Errorf("wrong email address")
		}
		tmp := dbUser.Get(user.Email, &user.User{Email: input.Email})
		if tmp.Err != nil {
			return replace, tmp.Err
		}
		if len(tmp.User) != 0 {
			log.Println("error: this email is already registered: " + input.Email)
			return replace, fmt.Errorf("error: this email is already registered")
		}

		mailToken, _ := toolToken.Generate(4)
		replace.Email = input.Email
		replace.MailVerify = false
		replace.MailToken = mailToken
	}

	//Pass
	if input.Pass == "" {
		replace.Pass = serverData.Pass
	} else {
		replace.Pass = input.Pass
	}

	if serverData.Status == 0 && updateInfo == 5 {
		replace.Status = 1
	} else if serverData.Status == 0 && updateInfo < 5 {
		return replace, fmt.Errorf("lack of information")
	}

	return replace, nil
}

func updateAdminUser(input, replace user.User) (user.User, error) {
	//Name
	if input.Name != "" {
		replace.Name = input.Name
	}

	//E-Mail
	if input.Email != "" {
		if !strings.Contains(input.Email, "@") {
			return user.User{}, fmt.Errorf("wrong email address")
		}
		tmp := dbUser.Get(user.Email, &user.User{Email: input.Email})
		if tmp.Err != nil {
			return replace, tmp.Err
		}
		if len(tmp.User) != 0 {
			log.Println("error: this email is already registered: " + input.Email)
			return replace, fmt.Errorf("error: this email is already registered")
		}

		mailToken, _ := toolToken.Generate(4)
		replace.Email = input.Email
		replace.MailVerify = false
		replace.MailToken = mailToken
	}

	//Pass
	if input.Pass != "" {
		replace.Pass = input.Pass
	}

	//Status
	if input.Status != replace.Status {
		replace.Status = input.Status
	}

	//Level
	if input.Level != replace.Level {
		replace.Level = input.Level
	}

	return replace, nil
}
