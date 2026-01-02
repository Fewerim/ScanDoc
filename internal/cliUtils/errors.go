package cliUtils

import (
	"fmt"
	"proWeb/internal/exitCodes"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/sys/windows/registry"
)

type ExitCoder interface {
	ExitCode() int
}

// AppError - структура для вывода и формирования ошибок приложения
type AppError struct {
	exitCode    int
	userMessage string
}

type FileError struct {
	FileName string
	Err      error
}

type FileErrors []FileError

// newAppError - конструктор для создания новой ошибки приложения
func newAppError(exitCode int, message string) *AppError {
	return &AppError{exitCode, message}
}

// Error - возвращает текст ошибки
func (err *AppError) Error() string {
	return err.userMessage
}

// ExitCode возвращает поле кода возврата
func (err *AppError) ExitCode() int {
	return err.exitCode
}

// ToString - переводит ошибку в строку для чтения пользователя
func (err *AppError) ToString() string {
	result := fmt.Sprintf("Код ошибки: [%d]\nОписание ошибки: %s", err.exitCode, err.userMessage)

	return color.RedString(result)
}

// UserError - возвращает пользовательскую ошибку приложения (1)
func UserError(message string) *AppError {
	return newAppError(exitCodes.UserError, message)
}

// ServerError - возвращает серверную ошибку приложения (2)
func ServerError(message string) *AppError {
	return newAppError(exitCodes.ServerError, message)
}

// InternalError - возвращает внутреннюю ошибку приложения (3)
func InternalError(message string) *AppError {
	return newAppError(exitCodes.InternalError, message)
}

// GetExitCode - возвращает статус ExitCode ошибки или defaultCode, если ошибка соответствует интерфейсу ExitCoder
func GetExitCode(err error, defaultCode int) int {
	if exitCoder, ok := err.(ExitCoder); ok {
		return exitCoder.ExitCode()
	}
	return defaultCode
}

func FilesNotProcessed(filesErrs []FileError) {
	filesNotSuccess := FileErrors{}
	filesNotSuccess = filesErrs
	filesNotSuccess.PrintNotSuccess()
}

func CheckIsAutorunCorrect() (bool, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `Software\Microsoft\Command Processor`, registry.QUERY_VALUE|registry.READ)
	if err != nil {
		return false, err
	}
	defer key.Close()

	autorun, _, err := key.GetStringValue("Autorun")
	if err != nil {
		fmt.Println("Autorun не установлен:", err)
		return true, err
	}

	if autorun != "@chcp 65001>nul" && autorun != "chcp 65001" {
		return true, nil
	}

	return false, nil
}

func PrintInstruction() {
	var sb strings.Builder
	sb.WriteString("Для корректной работы приложения установите UTF-8 кодировку в консоли\n")
	sb.WriteString("Вот как это можно сделать:\n")
	sb.WriteString("нажмите Win + R\n")
	sb.WriteString("Введите regedit\n")
	sb.WriteString("Нажмите Ок\n")
	sb.WriteString("В появившемся окне, найдите строчку Autorun, кликнете по ней два раза\n")
	sb.WriteString("Поменяйте значение на @chcp 65001>nul\n")
	sb.WriteString("Нажмите Ок и запустите команду заново\n")

	instructions := sb.String()
	color.Blue(instructions)
}
