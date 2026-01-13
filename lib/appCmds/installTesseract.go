package appCmds

import (
	"proWeb/lib/appUtils"
	"proWeb/lib/tesseract"
	"strings"

	"github.com/fatih/color"
)

// InstallTesseract - устанавливает Tesseract
func (a *App) InstallTesseract() error {
	if err := tesseract.CheckTesseract(); err == nil {
		color.Blue("Tesseract успешно установлен и добавлен в PATH")
		return nil
	}

	if err := appUtils.InstallTesseract(); err != nil {
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
