package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/iancoleman/orderedmap"
	"github.com/spf13/cobra"
)

const (
	CommandHelp     = "\\bin\\proweb.exe help"
	CommandInit     = "\\bin\\proweb.exe init"
	CommandRunOnce  = "\\bin\\proweb.exe run_once [путь к файлу] [название нового файла]"
	CommandRunMulty = "\\bin\\proweb.exe run_multi [директория]"
)

var commandsHelp = orderedmap.New()

func initCommands() {
	commandsHelp.Set("-Показывает список всех команд", CommandHelp)
	commandsHelp.Set("-Подтягивает все необходимые зависимости для работы приложения", CommandInit)
	commandsHelp.Set("-Обрабатывает один документ", CommandRunOnce)
	commandsHelp.Set("-Обрабатывает пакет документов", CommandRunMulty)
}

// helper - выводит список всех существующих команд
func (a *App) helper(cmd *cobra.Command, args []string) {
	const operation = "cli.helper"

	a.Log.Info(operation, "начало выполнения команды справки")

	initCommands()

	a.Log.Info(operation, "вывод приветственных сообщений и всех возможных команд")

	color.Blue("Добро пожаловать в сервис распознавания бухгалтерских документов")
	color.Blue("Вот список доступных команд (*использовать без скобок):")
	keys := commandsHelp.Keys()
	for _, k := range keys {
		v, _ := commandsHelp.Get(k)
		fmt.Printf("%s\t%s\n", v.(string), k)
	}

	a.Log.Info(operation, "справка успешно отображена")
	return
}

func newHelperCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:   "help",
		Short: "Список команд",
		Args:  cobra.NoArgs,
		Run:   a.helper,
	}
}
