/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/gingray/swisstools/pkg/common"
	"github.com/gingray/swisstools/pkg/jira"
	"github.com/gingray/swisstools/pkg/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// jiraCmd represents the jira command
var jiraCmd = &cobra.Command{
	Use:   "jira",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := common.Config{}
		err := viper.Unmarshal(&cfg)
		if err != nil {
			log.Error(err)
			return
		}
		jiraService := jira.NewJira(&cfg, ui.NewTableView())
		jiraService.GetIssues()

	},
}

func init() {
	rootCmd.AddCommand(jiraCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jiraCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jiraCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
