package cmd

import (
	"fmt"
	"os"

	"github.com/ldez/go-git-cmd-wrapper/clone"
	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "init cli",
		Long:  `init cli`,
		Run:   runInitCommand,
	}
}

func runInitCommand(cmd *cobra.Command, args []string) {
	// Clone default repositories
	reposDir := viper.GetString("repos")
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Something bad happened!")
		return
	}
	os.Chdir(reposDir)
	repositories := viper.GetStringMap("repositories")
	for _, repo := range repositories {
		output, _ := git.Clone(clone.Repository(repo.(map[string]interface{})["origin"].(string)), clone.Progress)
		fmt.Println(output)
	}
	os.Chdir(cwd)

	// Run sync command once
	var syncCommand = NewSyncCommand()
	syncCommand.Run(cmd, []string{})
}

func init() {
	var initCommand = NewInitCommand()

	RootCmd.AddCommand(initCommand)
}
