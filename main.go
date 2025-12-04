package main

import (
	"encoding/json"
	"fmt"
	"proWeb/parser"
)

func main() {
	//data := typesJSON.Torg12{}
	//err := files.SaveFileToDirectory("testTorg", "Торг-12", data)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//var m *orderedmap.OrderedMap
	//m = orderedmap.New()
	//m.Set("lineNumber", 1)
	//m.Set("Товар", "Печенье какое-то")
	//m.Set("Тип операции", "Продажа")
	//m.Set("Количество, шт", "20")
	//jsonData := []byte(`{
	//	"table": [
	//		{
	//			"coordinates": [[19, 23], [233, 64]],
	//			"text": "Товар"
	//		},
	//		{
	//			"coordinates": [[233, 23], [447, 64]],
	//			"text": "Тип операции"
	//		},
	//		{
	//			"coordinates": [[447, 23], [661, 64]],
	//			"text": "Количество, шт"
	//		},
	//		{
	//			"coordinates": [[19, 64], [233, 105]],
	//			"text": "Печенье какое-то"
	//		},
	//		{
	//			"coordinates": [[233, 64], [447, 105]],
	//			"text": "Продажа"
	//		},
	//		{
	//			"coordinates": [[447, 64], [661, 105]],
	//			"text": "20"
	//		},
	//		{
	//			"coordinates": [[19, 105], [233, 144]],
	//			"text": "Помидоры"
	//		},
	//		{
	//			"coordinates": [[233, 105], [447, 144]],
	//			"text": "Поступление"
	//		},
	//		{
	//			"coordinates": [[447, 105], [661, 144]],
	//			"text": "228"
	//		}
	//	]
	//}`)
	//jsonData := []byte(`{
	//"table": [
	//    {
	//        "coordinates": [[19, 23], [233, 64]],
	//        "text": "Товар"
	//    },
	//    {
	//        "coordinates": [[233, 23], [447, 64]],
	//        "text": "Тип операции"
	//    },
	//    {
	//        "coordinates": [[447, 23], [661, 64]],
	//        "text": "Количество, шт"
	//    },
	//    {
	//        "coordinates": [[661, 23], [800, 64]],
	//        "text": "Еще столбец"
	//    },
	//    {
	//        "coordinates": [[19, 64], [233, 105]],
	//        "text": "Печенье какое-то"
	//    },
	//    {
	//        "coordinates": [[233, 64], [447, 105]],
	//        "text": "Продажа"
	//    },
	//    {
	//        "coordinates": [[447, 64], [661, 105]],
	//        "text": "20"
	//    },
	//    {
	//        "coordinates": [[661, 64], [800, 105]],
	//        "text": "че нибудь"
	//    },
	//    {
	//        "coordinates": [[19, 105], [233, 144]],
	//        "text": "Помидоры"
	//    },
	//    {
	//        "coordinates": [[233, 105], [447, 144]],
	//        "text": "Поступление"
	//    },
	//    {
	//        "coordinates": [[447, 105], [661, 144]],
	//        "text": "228"
	//    },
	//    {
	//        "coordinates": [[661, 105], [800, 144]],
	//        "text": "еще че нибудь"
	//    }
	//]
	//}`)
	jsonData := []byte(`{
    "table": [
        {
            "coordinates": [[19, 23], [233, 64]],
            "text": "Товар"
        },
        {
            "coordinates": [[233, 23], [447, 64]],
            "text": "Тип операции"
        },
        {
            "coordinates": [[19, 64], [233, 105]],
            "text": "Печенье какое-то"
        },
        {
            "coordinates": [[233, 64], [447, 105]],
            "text": "Продажа"
        },
        {
            "coordinates": [[19, 105], [233, 144]],
            "text": "Помидоры"
        },
        {
            "coordinates": [[233, 105], [447, 144]],
            "text": "Поступление"
        }
    ]
	}`)
	items, err := parser.ProcessJsonTable(jsonData)
	if err != nil {
		fmt.Println(err)
	}
	jsonD, errs := json.MarshalIndent(items, "", "    ")
	if errs != nil {
		fmt.Println(errs)
	}
	fmt.Println(string(jsonD))
}
