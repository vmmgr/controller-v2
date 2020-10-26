package cmd

import (
	"github.com/spf13/cobra"
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

		store.InitDB()

		log.Println("end")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
