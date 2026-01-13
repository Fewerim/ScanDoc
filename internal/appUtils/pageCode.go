package appUtils

import (
	"strings"

	"github.com/fatih/color"
	"golang.org/x/sys/windows/registry"
)

const (
	pngCode1 = "@chcp 65001>nul"
	pngCode2 = "chcp 65001"
)

// CheckIsAutorunCorrect - проверяет установлена ли кодировка UTF-8, если нет, то выводит инструкцию для ее установки
func CheckIsAutorunCorrect() error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `Software\Microsoft\Command Processor`, registry.QUERY_VALUE|registry.READ)
	if err != nil {
		return InternalError("не удалось открыть ключ реестра Command Processor")
	}
	defer key.Close()

	autorun, _, err := key.GetStringValue("AutoRun")
	if err != nil {
		PrintInstructionToSetupAutoRun()
		return UserError("autorun параметр не найден в реестре")
	}

	if autorun != pngCode1 && autorun != pngCode2 {
		PrintInstructionToSetupAutoRun()
		return UserError("кодировка UTF-8 не установлена в autorun")
	}

	return nil
}

// PrintInstructionToSetupAutoRun - выводит в консоль инструкцию для установки autorun кодировки UTF-8
func PrintInstructionToSetupAutoRun() {
	var sb strings.Builder
	sb.WriteString("Для корректной работы приложения установите UTF-8 кодировку в консоли\n")
	sb.WriteString("Вот как это можно сделать:\n")
	sb.WriteString("Нажмите Win + R\n")
	sb.WriteString("Введите regedit\n")
	sb.WriteString("Нажмите Ок\n")
	sb.WriteString("В появившемся окне, перейдите в HKEY_LOCAL_MACHINE\\SOFTWARE\\Microsoft\\Command Processor\n")
	sb.WriteString("Найдите строковый параметр Autorun (тип REG_SZ), если его нет, создайте строковый параметр\n")
	sb.WriteString("Установите значение на @chcp 65001>nul\n")
	sb.WriteString("Нажмите Ок и перезапустите командную строку\n")

	instructions := sb.String()
	color.Blue(instructions)
}
