package cmd

import (
	"cheat/db"
	"cheat/term"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show [repo] [cheat]",
		Short: "show cheatsheet",
		Long:  `Show the cheatsheet if it exists.`,
		Run:   runShowCommand,
	}
}

func runShowCommand(cmd *cobra.Command, args []string) {
	if isList, _ := cmd.PersistentFlags().GetBool("list"); isList {
		showAvailableCheats()
		return
	}

	if len(args) < 1 {
		cmd.Help()
		return
	}

	findAndShowCheat(cmd, args)
}

func findAndShowCheat(cmd *cobra.Command, args []string) {
	fdb := db.GetDatabase("db", "")
	fdb.Read()

	var commands []*db.CheatCommand

	if len(args) > 1 {
		commands = fdb.FindMatchedCheatsForRepo(args[0], args[1])
	} else {
		commands = fdb.FindMatchedCheats(args[0])
	}

	if len(commands) < 1 {
		fmt.Println("The cheatsheet you asked does not exist :(")
		cmd.Help()
		return
	}

	if len(commands) > 1 {
		fmt.Println("More than one command matched. Specify the repo")
		cmd.Help()
		return
	}

	dir := getFileDirectory(commands[0].Repo, commands[0].FileName)
	term.RenderMarkdownFile(dir)
}

func showAvailableCheats() {
	fdb := db.GetDatabase("db", "")
	fdb.Read()

	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	commands := fdb.GetAllCommands()

	ch := make(chan string)
	go func() {
		for repo, cheats := range commands {
			repoLine := fmt.Sprintf("%s - %s", yellow("REPO"), red(repo))
			ch <- repoLine
			for _, cheat := range cheats {
				cheatLine := fmt.Sprintf("%s - Fuck", cheat.Name)
				ch <- cheatLine
			}
		}
		close(ch)
	}()

	term.OutputStream(ch)
}

func getFileDirectory(repoName string, fileName string) string {
	reposDir := viper.GetString("repos")
	return fmt.Sprintf("%s/%s/cheats/%s", reposDir, repoName, fileName)
}

func init() {
	var showCommand = NewShowCommand()
	showCommand.PersistentFlags().BoolP("list", "l", false, "List all subcommands")
	showCommand.PersistentFlags().Bool("repositories", false, "List all subcommands")
	showCommand.PersistentFlags().BoolP("repo", "r", false, "Specify the repository")

	RootCmd.AddCommand(showCommand)
}
