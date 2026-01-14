package main

import (
	"fmt"
	"proWeb/internal/appUtils"
)

func main() {
	var filePath string
	fmt.Scan(&filePath)

	img, err := appUtils.ConvertPdfToImg(filePath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(img)
}
