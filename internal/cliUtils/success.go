package cliUtils

import (
	"fmt"
	"time"
)

const (
	success = 1
	fail    = -1
)

type AppSuccess struct {
	Status  int
	Message string
	Time    time.Duration
}

func newAppSuccess(status int, message string, time time.Duration) *AppSuccess {
	return &AppSuccess{status, message, time}
}

func (app *AppSuccess) ToString() string {
	result := fmt.Sprintf("Статус выполнения: %d\nВремя выполнения: %.3fs\nОписание результата: %s", app.Status, app.Time.Seconds(), app.Message)
	return result
}

func Success(message string, time time.Duration) *AppSuccess {
	return newAppSuccess(success, message, time)
}

func Error(message string, time time.Duration) *AppSuccess {
	return newAppSuccess(fail, message, time)
}
