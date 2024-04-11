/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Charlie-Root/npv/pkg/export"
	_ "github.com/Charlie-Root/npv/statik"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var port string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Fire up the server and see the results in the browser",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("started server on http://localhost:" + port)

		statikFS, err := fs.New()
		if err != nil {
			log.Fatal(err)
		}
		http.Handle("/", http.FileServer(statikFS))
		http.HandleFunc("/graph", func(writer http.ResponseWriter, request *http.Request) {
			logger.Debug("Gather graph information")
			data, _ := export.ParseGraphData(database)
			logger.Debug("Set!")
			writer.Header().Set("Content-Type", "application/json")
			if _, err := writer.Write(data); err != nil {
				fmt.Printf("failed to write response data: %s", err)
			}
		})
		if err := http.ListenAndServe(getPort(), nil); err != nil {
			panic(fmt.Sprintf("failed to start server: %s", err))
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().StringVarP(&port, "port", "p", "3000", "Port on which the server should listen")
	if err := viper.BindPFlag("port", serveCmd.PersistentFlags().Lookup("port")); err != nil {
		panic(fmt.Sprintf("faild to bind port flag: %s", err))
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func getPort() string {
	return fmt.Sprintf(":%s", port)
}
