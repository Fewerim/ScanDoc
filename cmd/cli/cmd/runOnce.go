package cmd

import (
	"fmt"
	cliUtils "proWeb/internal/cliUtils"
	"proWeb/internal/files"

	"github.com/spf13/cobra"
)

const port = 3210

// onceFile - функция, которая проверяет входные данные, создает подключение к серверу,
// обрабатывает и сохраняет файл локально
func onceFile(cmd *cobra.Command, args []string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = cliUtils.InternalError(fmt.Sprintf("внутренняя ошибка: %v", r))
		}
	}()

	filePath := args[0]
	createdFileName := args[1]

	if !files.StorageExists() {
		if err = files.CreateStorageJSON(); err != nil {
			return cliUtils.InternalError("ошибка создания локального хранилища")
		}
	}

	if err = cliUtils.ValidateExtensionFile(filePath); err != nil {
		return cliUtils.UserError(err.Error())
	}

	if err = cliUtils.CheckExistsFile(filePath); err != nil {
		return cliUtils.UserError(err.Error())
	}

	//TODO: запуск сервиса и ожидание ответа
	_, err = cliUtils.ProcessOnceFile(filePath, createdFileName, port)
	if err != nil {
		return cliUtils.ServerError(err.Error())
	}

	fmt.Println(cliUtils.Success("программа завершилась успешно").ToString())
	return nil
}

var runOnce = &cobra.Command{
	Use:           "run_once",
	Short:         "Команда для обработки одного файла: runOnce [путь к файлу] [название будущего файла]",
	Args:          cobra.ExactArgs(2),
	RunE:          onceFile,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	rootCmd.AddCommand(runOnce)
}
