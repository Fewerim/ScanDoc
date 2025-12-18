package cmd

import (
	"fmt"
	cliUtils "proWeb/internal/cliUtils"
	"proWeb/internal/cliUtils/cliWorks"
	"proWeb/internal/files"

	"github.com/spf13/cobra"
)

// onceFile - проверяет входные данные, создает подключение к серверу,
// обрабатывает и сохраняет файл локально
func (a *App) onceFile(cmd *cobra.Command, args []string) (err error) {
	const operation = "cli.onceFile"
	a.Log.Info(operation, "начало обработки файла")

	defer func() {
		if r := recover(); r != nil {
			err = cliUtils.InternalError(fmt.Sprintf("внутренняя ошибка: %v", r))
			a.Log.Error(operation, "операция завершена с паникой", 3)
		}
	}()

	filePath := args[0]
	createdFileName := args[1]

	files.InitStorage(a.Cfg.StoragePath)

	if !files.StorageExists() {
		a.Log.Info(operation, "локальное хранилище не существует, создание нового")
		if err = files.CreateStorageJSON(); err != nil {
			a.Log.Error(operation, fmt.Sprintf("ошибка создания локального хранилища: %v", err), 3)
			return cliUtils.InternalError("ошибка создания локального хранилища")
		}
		a.Log.Info(operation, "локальное хранилище успешно создано")
	}

	if err = cliUtils.ValidateExtensionFile(filePath); err != nil {
		a.Log.Error(operation, err.Error(), 1)
		return err
	}

	if err = cliUtils.CheckExistsFile(filePath); err != nil {
		a.Log.Error(operation, err.Error(), 1)
		return err
	}

	_, err = cliWorks.ProcessOnceFile(filePath, createdFileName, a.Cfg)
	if err != nil {
		a.Log.Error(operation, err.Error(), 2)
		return err
	}

	//TODO: формирование ответа для пользователя

	fmt.Println(cliUtils.Success("программа завершилась успешно").ToString())
	a.Log.Info(operation, "операция завершена успешно")
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
