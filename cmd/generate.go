/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/Charlie-Root/npv/pkg/asn"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a hosts.json from any ASN",
	Long: `Create a hosts file (hosts.json) from all resources in a given ASN:
	go run ./main.go generate <ASN>`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate called")

		Int, _ := strconv.Atoi(args[0])
		
		asn.GenerateFile(Int)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
