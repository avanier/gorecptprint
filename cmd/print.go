package cmd

import (
	"github.com/avanier/gorecptprint/cmd/lib/print"
	"github.com/spf13/cobra"
)

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "print performs printy stuff",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		print.Print()
	},
}
