package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core"
	toolToken "github.com/vmmgr/controller/pkg/api/core/tool/token"
	"github.com/vmmgr/controller/pkg/api/core/user"
	dbUser "github.com/vmmgr/controller/pkg/api/store/user/v0"
	"log"
	"strings"
)

func replaceUser(serverData, input, replace core.User) (core.User, error) {
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
			return core.User{}, fmt.Errorf("wrong email address")
		}
		tmp := dbUser.Get(user.Email, &core.User{Email: input.Email})
		if tmp.Err != nil {
			return replace, tmp.Err
		}
		if len(tmp.User) != 0 {
			log.Println("error: this email is already registered: " + input.Email)
			return replace, fmt.Errorf("error: this email is already registered")
		}

		mailToken, _ := toolToken.Generate(4)
		replace.Email = input.Email
		replace.MailVerify = &[]bool{false}[0]
		replace.MailToken = mailToken
	}

	//Pass
	if input.Pass == "" {
		replace.Pass = serverData.Pass
	} else {
		replace.Pass = input.Pass
	}

	return replace, nil
}

func updateAdminUser(input, replace core.User) (core.User, error) {
	//Name
	if input.Name != "" {
		replace.Name = input.Name
	}

	//E-Mail
	if input.Email != "" {
		if !strings.Contains(input.Email, "@") {
			return core.User{}, fmt.Errorf("wrong email address")
		}
		tmp := dbUser.Get(user.Email, &core.User{Email: input.Email})
		if tmp.Err != nil {
			return replace, tmp.Err
		}
		if len(tmp.User) != 0 {
			log.Println("error: this email is already registered: " + input.Email)
			return replace, fmt.Errorf("error: this email is already registered")
		}

		mailToken, _ := toolToken.Generate(4)
		replace.Email = input.Email
		replace.MailVerify = &[]bool{false}[0]
		replace.MailToken = mailToken
	}

	//Pass
	if input.Pass != "" {
		replace.Pass = input.Pass
	}

	//Level
	if input.Level != replace.Level {
		replace.Level = input.Level
	}

	return replace, nil
}
