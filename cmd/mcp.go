package cmd

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/gingray/swisstools/pkg/common"
	"github.com/gingray/swisstools/pkg/mcp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setupMCP = `
{
  "mcpServers": {
    "swisstools-mcp": {
      "url": "http://localhost:8181/"
    }
  }
}
`
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "mcp server for exposed operations",
	Long:  `mcp server for exposed operations`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := common.Config{}
		err := viper.Unmarshal(&cfg)
		if err != nil {
			fmt.Println(setupMCP)
			log.Error(err)
			return
		}
		server := mcp.NewMCPServer(&cfg)
		err = server.Run()
		if err != nil {
			log.Error(err)
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jiraCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jiraCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
