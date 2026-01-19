package appUtils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"proWeb/internal/appUtils/command"
	"proWeb/internal/tesseract"
	"syscall"
)

// CreateVenv - создает виртуальное окружения для python
func CreateVenv(pathToCreate string) error {
	venvPath := filepath.Join(pathToCreate, ".venv")

	pipPath := filepath.Join(venvPath, "Scripts", "pip.exe")
	if _, err := os.Stat(pipPath); err == nil {
		return nil
	}

	py, err := findSystemPython()
	if err != nil {
		return InternalError(err.Error())
	}

	cmd := command.Command(py, "-m", "venv", venvPath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = pathToCreate

	if err = cmd.Run(); err != nil {
		log.Printf(err.Error())
		info := fmt.Sprintf("ошибка создания venv: %v", err)
		return InternalError(info)
	}

	if _, err := os.Stat(pipPath); os.IsNotExist(err) {
		venvPython := filepath.Join(pathToCreate, ".venv", "Scripts", "python.exe")

		fixCmd := command.Command(venvPython, "-m", "ensurepip", "--upgrade", "--default-pip")
		fixCmd.Dir = pathToCreate
		fixCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

		var out bytes.Buffer
		fixCmd.Stdout = &out

		if err = fixCmd.Run(); err != nil {
			panic("не удалось установить pip")
		}
	}

	return nil
}

// InstallRequirements - устанавливает все необходимые зависимости для python скрипта
func InstallRequirements(pathToVenv, pathToScript string) error {
	reqsPath := filepath.Join(pathToScript, "requirements.txt")
	if _, err := os.Stat(reqsPath); os.IsNotExist(err) {
		return UserError("файл requirements.txt не найден")
	}

	pipPath := filepath.Join(pathToVenv, ".venv", "Scripts", "pip.exe")
	if _, err := os.Stat(pipPath); os.IsNotExist(err) {
		info := fmt.Sprintf("pip не найден в .venv: %s", pipPath)
		return UserError(info)
	}

	if !checkInternetConnection() {
		return UserError("не удается подключиться к PyPi. проверьте подключение к интернету")
	}

	cmd := command.Command(pipPath, "install", "-r", reqsPath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return InternalError(fmt.Sprint("зависимости не установлены"))
	}
	return nil
}

// InstallTesseract - устанавливает Tesseract на локальную машину пользователя
func InstallTesseract() error {
	if _, err := exec.LookPath("tesseract"); err == nil {
		return nil
	}

	tempDir := filepath.Join(os.TempDir(), "tesseract-installer")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return InternalError(fmt.Sprintf("не удалось создать каталог: %v", err))
	}

	if !checkInternetConnection() {
		return UserError("проверьте подключение к интернету")
	}

	InfoMessage("Установка Tesseract")
	if err := tesseract.Install(tempDir); err != nil {
		return InternalError(err.Error())
	}

	return nil
}

// findSystemPython - ищет Python в системе, если не находит, возвращает ошибку
func findSystemPython() (string, error) {
	if path, err := exec.LookPath("py"); err == nil {
		return path, nil
	}

	if path, err := exec.LookPath("python3"); err == nil {
		return path, nil
	}

	if path, err := exec.LookPath("python"); err == nil {
		cmd := command.Command("python", "--version")
		out, _ := cmd.CombinedOutput()
		if !bytes.Contains(out, []byte("Microsoft Store")) {
			return path, nil
		}
	}

	return "", fmt.Errorf("реальный Python не найден")
}

// checkInternetConnection - проверяет наличие интернета через попытку подключиться к PyPi
func checkInternetConnection() bool {
	cmd := command.Command("ping", "-n", "1", "-w", "3000", "pypi.org")
	return cmd.Run() == nil
}

// CheckInitWasUsed - проверяет, был ли уже запущен init, чтобы все зависисмости были установлены
func CheckInitWasUsed(pathToVenv string) (bool, error) {
	venvPath := filepath.Join(pathToVenv, ".venv", "Lib", "site-packages")

	entries, err := os.ReadDir(venvPath)
	if err != nil {
		return false, fmt.Errorf("ошибка чтения папки с зависимостями или ее отсутствие: %w", err)
	}

	for _, entry := range entries {
		if entry.Name() == "yaml" && entry.IsDir() {
			return true, nil
		}
	}

	return false, nil
}
