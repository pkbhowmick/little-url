package cmd

import (
	"fmt"
	"github.com/pkbhowmick/url-lite/apis"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(genCmd)
}

var genCmd = &cobra.Command{
	Use:   "gen-key",
	Long:  "Generate unique key",
	Short: "Generate unique key",
	RunE: func(cmd *cobra.Command, args []string) error {
		str, err := apis.GenKey()
		if err != nil {
			return err
		}
		fmt.Println(str)
		return nil
	},
}
