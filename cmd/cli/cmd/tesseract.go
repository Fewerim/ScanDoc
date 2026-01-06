package cmd

import (
	"proWeb/internal/cliUtils"
	"proWeb/internal/tesseract"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func (a *App) installTesseract(cmd *cobra.Command, args []string) error {
	if err := tesseract.CheckTesseract(); err == nil {
		color.Blue("Tesseract успешно установлен и добавлен в PATH")
		return nil
	}

	if err := cliUtils.InstallTesseract(); err != nil {
		return err
	}

	if !tesseract.IsTesseractInstalled() {
		color.Blue("Завершите мастер установки и повторите попытку")
		return nil
	}

	//TODO: linux
	info := strings.Builder{}
	info.WriteString("ВНИМАНИЕ: для корректной работы приложения, установите tesseract в PATH\n")
	info.WriteString("windows: можете попробовать команду 'setx /M PATH \"<путь к каталогу с tesseract>;%PATH%\"\n'")
	info.WriteString("--примечания:\n\tдля успешного добавления в PATH требуются права администратора\n\tпосле добавления в PATH перезапустите консоль")

	color.Blue(info.String())
	return nil
}

func newInstallTesseract(a *App) *cobra.Command {
	return &cobra.Command{
		Use:           "install_tesseract",
		Short:         "Запускает установщик tesseract",
		Example:       "scanner.exe install_tesseract\nзапускает установщик tesseract",
		Args:          cobra.NoArgs,
		RunE:          a.installTesseract,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
}
