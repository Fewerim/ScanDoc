package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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
	dir := filepath.Dir(filePath)
	if dir == "" && dir != "." {
		if err := os.MkdirAll(dir, 0775); err != nil {
			return nil, fmt.Errorf("failed to create log directory '%s': %w", dir, err)
		}
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file '%s': %w", filePath, err)
	}
	return New(file), nil
}

// Info - выводит операцию, где был выполнен лог и информацию об операции
// Пример вывода: [INFO] operation: op | message: message
func (s *Log) Info(op, message string) {
	s.l.Printf("[INFO] operation: %s | message: %s\n", op, message)
}

// Error - выводит операцию, где произошла ошибка и информацию об ошибке выполнения кода
// Пример вывода: [ERROR] operation: op | message: message
func (s *Log) Error(op, message string, errorType int) {
	s.l.Printf("[ERROR] operation: %s | error type: %d | message: %s\n", op, errorType, message)
}

// Debug - выводит операцию, где происходит debug и информацию для пользователя
// Пример вывода: [DEBUG] operation: op | message:
func (s *Log) Debug(op, message string) {
	s.l.Printf("[DEBUG] operation: %s | message: %s", op, message)
}
