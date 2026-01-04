package parser

import (
	"github.com/iancoleman/orderedmap"
)

// ParseNewTable - парсит поле table в читаемый вид
func ParseNewTable(jsonData []byte) (*Items, error) {
	table, err := ParseTable(jsonData)
	if err != nil {
		return nil, err
	}

	if len(table.Rows) == 0 {
		return &Items{}, nil
	}

	items := GetAllLines(table)

	return items, nil
}

func GetAllLines(table *Table) *Items {
	titlesNames, fin := GetAllTitlesNames(table)

	lineNumber := 1
	var items Items
	for s := fin; s < len(table.Rows); s = s + len(titlesNames) {
		var item Item
		item.LineNumber = lineNumber
		item.Fields = orderedmap.New()
		numTit := 0
		for v := s; v < s+len(titlesNames); v++ {
			item.Fields.Set(titlesNames[numTit], table.Rows[v].Text)
			numTit++
		}
		items.Items = append(items.Items, item)
		lineNumber++
	}
	return &items
}

func GetAllTitlesNames(table *Table) ([]string, int) {
	var titles []Cell
	var doubleTitles []Cell
	var titlesNames []string
	var fin int
	secondColX := table.Rows[0].Coordinates[0][0]
	secondColY := table.Rows[0].Coordinates[1][1]
	secondCol := []int{secondColX, secondColY}
	for i := 0; i < len(table.Rows); i++ {
		if secondCol[0] == table.Rows[i].Coordinates[0][0] && secondCol[1] == table.Rows[i].Coordinates[0][1] {
			fin = i
			break
		} else {
			doubleTitles = append(doubleTitles, table.Rows[i])
		}
	}
	for i := 0; i < (len(doubleTitles) / 2); i++ {
		titles = append(titles, doubleTitles[i])
	}

	var first = true

	for z := 0; z < len(titles); z++ {
		var currentTitle string = titles[z].Text
		var findCol []int
		if first {
			findColX := table.Rows[z].Coordinates[0][0]
			findColY := table.Rows[z].Coordinates[1][1]
			findCol = []int{findColX, findColY}
		} else {
			findColX := table.Rows[z].Coordinates[1][0]
			findColY := table.Rows[z].Coordinates[1][1]
			findCol = []int{findColX, findColY}
		}
		for k := (len(doubleTitles) / 2); k < len(doubleTitles); k++ {
			if first {
				if findCol[0] == doubleTitles[k].Coordinates[0][0] && findCol[1] == doubleTitles[k].Coordinates[0][1] {
					currentTitle = currentTitle + ", " + doubleTitles[k].Text
					first = !first
					break
				}
			} else {
				if findCol[0] == doubleTitles[k].Coordinates[1][0] && findCol[1] == doubleTitles[k].Coordinates[0][1] {
					currentTitle = currentTitle + ", " + doubleTitles[k].Text
					first = !first
					break
				}
			}
		}
		titlesNames = append(titlesNames, currentTitle)
	}
	return titlesNames, fin
}
