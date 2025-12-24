package cmd

import (
	"fmt"
	"log"
	"proWeb/internal/cliUtils"
	"proWeb/internal/cliUtils/cliWorks"
	"proWeb/internal/files"
	"time"

	"github.com/spf13/cobra"
)

// multiFiles - проверяет входные данные, создает подключение к серверу,
// обрабатывает и сохраняет файлы локально
func (a *App) multiFiles(cmd *cobra.Command, args []string) (err error) {
	const operation = "cli.multiFiles"

	a.Log.Info(operation, "начало обработки директории файлов")

	start := time.Now()

	defer func() {
		if r := recover(); r != nil {
			err = cliUtils.InternalError(fmt.Sprintf("внутренняя ошибка: %v", r))
			a.Log.Error(operation, "операция завершена с паникой", 3)
		}
	}()

	dirPath := args[0]
	log.Println("система одновременно может обрабатывать не более 5 файлов, если файлов больше, то на это требуется больше времени")

	files.InitStorage(a.Cfg.StoragePath)

	if !files.StorageExists() {
		a.Log.Info(operation, "локальное хранилище не существует, создание нового")
		if err = files.CreateStorageJSON(); err != nil {
			a.Log.Error(operation, fmt.Sprintf("ошибка создания локального хранилища: %v", err), 3)
			return cliUtils.InternalError("ошибка создания локального хранилища")
		}
		a.Log.Info(operation, "локальное хранилище успешно создано")
	}

	if err = cliUtils.CheckExistsPath(dirPath); err != nil {
		a.Log.Error(operation, err.Error(), 1)
		return err
	}

	result, err, errs := cliWorks.MultiProcessFiles(dirPath, a.Cfg)
	if err != nil {
		a.Log.Error(operation, err.Error(), 2)
		return err
	}

	if errs != nil {
		for _, fileErr := range errs {
			a.Log.Error(operation, fileErr.Error(), 2)
		}
	}
	elapsed := time.Since(start)
	result.SetElapsedTime(elapsed)

	cliUtils.NewSuccess(&result).PrintSuccess()
	a.Log.Info(operation, fmt.Sprintf("операция завершена, время выполнения: %.3fs", result.Elapsed))
	return nil
}

func newMultiRunCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:           "run_multi",
		Short:         "Команда для обработки всех файлов в директории: run_multi [директория]",
		Args:          cobra.ExactArgs(1),
		RunE:          a.multiFiles,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
}
