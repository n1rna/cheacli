package cmd

import (
	"cheat/utils"
	"fmt"
	"os"
	"os/exec"

	"github.com/ldez/go-git-cmd-wrapper/clone"
	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize cheatcli",
		Long:  `Install main cheats repo and initialize the cheats database`,
		Run:   runInitCommand,
	}
}

func runInitCommand(cmd *cobra.Command, args []string) {
	// Clone default repositories
	reposDir := viper.GetString("repos")
	cwd, err := os.Getwd()
	if err != nil {
		utils.CommandError(cmd, "Something bad happened!", true)
		return
	}
	os.Chdir(reposDir)
	repositories := viper.GetStringMap("repositories")
	for _, repo := range repositories {
		git.Clone(
			clone.Repository(repo.(map[string]interface{})["origin"].(string)),
			clone.Progress,
			git.CmdExecutor(streamingGitExecutor),
		)
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

func streamingGitExecutor(name string, debug bool, args ...string) (string, error) {

	cmd := exec.Command(name, args...)

	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return "", err
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}

	cmd.Wait()

	return "", nil
}
