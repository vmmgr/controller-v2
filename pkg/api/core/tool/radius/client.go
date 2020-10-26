package radius

import (
	"context"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"log"
	"strconv"
)

func Client() {
	packet := radius.New(radius.CodeAccessRequest, []byte(`secret`))
	rfc2865.UserName_SetString(packet, config.Conf.Radius.User)
	rfc2865.UserPassword_SetString(packet, config.Conf.Radius.Pass)
	response, err := radius.Exchange(context.Background(), packet, config.Conf.Radius.Host+":"+strconv.Itoa(config.Conf.Radius.Port))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Code:", response.Code)
}
