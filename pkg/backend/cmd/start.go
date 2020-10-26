package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vmmgr/controller/pkg/api"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	logging "github.com/vmmgr/controller/pkg/api/core/tool/log"
	"log"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start controller server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		confPath, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		if config.GetConfig(confPath) != nil {
			log.Fatalf("error config process |%v", err)
		}

		logging.WriteLog("------Application Start(Controller)------")

		api.ControllerRestAPI()
		log.Println("end")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringP("config", "c", "", "config path")
}
