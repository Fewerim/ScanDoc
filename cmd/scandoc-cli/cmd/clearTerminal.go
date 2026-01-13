package cmd

import (
	"github.com/spf13/cobra"
)

func newClearCmd(a *AppCLI) *cobra.Command {
	return &cobra.Command{
		Use:     "clear_terminal",
		Short:   "Очистить терминал",
		Long:    "Очищает экран терминала (работает в CMD, PowerShell, bash, zsh)",
		Example: "scandoc.exe clear\nочистит экран терминала",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			const op = ".clear_terminal"
			return a.App.ClearTerminal(a.Name + op)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
}
