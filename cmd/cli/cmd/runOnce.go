package cmd

import (
	"fmt"
	"path/filepath"
	cliUtils "proWeb/internal/cliUtils"
	"proWeb/internal/cliUtils/cliWorks"
	"proWeb/internal/exitCodes"
	"proWeb/internal/tesseract"
	"time"

	"github.com/spf13/cobra"
)

// onceFile - проверяет входные данные, создает подключение к серверу,
// обрабатывает и сохраняет файл локально
func (a *App) onceFile(filePath, createdFileName string) (err error) {
	const operation = "cli.onceFile"
	if err := tesseract.CheckTesseract(); err != nil {
		a.Log.Error(operation, "tesseract не добавлен в PATH", exitCodes.UserError)
		return cliUtils.UserError("tesseract не добавлен в PATH")
	}

	if err := cliUtils.CheckIsAutorunCorrect(); err != nil {
		a.Log.Error(operation, err.Error(), cliUtils.GetExitCode(err, exitCodes.UserError))
		return err
	}

	a.Log.Info(operation, "Команда начала свое выполнение")

	start := time.Now()

	defer func() {
		if r := recover(); r != nil {
			err = cliUtils.InternalError(fmt.Sprintf("внутренняя ошибка: %v", r))
			a.Log.Error(operation, "операция завершена с паникой", exitCodes.InternalError)
		}
	}()

	if filePath == "" {
		a.Log.Error(operation, "не указан путь к файлу", exitCodes.UserError)
		return cliUtils.UserError("укажите путь к файлу")
	}

	if createdFileName == "" {
		a.Log.Error(operation, "не указано имя нового файла", exitCodes.UserError)
		return cliUtils.UserError("не указано имя нового файла")
	}

	if err := a.CheckStorageJSON(); err != nil {
		a.Log.Error(operation, err.Error(), cliUtils.GetExitCode(err, exitCodes.InternalError))
		return cliUtils.InternalError(err.Error())
	}

	if err = cliUtils.ValidateExtensionFile(filePath); err != nil {
		a.Log.Error(operation, err.Error(), cliUtils.GetExitCode(err, exitCodes.UserError), filepath.Base(filePath))
		return err
	}

	if err = cliUtils.CheckExistsFile(filePath); err != nil {
		a.Log.Error(operation, err.Error(), cliUtils.GetExitCode(err, exitCodes.UserError))
		return err
	}

	a.Log.Info(operation, "начало обработки файла")

	result, err := cliWorks.ProcessOnceFile(filePath, createdFileName, a.Cfg)
	if err != nil {
		a.Log.Error(operation, err.Error(), cliUtils.GetExitCode(err, exitCodes.ServerError), filepath.Base(filePath))
		return err
	}
	elapsed := time.Since(start)
	result.SetElapsedTime(elapsed)

	cliUtils.NewSuccess(&result).PrintSuccess()
	a.Log.Info(operation, fmt.Sprintf("операция завершена, время выполнения: %.3fs", result.GetElapsedTime()))
	return nil
}

// newRunOnceCmd - обертка над runOnce, чтобы можно было использовать логгер и конфиг
func newRunOnceCmd(a *App) *cobra.Command {
	var pathToFile, nameNewFile string

	cmd := &cobra.Command{
		Use:     "run_once",
		Short:   "Команда для обработки одного файла: run_once -f='путь к файлу' -n='название будущего файла'",
		Example: "scanner.exe run_once --file='.test/scan.jpg' --name='result'\nотправит на обработку файл 'scan.jpg' и сохранит результат под именем 'result.json'",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.onceFile(pathToFile, nameNewFile)
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
