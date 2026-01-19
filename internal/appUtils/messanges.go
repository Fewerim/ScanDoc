package appUtils

import "github.com/fatih/color"

// InfoMessage - информационные сообщения системы
func InfoMessage(text string) {
	color.Blue(text)
}

// SuccessMessage - сообщения системы об успехе
func SuccessMessage(text string) {
	color.Green(text)
}

// FailMessage - сообщения системы об ошибке
func FailMessage(text string) {
	color.Red(text)
}
