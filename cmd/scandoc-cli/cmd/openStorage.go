package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newOpenStorageCmd(a *AppCLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "open_storage",
		Short:   "Открывает локальное хранилище",
		Example: "scandoc.exe open_storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			const op = ".open_storage"

			clearFlag, err := cmd.Flags().GetBool("clear")
			if err != nil {
				return fmt.Errorf("ошибка чтения флага: %v", err)
			}
			return a.App.OpenStorage(a.Name+op, clearFlag)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().BoolP("clear", "c", false, "очистка папки storage")
	return cmd
}
