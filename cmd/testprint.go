package cmd

import (
	"github.com/avanier/gorecptprint/lib/comm"
	"github.com/avanier/gorecptprint/lib/print"
	"github.com/spf13/cobra"
)

var testPrintCmd = &cobra.Command{
	Use:   "test-print",
	Short: "performs a test print",
	Run: func(cmd *cobra.Command, args []string) {
		comm.Init()
		print.PrintTest()
	},
}

// func init() {
// 	var defaultSerialPort = "/dev/ttyS0"
// 	var defaultWriteChunkSize = uint(192)

// 	testPrintCmd.PersistentFlags().String("port", defaultSerialPort, "serial port to use for communication")
// 	testPrintCmd.PersistentFlags().Uint("write-chunk-size", defaultWriteChunkSize, "the size in bytes of every chunk of data sent to the printer")

// 	viper.BindPFlags(testPrintCmd.PersistentFlags())
// 	viper.BindPFlags(testPrintCmd.Flags())
// }
