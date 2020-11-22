package cmd

import (
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
	"github.com/vmmgr/controller/pkg/api/core/group"
	"github.com/vmmgr/controller/pkg/api/core/node"
	nodeNIC "github.com/vmmgr/controller/pkg/api/core/node/nic"
	nodeStorage "github.com/vmmgr/controller/pkg/api/core/node/storage"
	"github.com/vmmgr/controller/pkg/api/core/notice"
	"github.com/vmmgr/controller/pkg/api/core/region"
	"github.com/vmmgr/controller/pkg/api/core/region/zone"
	"github.com/vmmgr/controller/pkg/api/core/support/chat"
	"github.com/vmmgr/controller/pkg/api/core/support/ticket"
	"github.com/vmmgr/controller/pkg/api/core/token"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"github.com/vmmgr/controller/pkg/api/core/user"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"log"
	"strconv"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init database",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		confPath, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		if config.GetConfig(confPath) != nil {
			log.Fatalf("error config process |%v", err)
		}

		db, err := gorm.Open("mysql", config.Conf.DB.User+":"+config.Conf.DB.Pass+"@"+
			"tcp("+config.Conf.DB.IP+":"+strconv.Itoa(config.Conf.DB.Port)+")"+"/"+config.Conf.DB.DBName+"?parseTime=true")
		if err != nil {
			panic(err)
		}
		result := db.AutoMigrate(&user.User{}, &group.Group{}, &token.Token{}, &notice.Notice{},
			&ticket.Ticket{}, &chat.Chat{}, &region.Region{}, &zone.Zone{}, &node.Node{}, &nodeStorage.Storage{},
			&nodeNIC.NIC{}, &vm.VM{})
		log.Println(result.Error)
		log.Println("end")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
