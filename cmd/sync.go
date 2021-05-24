package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RepoManifest struct {
	Name    string
	Version string
}

func NewSyncCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "sync repos",
		Long:  `sync repos`,
		Run:   runSyncCommand,
	}
}

func runSyncCommand(cmd *cobra.Command, args []string) {
	reposDir := viper.GetString("repos")
	commandsDbDir := viper.GetString("db")

	commandsDB, err := os.OpenFile(fmt.Sprintf("%s/%s", commandsDbDir, "db.update"), os.O_RDWR|os.O_CREATE, 0755)
	commandsDB.Truncate(0)

	if err != nil {
		fmt.Println(err)
	}
	defer commandsDB.Close()

	repos, err := ioutil.ReadDir(reposDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, repo := range repos {
		if !repo.IsDir() {
			continue
		}
		manifestJson, err := os.Open(fmt.Sprintf("%s/%s/%s", reposDir, repo.Name(), "manifest.json"))
		if err != nil {
			fmt.Println(err)
		}
		defer manifestJson.Close()

		byteValue, _ := ioutil.ReadAll(manifestJson)
		var manifest RepoManifest

		json.Unmarshal(byteValue, &manifest)

		cheatsDir := fmt.Sprintf("%s/%s/%s", reposDir, repo.Name(), "cheats")

		cheats, err := ioutil.ReadDir(cheatsDir)
		if err != nil {
			log.Fatal(err)
		}

		for _, cheat := range cheats {
			commandName := strings.Replace(cheat.Name(), ".md", "", -1)
			commandsDB.Write(
				[]byte(fmt.Sprintf("[%s] %s:%s:%s\n", repo.Name(), commandName, cheat.Name(), "2")),
			)
		}
	}

	cleanUpAfterSync()
}

func cleanUpAfterSync() {
	commandsDbDir := viper.GetString("db")

	original := fmt.Sprintf("%s/%s", commandsDbDir, "db")
	updated := fmt.Sprintf("%s/%s", commandsDbDir, "db.update")

	os.Rename(updated, original)
}

func init() {
	var syncCommand = NewSyncCommand()
	syncCommand.Flags().StringP("repo", "r", "", "Sync only specific repository")

	RootCmd.AddCommand(syncCommand)
}
