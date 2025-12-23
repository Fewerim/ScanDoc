package parser

import (
	"encoding/json"
	"fmt"
)

// GetTableFromJson - достает объект table(если оно есть) из json файла
func GetTableFromJson(data interface{}) ([]byte, error) {
	var jsonBytes []byte
	var err error

	switch v := data.(type) {
	case []byte:
		jsonBytes = v
	case string:
		jsonBytes = []byte(v)
	default:
		jsonBytes, err = json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("ошибка сериализации данных: %v", err)
		}
	}

	var jsonData map[string]json.RawMessage
	err = json.Unmarshal(jsonBytes, &jsonData)
	if err != nil {
		return nil, err
	}

	tableBytes, exists := jsonData["table"]
	if !exists {
		return nil, fmt.Errorf("поле 'table' не найдено в JSON")
	}

	if tableBytes[0] == '[' {
		tableBytes = []byte(`{"table":` + string(tableBytes) + `}`)
	}

	return tableBytes, nil
}
