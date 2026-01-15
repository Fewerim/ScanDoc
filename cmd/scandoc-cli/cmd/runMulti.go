package cmd

import (
	"github.com/spf13/cobra"
)

func newMultiRunCmd(a *AppCLI) *cobra.Command {
	var directory string
	var createdFolderName string

	cmd := &cobra.Command{
		Use:     "run_multi",
		Short:   "Команда для обработки всех файлов в директории: run_multi -d='директория' -n='название папки'",
		Example: "scandoc.exe run_multi --dir='./packageToScan' --name='test'\nотправит пакет файлов на обработку, результаты будут сохранены в подпапку с именем, переданным во флаге 'name', в локальное хранилище под теми же именами",
		RunE: func(cmd *cobra.Command, args []string) error {
			const op = ".run_multi"
			return a.AppCmds.MultiFiles(a.Name+op, directory, createdFolderName)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().StringVarP(&directory, "dir", "d", "", "путь к директории, которую требуется обработать")
	cmd.Flags().StringVarP(&createdFolderName, "name", "n", "", "имя папки, в которую будут сохранены обработанные файлы")

	cmd.MarkFlagRequired("dir")
	cmd.MarkFlagRequired("name")

	return cmd
}
