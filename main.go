package main

import (
	"fmt"
	"proWeb/lib/cliUtils"
)

func main() {
	b, err := cliUtils.CheckInitWasUsed()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)
}
