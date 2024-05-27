/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/kaffeed/bingoscape/db"
	"github.com/kaffeed/bingoscape/services"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add management user",
	Long: `Adds a new management User to 
to quickly create a Cobra application.`,
	ValidArgs: []string{"name", "password"},
	Run: func(cmd *cobra.Command, args []string) {

		store, err := db.NewStore(os.Getenv("DB"))
		if err != nil {
			fmt.Printf("%v", err)
			return
		}

		us := services.NewUserServices(store)
		us.CreateUser(services.User{
			Username:     "",
			Password:     "",
			IsManagement: true,
		})

		fmt.Println(services.User{})
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
