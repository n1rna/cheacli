package cmd

import (
	"cheat/db"
	"cheat/term"
	"cheat/utils"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show [repo] [cheat]",
		Short: "Show the requested cheatsheet",
		Long:  `Show the cheatsheet if it exists in one of the repos.`,
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
		utils.CommandError(cmd, "The cheatsheet you asked does not exist :(", true)
		return
	}

	if len(commands) > 1 {
		utils.CommandError(cmd, "More than one command matched. Specify the repo", true)
		return
	}

	dir := getFileDirectory(commands[0].Repo, commands[0].FileName)
	term.RenderMarkdownFile(dir)
}

func showAvailableCheats() {
	fdb := db.GetDatabase("db", "")
	fdb.Read()

	commands := fdb.GetAllCommands()

	ch := make(chan string)
	go func() {
		for repo, cheats := range commands {
			repoLine := fmt.Sprintf("%s %s", utils.Yellow("[Repository]"), utils.Red(repo))
			maxLineLength := len(cheats[0].Name)
			for _, cheat := range cheats {
				if len(cheat.Name) > maxLineLength {
					maxLineLength = len(cheat.Name)
				}
			}
			ch <- repoLine
			for i := 0; i < len(cheats); i++ {
				color := utils.Blue
				if i%4 == 0 {
					color = utils.Green
				}
				var cheatLine string
				if i < len(cheats)-1 {
					padding := strings.Repeat(" ", maxLineLength-len(cheats[i].Name)+1)
					cheatLine = fmt.Sprintf("%s%s%s", color(cheats[i].Name), padding, color(cheats[i+1].Name))
					i++
				} else {
					cheatLine = fmt.Sprintf("%s", cheats[i].Name)
				}
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
