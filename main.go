package main

import (
	"fmt"
	"log"
	"proWeb/internal/storage"
)

func main() {
	pathToStorage := "C:\\Projects\\PP2025-ProWeb\\storageJSONs"
	files, err := storage.GetStorageFiles(pathToStorage)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
}
