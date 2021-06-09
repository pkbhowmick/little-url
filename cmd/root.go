package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:                        "url-lite",
	Short:                      "Root command for the url-lite project",
	Long:                       "Root command for the url-lite project",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Welcome!")
	},
}

func Execute()  {
	if err:= rootCmd.Execute(); err !=nil {
		os.Exit(1)
	}
}
