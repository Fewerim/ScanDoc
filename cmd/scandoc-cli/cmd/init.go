package cmd

import (
	"github.com/spf13/cobra"
)

func newInitAppCmd(a *AppCLI) *cobra.Command {
	return &cobra.Command{
		Use:     "init",
		Short:   "Устанавливает необходимые зависимости для корректной работы приложения",
		Example: "scandoc.exe init\nустановит необходимые зависимости и локальное хранилище",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			const op = ".init"
			return a.AppCmds.InitApp(a.Name + op)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
}
