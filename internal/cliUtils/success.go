package cliUtils

import "fmt"

const (
	success = 1
	fail    = -1
)

type AppSuccess struct {
	Status  int
	Message string
}

func newAppSuccess(status int, message string) *AppSuccess {
	return &AppSuccess{status, message}
}

func (app *AppSuccess) ToString() string {
	result := fmt.Sprintf("Статус выполнения: %d\nОписание результата: %s", app.Status, app.Message)
	return result
}

func Success(message string) *AppSuccess {
	return newAppSuccess(success, message)
}

func Error(message string) *AppSuccess {
	return newAppSuccess(fail, message)
}
