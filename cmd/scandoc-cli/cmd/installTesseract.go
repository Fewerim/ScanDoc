package cmd

import (
	"github.com/spf13/cobra"
)

func newInstallTesseractCmd(a *AppCLI) *cobra.Command {
	return &cobra.Command{
		Use:     "install_tesseract",
		Short:   "Запускает установщик tesseract",
		Example: "scandoc.exe install_tesseract\nзапускает установщик tesseract",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.AppCmds.InstallTesseract()
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
}
