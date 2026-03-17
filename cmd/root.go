/*
Copyright © 2025 Erin Atkinson
*/
package cmd

import (
	"os"

	"github.com/erindatkinson/slack-emojinator/internal/utilities"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "slack-emojinator",
	Short: "A tool to bulk import and export slack emojis",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	initConfig()
	rootCmd.PersistentFlags().StringVarP(&browser, "browser", "b", utilities.ConfigOrEnv("slack", "browser"), "browser to look for token")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", utilities.ConfigOrEnv("slack", "profile"), "profile to look for token")
	rootCmd.PersistentFlags().StringVarP(&subdomain, "subdomain", "s", utilities.ConfigOrEnv("slack", "subdomain"), "what subdomain to pull a slack token for")
}

// initConfig reads in config file
func initConfig() {
	viper.SetConfigName(".config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/slack-emojinator/")
	viper.AddConfigPath("$HOME/.emojinator")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
}
