package main

import (
	"fmt"
	"proWeb/internal/appUtils"
)

func main() {
	dir, err := appUtils.ConvertPdfToImg("C:\\Users\\Administrator\\Downloads\\Telegram Desktop\\4ТОРГ-12.pdf")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("JPG сохранены в:", dir)
}
