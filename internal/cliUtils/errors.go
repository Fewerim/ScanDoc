package cliUtils

import (
	"fmt"
	"proWeb/internal/exitCodes"

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

// UserError - возвращает пользовательскую ошибку приложения
func UserError(message string) *AppError {
	return newAppError(exitCodes.UserError, message)
}

// ServerError - возвращает серверную ошибку приложения
func ServerError(message string) *AppError {
	return newAppError(exitCodes.ServerError, message)
}

// InternalError - возвращает внутреннюю ошибку приложения
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
