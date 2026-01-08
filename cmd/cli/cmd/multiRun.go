package cmd

import (
	"errors"
	"fmt"
	"proWeb/internal/cliUtils"
	"proWeb/internal/cliUtils/cliWorks"
	"proWeb/internal/exitCodes"
	"proWeb/internal/tesseract"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// multiFiles - проверяет входные данные, создает подключение к серверу,
// обрабатывает и сохраняет файлы локально
func (a *App) multiFiles(directory string) (err error) {
	const operation = "cli.multiFiles"

	color.Blue("Команда run_multi начала свое выполнение, валидация входных данных и проверка предварительных условий")

	if err := tesseract.CheckTesseract(); err != nil {
		a.Log.Error(operation, "tesseract не добавлен в PATH", exitCodes.UserError)
		return cliUtils.UserError("tesseract не добавлен в PATH")
	}

	if err := cliUtils.CheckIsAutorunCorrect(); err != nil {
		a.Log.Error(operation, err.Error(), cliUtils.GetExitCode(err, exitCodes.UserError))
		return err
	}

	start := time.Now()

	a.Log.Info(operation, "Команда начала свое выполнение")

	defer func() {
		if r := recover(); r != nil {
			err = cliUtils.InternalError(fmt.Sprintf("внутренняя ошибка: %v", r))
			a.Log.Error(operation, "операция завершена с паникой", exitCodes.InternalError)
		}
	}()

	if err := a.CheckStorageJSON(); err != nil {
		a.Log.Error(operation, err.Error(), cliUtils.GetExitCode(err, exitCodes.InternalError))
		return cliUtils.InternalError(err.Error())
	}

	if err = cliUtils.CheckExistsPath(directory); err != nil {
		a.Log.Error(operation, err.Error(), cliUtils.GetExitCode(err, exitCodes.UserError))
		return err
	}

	a.Log.Info(operation, "начало обработки директории файлов")
	color.Blue("Начало обработки директории файлов")

	result, err, errs := cliWorks.MultiProcessFiles(directory, a.Cfg)
	if err != nil {
		a.Log.Error(operation, err.Error(), cliUtils.GetExitCode(err, exitCodes.ServerError))
		return err
	}

	if errs != nil {
		for _, fileErr := range errs {
			a.Log.Error(operation, fileErr.Err.Error(), cliUtils.GetExitCode(err, exitCodes.ServerError), fileErr.FileName)
		}
	}
	elapsed := time.Since(start)
	result.SetElapsedTime(elapsed)

	if len(result.Results) == 0 {
		err = errors.New("в директории нет подходящих файлов для обработки")
		a.Log.Error(operation, err.Error(), cliUtils.GetExitCode(err, exitCodes.InternalError))
		return err
	}

	cliUtils.NewSuccess(&result).PrintSuccess()
	if errs != nil {
		cliUtils.FilesNotProcessed(errs)
	}

	a.Log.Info(operation, fmt.Sprintf("операция завершена, время выполнения: %.3fs", result.GetElapsedTime()))
	return nil
}

func newMultiRunCmd(a *App) *cobra.Command {
	var directory string

	cmd := &cobra.Command{
		Use:     "run_multi",
		Short:   "Команда для обработки всех файлов в директории: run_multi -d='директория'",
		Example: "scandoc.exe run_multi --dir='./packageToScan'\nотправит пакет файлов на обработку, результаты будут сохранены в локальное хранилище под теми же именами",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.multiFiles(directory)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().StringVarP(&directory, "dir", "d", "", "путь к директории, которую требуется обработать")
	cmd.MarkFlagRequired("dir")

	return cmd
}
