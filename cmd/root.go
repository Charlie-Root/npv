/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"time"

	"github.com/Charlie-Root/mtrview/pkg/config"
	"github.com/Charlie-Root/mtrview/pkg/db"
	"github.com/Charlie-Root/mtrview/pkg/logging"
	"github.com/spf13/cobra"
)

var logger = logging.NewLogger("cmd")
var c, _ = config.LoadYAML("config.yaml")
var database, _ = db.Open(c.DbType, "file:"+c.DbFilename)


var (
			COUNT            = int(c.MTRView.Count)
			TIMEOUT          = time.Duration(c.MTRView.Timeout) * time.Millisecond
			INTERVAL         = time.Duration(c.MTRView.Interval) * time.Millisecond
			HOP_SLEEP        = time.Nanosecond
			MAX_HOPS         = int(c.MTRView.MaxHops)
			MAX_UNKNOWN_HOPS = int(c.MTRView.MaxHopsUnknown)
			RING_BUFFER_SIZE = int(c.MTRView.Ringbuffer)
			PTR_LOOKUP       = bool(c.MTRView.PtrLookup)
			srcAddr          = ""
		)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mtrview",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mtrview.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	
	database.CreateTables()

}
