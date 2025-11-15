package main

import (
	"log"
	"proWeb/files"
)

func main() {
	if !files.StorageExists() {
		if err := files.CreateStorageJSON(); err != nil {
			log.Fatal(err)
		}
	}
	datas := []files.TestJson{}
	data1 := files.TestJson{
		Name:       "Danila",
		Age:        25,
		Profession: "Golang programmer",
		Skills: []string{
			"SQL",
			"Kafka",
			"Python",
		},
	}
	data2 := files.TestJson{
		Name:       "Danila",
		Age:        25,
		Profession: "Golang programmer",
		Skills: []string{
			"SQL",
			"Kafka",
			"Python",
		},
	}

	datas = append(datas, data1, data2)
	err := files.SaveFileToDirectory("jssdo", "УТП", datas)
	if err != nil {
		log.Fatal(err)
	}
}
