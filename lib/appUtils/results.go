package appUtils

import (
	"fmt"
	"path/filepath"
	"proWeb/lib/typesJSON/typesUtils"
	"strings"
	"time"
)

// Result - результат выполнения CLI команды (имя созданного файла, время создания)
type Result struct {
	FileName  string
	Location  string
	CreatedAt time.Time
	DocType   string
}

type OnceProcessResult struct {
	Result  Result
	Elapsed time.Duration
}

type MultiProcessResult struct {
	Results []Result
	Elapsed time.Duration
}

type InitResult struct {
	Message string
}

// CreateResult - конструктор для создания результата выполнения CLI команды
func CreateResult(fileName, docType, location string) Result {
	return Result{
		FileName:  fileName,
		Location:  location,
		CreatedAt: time.Now(),
		DocType:   typesUtils.TypesDoc[docType],
	}
}

// CreateResultWithFolder - создает результат с подпапкой в storage
func CreateResultWithFolder(fileName, docType, storagePath, folderName string) Result {
	fullPath := storagePath
	if folderName != "" {
		fullPath = filepath.Join(storagePath, folderName)
	}

	return CreateResult(fileName, docType, fullPath)
}

// CreateOnceProcessResult - конструктор для создания результата выполнения CLI команды run_once
func CreateOnceProcessResult(fileName, docType, location string) OnceProcessResult {
	return OnceProcessResult{
		Result:  CreateResult(fileName, docType, location),
		Elapsed: 0,
	}
}

// ToString - возвращает строку для вывода результата выполнения CLI команды run_once
func (res *OnceProcessResult) ToString() string {
	date := res.Result.CreatedAt.Format("2006/01/02 | 15:04:05")

	s := fmt.Sprintf("В хранилище:\t%s\nСоздан файл:\t%s\nТип документа:\t%s\nВремя создания:\t%s\nВремя выполнения:\t%.3fs",
		res.Result.Location, res.Result.FileName, res.Result.DocType, date, res.Elapsed.Seconds())
	return s
}

// SetElapsedTime - устанавливает прошедшее время внутрь объекта
func (res *OnceProcessResult) SetElapsedTime(elapsed time.Duration) {
	res.Elapsed = elapsed
}

// GetElapsedTime - возвращает прошедшее время
func (res *OnceProcessResult) GetElapsedTime() float64 {
	return res.Elapsed.Seconds()
}

// CreateMultiProcessResult - конструктор для создания результата выполнения CLI команды run_multi
func CreateMultiProcessResult() MultiProcessResult {
	return MultiProcessResult{
		Results: make([]Result, 0),
		Elapsed: 0,
	}
}

// ToString - возвращает строку для вывода результата выполнения CLI команды run_multi
func (res *MultiProcessResult) ToString() string {
	date := res.Results[len(res.Results)-1].CreatedAt.Format("2006/01/02 | 15:04:05")
	location := res.Results[0].Location

	filesList := make([]string, len(res.Results))
	for i, result := range res.Results {
		filesList[i] = result.FileName
	}
	files := strings.Join(filesList, ", ")
	s := fmt.Sprintf("В хранилище:\t%s\nСозданы файлы:\t%s\nВремя создания:\t%s\nВремя выполнения:\t%.3fs",
		location, files, date, res.Elapsed.Seconds())
	return s
}

// SetElapsedTime - устанавливает прошедшее время внутрь объекта
func (res *MultiProcessResult) SetElapsedTime(elapsed time.Duration) {
	res.Elapsed = elapsed
}

// GetElapsedTime - возвращает прошедшее время
func (res *MultiProcessResult) GetElapsedTime() float64 {
	return res.Elapsed.Seconds()
}

// SetResult - устанавливает Result внутрь объекта
func (res *MultiProcessResult) SetResult(result Result) {
	res.Results = append(res.Results, result)
}

// CreateInitResult - конструктор для создания результата выполнения CLI команды init
func CreateInitResult(textMsg string) InitResult {
	return InitResult{
		Message: textMsg,
	}
}

// ToString - возвращает строку для вывода результата выполнения CLI команды init
func (res *InitResult) ToString() string {
	return res.Message
}
