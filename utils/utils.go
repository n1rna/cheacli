package utils

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	Green = color.New(color.FgGreen).SprintFunc()
	Red   = color.New(color.FgRed).SprintFunc()
)

func CommandError(cmd *cobra.Command, err string, showHelp bool) {
	fmt.Println(Red("[Error]"), err)
	if showHelp {
		cmd.Help()
	}
}
