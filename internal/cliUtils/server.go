package cliUtils

import (
	"fmt"
	"proWeb/internal/files"
	"time"
)

type Response struct {
	File string
	Time time.Time
}

func createResponse(file string) Response {
	return Response{
		File: file,
		Time: time.Now(),
	}
}

func ProcessOnceFile(filePath, createdNameFile string, port int) (Response, error) {
	err := startPythonServer(fmt.Sprintf(":%d", port))
	if err != nil {
		info := fmt.Sprintf("ошибка при старте сервера: %v", err)
		return Response{}, ServerError(info)
	}

	response := createResponse(filePath)

	//TODO: добавить отправку файла на сервер
	//sendFileToServer(filePath, )

	err = files.SaveFileToStorage(createdNameFile, response)
	if err != nil {
		info := fmt.Sprintf("ошибка при попытке сохранить файл: %v", err)
		return Response{}, InternalError(info)
	}
	return response, nil
}

// startPythonServer - создает соединение с сервером
func startPythonServer(port string) error {
	time.Sleep(5 * time.Second)
	return nil
}

// sendFileToServer - отправляет файл на сервер и возвращает ответ
func sendFileToServer(filePath string) (Response, error) {
	panic("implement me")
}
