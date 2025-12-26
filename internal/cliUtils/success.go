package cliUtils

import (
	"fmt"
	"proWeb/internal/exitCodes"

	"github.com/fatih/color"
)

type Message interface {
	ToString() string
}

type AppSuccess struct {
	Status  int
	Message Message
}

func newAppSuccess(status int, message Message) *AppSuccess {
	return &AppSuccess{status, message}
}

// NewSuccess - конструктор, возвращающий статус успеха приложения
func NewSuccess(message Message) *AppSuccess {
	return newAppSuccess(exitCodes.Success, message)
}

// ToString - возвращает строку для вывода статуса успеха приложения
func (app *AppSuccess) ToString() string {
	result := fmt.Sprintf("Статус выполнения: %d\nОписание результата:\n%s", app.Status, app.Message.ToString())
	return color.GreenString(result)
}

// PrintSuccess - выводит в консоль статус успеха приложения
func (app *AppSuccess) PrintSuccess() {
	fmt.Println(app.ToString())
}
