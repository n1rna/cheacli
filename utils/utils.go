package utils

import (
	"fmt"
	"os"
	"syscall"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

var (
	Green  = color.New(color.FgGreen).SprintFunc()
	Red    = color.New(color.FgRed).SprintFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
	Blue   = color.New(color.FgBlue).SprintFunc()
)

func CommandError(cmd *cobra.Command, err string, showHelp bool) {
	fmt.Println(Red("[Error]"), err)
	if showHelp {
		cmd.Help()
	}
}

func GetWinSize() (int, int, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return 0, 0, err
	}
	ws, err := unix.IoctlGetWinsize(int(tty.Fd()), syscall.TIOCGWINSZ)
	err2 := tty.Close()
	if err != nil {
		return 0, 0, err
	}
	return int(ws.Col), int(ws.Row), err2
}
