package typesUtils

import (
	"errors"
	"fmt"
	"net/http"
)

func getDoctype(resp *http.Response) (string, error) {
	doctype := resp.Header.Get("document_type")
	if doctype == "" {
		return "", errors.New("тип документа не найден в заголовках ответа")
	}

	_, err := GetJSONStruct(doctype)
	if err != nil {
		return "", fmt.Errorf("невалидный тип документа: %w", err)
	}

	return doctype, nil
}
