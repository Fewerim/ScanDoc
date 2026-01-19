package appCmds

import (
	"proWeb/internal/appUtils"
	"proWeb/internal/tesseract"
	"strings"
)

// InstallTesseract - устанавливает Tesseract
func (a *App) InstallTesseract() error {
	if err := tesseract.CheckTesseract(); err == nil {
		appUtils.InfoMessage("Tesseract успешно установлен и добавлен в PATH")
		return nil
	}

	if err := appUtils.InstallTesseract(); err != nil {
		return err
	}

	if !tesseract.IsTesseractInstalled() {
		appUtils.InfoMessage("Завершите мастер установки и повторите попытку")
		return nil
	}

	//TODO: linux
	info := strings.Builder{}
	info.WriteString("ВНИМАНИЕ: для корректной работы приложения, установите tesseract в PATH\n")
	info.WriteString("windows: можете попробовать команду 'setx /M PATH \"<путь к каталогу с tesseract>;%PATH%\"\n'")
	info.WriteString("--примечания:\n\tдля успешного добавления в PATH требуются права администратора\n\tпосле добавления в PATH перезапустите консоль")

	appUtils.InfoMessage(info.String())
	return nil
}
