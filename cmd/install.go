package cmd

import (
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
		Use:   "install",
		Short: "install repo",
		Long:  `install new repository`,
		Run:   runInstallCommand,
	}
}

func runInstallCommand(cmd *cobra.Command, args []string) {
	// Install new repository
	origin := args[0]
	reposDir := viper.GetString("repos")
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Something bad happened!", err)
		return
	}
	os.Chdir(reposDir)
	output, _ := git.Clone(clone.Repository(origin), clone.Progress)
	fmt.Println(output)
	os.Chdir(cwd)

	re := regexp.MustCompile(`.*/(.*)\.git`)
	repoName := re.Find([]byte(origin))
	err = validateInstallation(string(repoName))
	if err != nil {
		fmt.Println(err)
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

func init() {
	var installCommand = NewInstallCommand()

	RootCmd.AddCommand(installCommand)
}
