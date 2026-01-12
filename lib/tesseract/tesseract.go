package tesseract

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"proWeb/lib/files"
	"runtime"

	"github.com/fatih/color"
)

const (
	errorOC = "данная ОС в данный момент не поддерживается"
)

// Install - полная установка Tesseract в Program Files, tempDir - временная папка для хранения установщика
func Install(tempDir string) error {
	if IsTesseractInstalled() {
		return nil
	}

	installPath, err := downloadTesseractInstaller(tempDir)
	if err != nil {
		return err
	}
	defer os.Remove(installPath)

	if err := setupTesseract(installPath); err != nil {
		return err
	}

	return nil
}

// IsTesseractInstalled проверяет, установлен ли Tesseract
func IsTesseractInstalled() bool {
	if err := CheckTesseract(); err == nil {
		return true
	}

	candidates := []string{
		`C:\Program Files\Tesseract-OCR`,
		`C:\Program Files (x86)\Tesseract-OCR`,
	}

	for _, dir := range candidates {
		if _, err := os.Stat(filepath.Join(dir, "tesseract.exe")); err == nil {
			return true
		}
	}

	return false
}

// downloadTesseractInstaller - скачивает установщик Tesseract в необходимую директорию, возвращает путь к установщику
func downloadTesseractInstaller(dir string) (installerPath string, err error) {
	if runtime.GOOS == "windows" {
		url := "https://github.com/tesseract-ocr/tesseract/releases/download/5.5.0/tesseract-ocr-w64-setup-5.5.0.20241111.exe"
		target := filepath.Join(dir, "tesseract.exe")

		if err := files.DownloadFile(url, target); err != nil {
			return "", fmt.Errorf("не удалось скачать .exe файл: %v", err)
		}

		if err := os.Chmod(target, 0775); err != nil {
			return "", fmt.Errorf("не удалось сделать файл исполняемым")
		}

		return target, nil
	}
	//TODO: сделать для Linux
	return "", fmt.Errorf(errorOC)
}

// setupTesseract - запускает установщик Tesseract и скачивает в ProgramFiles
func setupTesseract(installerPath string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf(errorOC)
	}

	cmd := exec.Command("powershell",
		"-Command",
		"Start-Process",
		fmt.Sprintf(`"%s"`, installerPath),
		"-Verb", "runAs")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	color.Blue("Запуск установщика Tesseract")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ошибка запуска установщика")
	}
	return nil
}

// CheckTesseract - проверяет наличие Tesseract в PATH
func CheckTesseract() error {
	_, err := exec.LookPath("tesseract")

	return err
}
