package parser

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrNotFoundTable = errors.New("поле 'table' не найдено в JSON")

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
		return nil, ErrNotFoundTable
	}

	if tableBytes[0] == '[' {
		tableBytes = []byte(`{"table":` + string(tableBytes) + `}`)
	}

	return tableBytes, nil
}

// UpdateTableInData - меняет поле table в json на читаемую таблицу
func UpdateTableInData(data interface{}) (interface{}, error) {
	dataBytes, errs := json.Marshal(data)
	if errs != nil {
		return nil, errs
	}

	var dataMap map[string]interface{}
	if err := json.Unmarshal(dataBytes, &dataMap); err != nil {
		return nil, err
	}

	tableBytes, err := GetTableFromJson(dataMap)
	if err != nil {
		return nil, err
	}

	items, err := ParseNewTable(tableBytes)
	if err != nil {
		return nil, err
	}

	dataMap["table"] = items

	return dataMap, nil
}
