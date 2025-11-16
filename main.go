package main

import (
	"fmt"
	"proWeb/files"
	"proWeb/typesJSON"
)

func main() {
	data := typesJSON.Torg12{}
	err := files.SaveFileToDirectory("testTorg", "Торг-12", data)
	if err != nil {
		fmt.Println(err)
	}
}
