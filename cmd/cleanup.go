/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// cleanupCmd represents the cleanup command
var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Cleaup all temporary files",
	Long: `Just a cleanup of the database file and stuff`,
	Run: func(cmd *cobra.Command, args []string) {

		err := os.Remove(c.DbFilename)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Deleted: " + c.DbFilename)
	},
}

func init() {
	rootCmd.AddCommand(cleanupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
