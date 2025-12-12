package cmd

import (
	"fmt"
	cliUtils "proWeb/internal/cliUtils"
	"proWeb/internal/cliUtils/cliWorks"
	"proWeb/internal/files"

	"github.com/spf13/cobra"
)

// onceFile - функция, которая проверяет входные данные, создает подключение к серверу,
// обрабатывает и сохраняет файл локально
func (a *App) onceFile(cmd *cobra.Command, args []string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = cliUtils.InternalError(fmt.Sprintf("внутренняя ошибка: %v", r))
		}
	}()

	filePath := args[0]
	createdFileName := args[1]

	files.InitStorage(a.Cfg.StoragePath)

	if !files.StorageExists() {
		if err = files.CreateStorageJSON(); err != nil {
			return cliUtils.InternalError("ошибка создания локального хранилища")
		}
	}

	if err = cliUtils.ValidateExtensionFile(filePath); err != nil {
		return err
	}

	if err = cliUtils.CheckExistsFile(filePath); err != nil {
		return err
	}

	_, err = cliWorks.ProcessOnceFile(filePath, createdFileName, a.Cfg)
	if err != nil {
		return err
	}

	//TODO: формирование ответа для пользователя

	fmt.Println(cliUtils.Success("программа завершилась успешно").ToString())
	return nil
}

// newRunOnceCmd - обертка над runOnce, чтобы можно было использовать логгер и конфиг
func newRunOnceCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:           "run_once",
		Short:         "Команда для обработки одного файла: run_once [путь к файлу] [название будущего файла]",
		Args:          cobra.ExactArgs(2),
		RunE:          a.onceFile,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
}
