package main

import (
	"fmt"
	"proWeb/lib/config"
)

func main() {
	projectRoot, _ := config.FindProjectRoot(".")
	fmt.Println(projectRoot)
}
