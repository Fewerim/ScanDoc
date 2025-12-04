package parser

import (
	"encoding/json"

	"github.com/iancoleman/orderedmap"
)

func ProcessJsonTable(jsonData []byte) (*Items, error) {
	table, err := ParseTable(jsonData)
	if err != nil {
		return nil, err
	}

	if len(table.Rows) == 0 {
		return &Items{}, nil
	}

	columns := ParseColumns(table)

	fieldsName := GetFieldsNames(columns)

	items := ParceRows(columns, fieldsName)
	return &items, nil
}

func ParseTable(jsonData []byte) (*Table, error) {
	var table Table
	if err := json.Unmarshal(jsonData, &table); err != nil {
		return nil, err
	}
	return &table, nil
}

func ParseColumns(table *Table) [][]Cell {
	var columns [][]Cell
	secondColX := table.Rows[0].Coordinates[0][0]
	secondColY := table.Rows[0].Coordinates[1][1]
	secondCol := []int{secondColX, secondColY}

	for i := 0; i < len(table.Rows); i++ {
		if secondCol[0] == table.Rows[i].Coordinates[0][0] && secondCol[1] == table.Rows[i].Coordinates[0][1] {
			break
		} else {
			row := make([]Cell, 0)
			row = append(row, table.Rows[i])
			next := []int{table.Rows[i].Coordinates[0][0], table.Rows[i].Coordinates[1][1]}
			for n := i + 1; n < len(table.Rows); n++ {
				if next[0] == table.Rows[n].Coordinates[0][0] && next[1] == table.Rows[n].Coordinates[0][1] {
					row = append(row, table.Rows[n])
					next = []int{table.Rows[n].Coordinates[0][0], table.Rows[n].Coordinates[1][1]}
				}
			}
			columns = append(columns, row)
		}
	}
	return columns
}

func GetFieldsNames(columns [][]Cell) []string {
	fieldsName := []string{}
	for _, item := range columns {
		fieldsName = append(fieldsName, item[0].Text)
	}
	return fieldsName
}

func ParceRows(columns [][]Cell, fieldsNames []string) Items {
	lineNumber := 1
	var items Items
	for z := 1; z < len(columns[0]); z++ {
		var item Item
		item.LineNumber = lineNumber
		item.Fields = orderedmap.New()
		for v := 0; v < len(fieldsNames); v++ {
			item.Fields.Set(fieldsNames[v], columns[v][z].Text)
		}
		items.Items = append(items.Items, item)
		lineNumber++
	}
	return items
}
