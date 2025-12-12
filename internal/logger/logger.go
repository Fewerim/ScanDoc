package logger

import (
	"io"
	"log"
)

// Logger - интерфейс, чтобы протягивать логгер по всему проекту, где он необходим
type Logger interface {
	Info(op, message string)
	Error(op, message string)
	Debug(op, message string)
}

type Log struct {
	l *log.Logger
}

// New - конструктор логгера с , формат вывода в файл: DATA | TIME | Message
func New(w io.Writer) *Log {
	return &Log{log.New(w, "", log.Ldate|log.Ltime)}
}

// NewFileLog - создает логгер, который пишет в указанный по пути файл
func NewFileLog(filePath string) (*Log, error) {
	//TODO: реализовать открытие файла по пути (os.OpenFile) также прописать флаги на тот случай, если файла нет
	// Возвращать конструктор логгера с переданным туда файлом
	return &Log{}, nil
}

// Info - выводит операцию, где был выполнен лог и информацию об операции
// Пример вывода: [INFO] operation: op | message: message
func (s *Log) Info(op, message string) {
	//TODO:
}

// Error - выводит операцию, где произошла ошибка и информацию об ошибке выполнения кода
// Пример вывода: [ERROR] operation: op | message: message
func (s *Log) Error(op, message string) {
	//TODO:
}

// Debug - выводи операцию, где происходит debug и информацию для пользователя
// Пример вывода: [DEBUG] operation: op | message:
func (s *Log) Debug(op, message string) {
	//TODO:
}
