package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"github.com/vmmgr/controller/pkg/api/store"
	"log"
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

		db, _ := store.ConnectDB()
		result := db.AutoMigrate(
			&core.User{},
			&core.Group{},
			&core.Token{},
			&core.Notice{},
			&core.Ticket{},
			&core.Chat{},
			&core.Region{},
			&core.Zone{},
			&core.Node{},
			&core.Storage{},
			&core.NIC{},
			&core.VM{},
			&core.ImaCon{},
			&core.Image{},
			&core.Template{},
			&core.TemplatePlan{},
			&core.IP{},
		)
		log.Println(result)
		log.Println("end")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
