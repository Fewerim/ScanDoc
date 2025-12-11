package cmd

import (
	"fmt"

	"github.com/iancoleman/orderedmap"
	"github.com/spf13/cobra"
)

const (
	CommandHelp    = "\\bin\\proweb.exe help"
	CommandRunOnce = "\\bin\\proweb.exe run_once [путь к файлу] [название нового файла]"
)

var commandsHelp = orderedmap.New()

func initCommands() {
	commandsHelp.Set("-Показывает список всех команд", CommandHelp)
	commandsHelp.Set("-Обрабатывает один документ", CommandRunOnce)
}

func helper(cmd *cobra.Command, args []string) {
	initCommands()
	fmt.Println("Добро пожаловать в сервис распознавания бухгалтерских документов")
	fmt.Println("Вот список доступных команд (*использовать без скобок):")
	keys := commandsHelp.Keys()
	for _, k := range keys {
		v, _ := commandsHelp.Get(k)
		fmt.Printf("%s\t%s\n", k, v.(string))
	}
	return
}

var help = &cobra.Command{
	Use:   "help",
	Short: "Список команд",
	Args:  cobra.NoArgs,
	Run:   helper,
}

func init() {
	rootCmd.AddCommand(help)
}
