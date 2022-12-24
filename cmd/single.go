/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"sync"

	"github.com/Charlie-Root/mtrview/pkg/mtr"
	"github.com/Charlie-Root/mtrview/pkg/parser"
	tm "github.com/buger/goterm"
	"github.com/spf13/cobra"
)

// singleCmd represents the single command
var singleCmd = &cobra.Command{
	Use:   "single",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {


		m, ch, err := mtr.NewMTR("telia.se", srcAddr, TIMEOUT, INTERVAL, HOP_SLEEP, MAX_HOPS, MAX_UNKNOWN_HOPS, RING_BUFFER_SIZE, PTR_LOOKUP)

		if err != nil {
			logger.Error(err.Error())
		}

		tm.Clear()
		mu := &sync.Mutex{}

		go func(ch chan struct{}) {
			for {
				mu.Lock()
				<-ch
				mu.Unlock()
			}
		}(ch)

		m.Run(ch, COUNT)

		parser.ParseResults(*m, *database)
		print("Results saved.")
	},
}

func init() {
	runCmd.AddCommand(singleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// singleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// singleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
