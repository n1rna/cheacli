package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "cheatcli",
	Short: "Cheat sheets in your console",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.SetDefault("repos", fmt.Sprintf("%s/.cheats/%s", home, "repos"))
	viper.SetDefault("db", fmt.Sprintf("%s/.cheats", home))

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.config/cheat/.cheatclirc)")
}

func initConfig() {
	viper.SetConfigType("json")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".config/cheat/.cheatclirc")
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("\nUsing config file:", viper.ConfigFileUsed())
	}
}
