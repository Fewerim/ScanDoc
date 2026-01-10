package cmd

import "github.com/spf13/cobra"

// newHelperCmd - меняет дефолтную справку на русский аналог
func newHelperCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:   "help",
		Short: "Справка по командам",
		Long:  "Показывает подробную справку по доступным командам",
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.Help()
		},
		Hidden: false,
	}
}
