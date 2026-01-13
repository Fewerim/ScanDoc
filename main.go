package main

import (
	"fmt"
	"proWeb/internal/appUtils"
)

func main() {
	b, err := appUtils.CheckInitWasUsed()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)
}
