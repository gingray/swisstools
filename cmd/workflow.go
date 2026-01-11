package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/gingray/swisstools/pkg/common"
	"github.com/gingray/swisstools/pkg/workflow"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Trigger a workflow by apply request to particular endpoint with payload",
	Long:  `Trigger a workflow by apply request to particular endpoint with payload`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := common.Config{}
		err := viper.Unmarshal(&cfg)
		if err != nil {
			log.Error(err)
			return
		}
		if len(args) == 0 {
			for k, _ := range cfg.Workflows {
				log.Infof("available workflow: [%s]", k)
			}
			return
		}
		workflowRequest := workflow.NewWorkflowRequest(args[0], &cfg)
		err = workflowRequest.MakeRequest()
		if err != nil {
			log.Error(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(workflowCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jiraCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jiraCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
