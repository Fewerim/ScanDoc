package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "proweb-cli",
	Short: "CLI для распознавания бухгалтерских документов",
}

func Execute() error {
	return rootCmd.Execute()
}
