package cmd

import (
	"github.com/spf13/cobra"
)

// newRunOnceCmd - обертка над runOnce, чтобы можно было использовать логгер и конфиг
func newRunOnceCmd(a *AppCLI) *cobra.Command {
	var pathToFile, nameNewFile string

	cmd := &cobra.Command{
		Use:     "run_once",
		Short:   "Команда для обработки одного файла: run_once -f='путь к файлу' -n='название будущего файла'",
		Example: "scandoc.exe run_once --file='.test/scan.jpg' --name='result'\nотправит на обработку файл 'scan.jpg' и сохранит результат под именем 'result.json'",
		RunE: func(cmd *cobra.Command, args []string) error {
			const op = ".run_once"

			return a.AppCmds.OnceFile(a.Name+op, pathToFile, nameNewFile)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.MarkFlagRequired("file")
	cmd.MarkFlagRequired("name")

	cmd.Flags().StringVarP(&pathToFile, "file", "f", "", "путь к файлу, который требуется обработать")
	cmd.Flags().StringVarP(&nameNewFile, "name", "n", "", "имя нового файла")

	return cmd
}
