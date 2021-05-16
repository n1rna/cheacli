package cmd

import (
	"cheatcli/term"
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var showCommand = &cobra.Command{
	Use:   "show [cheat]",
	Short: "show cheat",
	Long:  `show cheats`,
	Run: func(cmd *cobra.Command, args []string) {
		if isList, _ := cmd.PersistentFlags().GetBool("list"); isList {
			showAvailableCheats()
			return
		}

		if len(args) < 1 {
			cmd.Help()
			return
		}

		findAndShowCheat(args)
	},
}

func findAndShowCheat(args []string) {
	db := readCommandsDB()

	var re *regexp.Regexp

	if len(args) > 1 {
		re = regexp.MustCompile(fmt.Sprintf(`\[(%s)\]\s(%s)\:(.*\.md)\:(.*)`, args[0], args[1]))
	} else {
		re = regexp.MustCompile(fmt.Sprintf(`\[(.*)\]\s(%s)\:(.*\.md)\:(.*)`, args[0]))
	}

	matches := re.FindAllStringSubmatch(string(db), -1)

	if len(matches) < 1 {
		panic("The command you asked does not exist")
	}

	if len(matches) > 1 {
		panic("More than one command matched. Specify the repo.")
	}

	dir := getFileDirectory(matches[0][1], matches[0][3])
	term.ShowMDFile(dir)
}

func getFileDirectory(repoName string, fileName string) string {
	reposDir := viper.GetString("reposDir")
	return fmt.Sprintf("%s/%s/cheats/%s", reposDir, repoName, fileName)
}

type Cheat struct {
	Name     string
	FileName string
	Columns  string
}

func showAvailableCheats() {
	db := readCommandsDB()

	commandPattern := regexp.MustCompile(`\[(.*)\]\s(.*)\:(.*\.md)\:(.*)`)

	commands := make(map[string][]*Cheat)

	matches := commandPattern.FindAllStringSubmatch(string(db), -1)
	for _, v := range matches {
		commands[v[1]] = append(commands[v[1]], &Cheat{
			Name:     v[2],
			FileName: v[3],
			Columns:  v[4],
		})
	}
	fmt.Println(commands)

	ch := make(chan string)

	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

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

func readCommandsDB() []byte {
	commandsDbDir := viper.GetString("commandsDbDir")

	commandsDB, err := os.ReadFile(fmt.Sprintf("%s/%s", commandsDbDir, "commandsdb.update"))
	if err != nil {
		panic("Database file does not exist. Please sync first.")
	}
	return commandsDB
}

func init() {
	showCommand.PersistentFlags().Bool("list", false, "List all subcommands")
	showCommand.PersistentFlags().Bool("repo", false, "Specify the repository")

	RootCmd.AddCommand(showCommand)
}
