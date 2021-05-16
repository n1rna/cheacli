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

type Cheat struct {
	Name     string
	FileName string
	Columns  string
}

var showCommand = &cobra.Command{
	Use:   "show [repo] [cheat]",
	Short: "show cheatsheet",
	Long:  `Show the cheatsheet if it exists.`,
	Run:   runShowCommand,
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
	db, err := readCommandsDB()

	if err != nil {
		fmt.Println("Database could not be found. Please try syncing again.\n")
		RootCmd.Help()
		return
	}

	var re *regexp.Regexp

	if len(args) > 1 {
		re = regexp.MustCompile(fmt.Sprintf(`\[(%s)\]\s(%s)\:(.*\.md)\:(.*)`, args[0], args[1]))
	} else {
		re = regexp.MustCompile(fmt.Sprintf(`\[(.*)\]\s(%s)\:(.*\.md)\:(.*)`, args[0]))
	}

	matches := re.FindAllStringSubmatch(string(db), -1)

	if len(matches) < 1 {
		fmt.Println("The cheatsheet you asked does not exist :(\n")
		cmd.Help()
		return
	}

	if len(matches) > 1 {
		fmt.Println("More than one command matched. Specify the repo.\n")
		cmd.Help()
		return
	}

	dir := getFileDirectory(matches[0][1], matches[0][3])
	term.RenderMarkdownFile(dir)
}

func showAvailableCheats() {
	db, err := readCommandsDB()

	if err != nil {
		fmt.Println("Database could not be found. Please try syncing again.\n")
		RootCmd.Help()
		return
	}

	re := regexp.MustCompile(`\[(.*)\]\s(.*)\:(.*\.md)\:(.*)`)

	commands := make(map[string][]*Cheat)

	matches := re.FindAllStringSubmatch(string(db), -1)
	for _, v := range matches {
		commands[v[1]] = append(commands[v[1]], &Cheat{
			Name:     v[2],
			FileName: v[3],
			Columns:  v[4],
		})
	}

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

func getFileDirectory(repoName string, fileName string) string {
	reposDir := viper.GetString("repos")
	return fmt.Sprintf("%s/%s/cheats/%s", reposDir, repoName, fileName)
}

func readCommandsDB() (db []byte, err error) {
	dbDirectory := viper.GetString("db")

	db, err = os.ReadFile(fmt.Sprintf("%s/%s", dbDirectory, "db"))
	return
}

func init() {
	showCommand.PersistentFlags().BoolP("list", "l", false, "List all subcommands")
	showCommand.PersistentFlags().Bool("repositories", false, "List all subcommands")
	showCommand.PersistentFlags().BoolP("repo", "r", false, "Specify the repository")

	RootCmd.AddCommand(showCommand)
}
