package cmd

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/gingray/swisstools/pkg/api"
	"github.com/gingray/swisstools/pkg/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
)

// jiraCmd represents the jira command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Create an REST API for local machine",
	Long:  `Create an REST API for local machine`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := common.Config{}
		err := viper.Unmarshal(&cfg)
		if err != nil {
			log.Error(err)
			return
		}
		srv := api.CreateServer(&cfg)
		router := gin.Default()
		router.POST("/execute", srv.ExecuteCMD)
		router.GET("/health", srv.Health)

		log.Infof("Server starting on port %d", cfg.Api.Port)

		s := &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Api.Port),
			Handler: router,
		}
		err = s.ListenAndServe()
		if err != nil {
			log.Error(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jiraCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jiraCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
