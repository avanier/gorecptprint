package cmd

import (
	"github.com/avanier/gorecptprint/lib/comm"
	"github.com/avanier/gorecptprint/lib/print"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "print performs printy stuff",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		comm.Init()
		print.Print()
	},
}

func init() {
	var defaultSerialPort = "/dev/ttyS0"
	var defaultWriteChunkSize = uint(192)

	printCmd.PersistentFlags().String("port", defaultSerialPort, "serial port to use for communication")
	printCmd.PersistentFlags().Uint("write-chunk-size", defaultWriteChunkSize, "the size in bytes of every chunk of data sent to the printer")

	viper.BindPFlags(printCmd.PersistentFlags())
	viper.BindPFlags(printCmd.Flags())
}
