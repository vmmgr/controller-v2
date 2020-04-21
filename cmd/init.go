package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vmmgr/controller/db"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize",
	Long: `initialize command. For example:

database init: init database
client init:     init client`,
}
var initdbCmd = &cobra.Command{
	Use:   "db",
	Short: "db init",
	Long: `db init cmd
For example:

db init [Admin Name] [Admin Pass]`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("----Initalize Database----")
		result := db.InitDB()
		fmt.Println(result)
		db.AddDBUser(db.User{Name: args[0], Pass: args[1]})
		db.AddDBGroup(db.Group{
			Name:       "admin",
			Admin:      args[0],
			User:       args[0],
			MaxCPU:     100,
			MaxMem:     10240000,
			MaxStorage: 1000000000,
		})
		return nil
	},
}
var initNodeCmd = &cobra.Command{
	Use:   "client",
	Short: "client init",
	Long:  "client init cmd",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Not implemented")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.AddCommand(initdbCmd)
	initCmd.AddCommand(initNodeCmd)
}
