package appUtils

import (
	"fmt"
	"proWeb/internal/exitCodes"
	"strings"

	"github.com/fatih/color"
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

// ToString - переводит ошибки в текст
func (errs FileErrors) ToString() string {
	fileNames := make([]string, 0)
	for _, file := range errs {
		fileNames = append(fileNames, file.FileName)
	}
	res := "Файлы, которые не были обработаны: " + strings.Join(fileNames, ", ")
	return color.RedString(res)
}

// PrintErrors - выводит ошибки
func (errs FileErrors) PrintErrors() {
	fmt.Println(errs.ToString())
}

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

// FilesNotProcessed - возвращает ошибки файлов, которые не были обработаны
func FilesNotProcessed(filesErrs []FileError) {
	filesNotSuccess := FileErrors{}
	filesNotSuccess = filesErrs
	filesNotSuccess.PrintErrors()
}
