package typesUtils

import (
	"errors"
	"fmt"
	typesJSON2 "proWeb/internal/typesJSON"
	"strings"
)

type DocFactory func() interface{}

// docFactories - хранилище допустимых типов документов
var docFactories = map[string]DocFactory{
	"upd":     func() interface{} { return &typesJSON2.Upd{} },
	"invoice": func() interface{} { return &typesJSON2.TheInvoice{} },
	"torg12":  func() interface{} { return &typesJSON2.Torg12{} },
}

// normalize - нормализует строку переводя ее в верхний регистр и убирая лишние пробелы
func normalize(tDocument string) string {
	return strings.ToLower(strings.TrimSpace(tDocument))
}

// GetJSONStruct - фабрика, возвращает указатель на объект JSON структуры по переданному типу документа
func GetJSONStruct(typeOfDocument string) (data interface{}, err error) {
	tDocument := normalize(typeOfDocument)

	if tDocument == "" {
		return nil, errors.New("тип документа не был передан в результате запроса")
	}

	factory, ok := docFactories[tDocument]
	if !ok {
		return nil, fmt.Errorf("неизвестный тип документа %s", tDocument)
	}
	return factory(), nil
}
