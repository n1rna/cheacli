package cmd

import (
	"cheat/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/ldez/go-git-cmd-wrapper/clone"
	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewInstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "install [repositroy url]",
		Short: "Install a new repository",
		Long:  `Install new repositories from github by providing the git url`,
		Run:   runInstallCommand,
	}
}

func runInstallCommand(cmd *cobra.Command, args []string) {
	// Install new repository
	if len(args) == 0 {
		utils.CommandError(cmd, "Please provide a repository url", true)
		return
	}
	origin := args[0]
	reposDir := viper.GetString("repos")
	cwd, err := os.Getwd()
	if err != nil {
		utils.CommandError(cmd, fmt.Sprintf("Something bad happened! => %s", err), true)
		return
	}
	os.Chdir(reposDir)

	re := regexp.MustCompile(`.*\/(.*)\.git`)
	isValidUrl := re.Match([]byte(origin))
	repoName := re.FindAllStringSubmatch(origin, 1)[0][1]

	if !isValidUrl {
		utils.CommandError(cmd, "Please provide a valid repository url", true)
		return
	}

	git.Clone(clone.Repository(origin), clone.Progress, git.CmdExecutor(streamingGitExecutor))
	os.Chdir(cwd)

	err = validateInstallation(string(repoName))
	if err != nil {
		utils.CommandError(cmd, "The provided git repository is not a valid cheat repository", false)
		cleanupInstallation(repoName)
		return
	}
}

func validateInstallation(repoName string) (err error) {
	// Validate installation

	// Should have manifest and validate manifest
	reposDir := viper.GetString("repos")
	repoDirectory := fmt.Sprintf("%s/%s", reposDir, repoName)

	// Check directory exists
	if _, err := os.Stat(repoDirectory); os.IsNotExist(err) {
		return err
	}

	// Check manifest exists
	if _, err := os.Stat(fmt.Sprintf("%s/%s", repoDirectory, "manifest.json")); os.IsNotExist(err) {
		return err
	}

	// Check cheats directory exists
	if _, err := os.Stat(fmt.Sprintf("%s/%s", repoDirectory, "cheats")); os.IsNotExist(err) {
		return err
	}

	manifestJson, err := os.Open(fmt.Sprintf("%s/%s", repoDirectory, "manifest.json"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer manifestJson.Close()

	byteValue, _ := ioutil.ReadAll(manifestJson)
	var manifest RepoManifest

	if err := json.Unmarshal(byteValue, &manifest); err != nil {
		return err
	}

	return nil
}

func cleanupInstallation(repoName string) (err error) {
	reposDir := viper.GetString("repos")
	err = os.RemoveAll(fmt.Sprintf("%s/%s", reposDir, repoName))
	return err
}

func init() {
	var installCommand = NewInstallCommand()

	RootCmd.AddCommand(installCommand)
}
