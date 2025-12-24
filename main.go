package main

import (
	"fmt"
	"proWeb/internal/cliUtils"
	"time"
)

func main() {
	multi1 := cliUtils.MultiProcessResult{
		Results: []cliUtils.Result{
			{FileName: "file1.json", Location: "storageJSON", CreatedAt: time.Now()},
			{FileName: "file2.json", Location: "storageJSON", CreatedAt: time.Now().Add(time.Second)},
			{FileName: "file3.json", Location: "storageJSON", CreatedAt: time.Now().Add(2 * time.Second)},
		},
	}
	fmt.Print(multi1.ToString())
}
