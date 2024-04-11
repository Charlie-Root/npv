/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"

	api "github.com/Charlie-Root/npv/pkg/server"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run in server mode and expose the API endpoints",
	Long: `Running in server mode will start the server and expose the API endpoints.
so remote clients can connect to the server and interact with the API.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
		
		apiHandler := api.NewHandler(database)
		http.Handle("/", apiHandler.SetupRoutes())
		http.ListenAndServe(":8080", nil)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
