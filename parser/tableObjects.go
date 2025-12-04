package parser

import "github.com/iancoleman/orderedmap"

type Cell struct {
	Coordinates [][]int `json:"coordinates"`
	Text        string  `json:"text"`
}

//type Row struct {
//	Cells []Cell `json:"cells"`
//}

type Table struct {
	Rows []Cell `json:"table"`
}

type Item struct {
	LineNumber int                    `json:"lineNumber"`
	Fields     *orderedmap.OrderedMap `json:"fields"`
}

type Items struct {
	Items []Item `json:"items"`
}
