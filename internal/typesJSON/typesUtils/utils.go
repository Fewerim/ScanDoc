package typesUtils

import (
	"fmt"
	"proWeb/internal/typesJSON"
	"strings"
)

const (
	typeUPD     = "UPD"
	typeINVOICE = "INVOICE"
	typeTORG    = "TORG"
)

// validateDocumentType - нормализует и валидирует тип документа
// ТИПЫ ДОКУМЕНТОВ: UPD | TORG | INVOICE
func validateDocumentType(typeOfDocument string) (result string, err error) {
	tDocument := strings.ToUpper(typeOfDocument)

	switch tDocument {
	case typeUPD, typeINVOICE, typeTORG:
		return tDocument, nil
	default:
		return "", fmt.Errorf("неизвестный тип документа %v", typeOfDocument)
	}
}

// GetJSONStruct - фабрика, возвращает указатель на объект JSON структуры по переданному типу документа
func GetJSONStruct(typeOfDocument string) (data interface{}, err error) {
	tDocument, err := validateDocumentType(typeOfDocument)
	if err != nil {
		return nil, err
	}

	switch tDocument {
	case typeUPD:
		return &typesJSON.Upd{}, nil
	case typeINVOICE:
		return &typesJSON.TheInvoice{}, nil
	case typeTORG:
		return &typesJSON.Torg12{}, nil
	default:
		return nil, fmt.Errorf("неизвестный тип документа : %v", typeOfDocument)
	}
}
