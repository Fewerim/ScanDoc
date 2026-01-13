package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newOpenLogCmd(a *AppCLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "open_log",
		Short:   "Открывает файл с логами",
		Example: "scandoc.exe openLog",
		RunE: func(cmd *cobra.Command, args []string) error {
			const op = ".open_log"

			clearFlag, err := cmd.Flags().GetBool("clear")
			if err != nil {
				return fmt.Errorf("ошибка чтения флага: %v", err)
			}
			return a.App.Logs(a.Name+op, clearFlag)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().BoolP("clear", "c", false, "очистка файла с логами")
	return cmd
}
