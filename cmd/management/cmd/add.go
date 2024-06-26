/*
Copyright © 2024 Sebastian Schubert <sschubert932@gmail.com>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kaffeed/bingoscape/app/db"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var username, password string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add management user",
	Long: `Adds a new management User to 
to quickly create a Cobra application.`,
	ValidArgs: []string{"username", "password"},
	Run: func(cmd *cobra.Command, args []string) {
		if username == "" || password == "" {
			fmt.Println("username and password must be provided")
			return
		}

		ctx := context.Background()
		connpool, err := pgxpool.New(ctx, os.Getenv("DB_URL"))

		if err != nil {
			log.Fatalf("Error connecting to db: %#v", err)
		}
		defer connpool.Close()
		store := db.New(connpool)

		if err != nil {
			log.Fatalf("failed to create store: %s", err)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
		if err != nil {
			log.Fatalf("Error connecting to db: %#v", err)
		}

		err = store.CreateLogin(ctx, db.CreateLoginParams{
			Password:     string(hashedPassword),
			Name:         username,
			IsManagement: true,
		})
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

	addCmd.Flags().StringVarP(&username, "username", "u", "", "Username for login")
	addCmd.Flags().StringVarP(&password, "password", "p", "", "Password for login")

	// Mark flags as required
	addCmd.MarkFlagRequired("username")
	addCmd.MarkFlagRequired("password")
}
