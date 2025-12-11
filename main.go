package main

import (
	"log"
	"os"
)

func main() {
	wd, _ := os.Getwd()
	log.Println("WD:", wd)
}
