package cliUtils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"proWeb/lib/tesseract"
)

// CreateVenv - создает виртуальное окружения для python
func CreateVenv(pathToCreate string) error {
	venvPath := filepath.Join(pathToCreate, ".venv")

	if _, err := os.Stat(".venv"); err == nil {
		return nil
	}

	py, err := findSystemPython()
	if err != nil {
		return InternalError(err.Error())
	}

	cmd := exec.Command(py, "-m", "venv", venvPath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = pathToCreate

	if err = cmd.Run(); err != nil {
		info := fmt.Sprintf("ошибка создания venv: %w", err)
		return InternalError(info)
	}

	return nil
}

// InstallRequirements - устанавливает все необходимые зависимости для python скрипта
func InstallRequirements(pathToScript string) error {
	reqsPath := filepath.Join(pathToScript, "requirements.txt")
	if _, err := os.Stat(reqsPath); os.IsNotExist(err) {
		return UserError("файл requirements.txt не найден")
	}

	pipPath := filepath.Join(".venv", "Scripts", "pip.exe")
	if _, err := os.Stat(pipPath); os.IsNotExist(err) {
		info := fmt.Sprintf("pip не найден в .venv: %s", pipPath)
		return UserError(info)
	}

	if !checkInternetConnection() {
		return UserError("не удается подключиться к PyPi. проверьте подключение к интернету")
	}

	cmd := exec.Command(pipPath, "install", "-r", reqsPath)

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
		cmd := exec.Command("python", "--version")
		out, _ := cmd.CombinedOutput()
		if !bytes.Contains(out, []byte("Microsoft Store")) {
			return path, nil
		}
	}

	return "", fmt.Errorf("реальный Python не найден")
}

// checkInternetConnection - проверяет наличие интернета через попытку подключиться к PyPi
func checkInternetConnection() bool {
	cmd := exec.Command("ping", "-n", "1", "-w", "3000", "pypi.org")
	return cmd.Run() == nil
}
