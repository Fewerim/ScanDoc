package parser

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/iancoleman/orderedmap"
)

// AddTotalFields - подсчитывает и добавляет поле "итого" на основе таблицы
func AddTotalFields(data interface{}) (interface{}, error) {
	root, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("root is not map[string]interface{}")
	}

	tableAny, ok := root["table"]
	if !ok {
		return nil, errors.New("missing table")
	}

	itemsPtr, ok := tableAny.(*Items)
	if !ok {
		return nil, fmt.Errorf("table expected *parser.Items, got %T", tableAny)
	}

	var (
		totalQuantity         float64
		totalGrossWeight      float64
		totalNetWeight        float64
		totalAmountWithoutVAT float64
		totalVATAmount        float64
		totalAmountWithVAT    float64
	)

	for i, item := range itemsPtr.Items {
		if i == 0 {
			continue
		}

		fieldsOM := item.Fields

		keys := fieldsOM.Keys()

		fieldQuantity := keys[7]
		fieldGrossWeight := keys[8]
		fieldNetWeight := keys[9]
		fieldAmountNoVAT := keys[11]
		fieldVATAmount := keys[13]
		fieldAmountWith := keys[14]

		totalQuantity += findFieldValue(fieldQuantity, fieldsOM)
		totalGrossWeight += findFieldValue(fieldGrossWeight, fieldsOM)
		totalNetWeight += findFieldValue(fieldNetWeight, fieldsOM)
		totalAmountWithoutVAT += findFieldValue(fieldAmountNoVAT, fieldsOM)
		totalVATAmount += findFieldValue(fieldVATAmount, fieldsOM)
		totalAmountWithVAT += findFieldValue(fieldAmountWith, fieldsOM)
	}

	totalsOM := orderedmap.New()
	totalsOM.Set("totalQuantity", roundTo3(totalQuantity))
	totalsOM.Set("totalGrossWeight", roundTo3(totalGrossWeight))
	totalsOM.Set("totalNetWeight", roundTo3(totalNetWeight))
	totalsOM.Set("totalAmountWithoutVAT", roundTo3(totalAmountWithoutVAT))
	totalsOM.Set("totalVATAmount", roundTo3(totalVATAmount))
	totalsOM.Set("totalAmountWithVAT", roundTo3(totalAmountWithVAT))

	root["totals"] = totalsOM

	return root, nil
}

// findFieldValue - достает значение поля и из упорядоченной мапы и конвертирует его в float64
func findFieldValue(key string, fieldsOM *orderedmap.OrderedMap) float64 {
	if v, ok := fieldsOM.Get(key); ok {
		switch val := v.(type) {
		case string:
			if val == "" {
				return 0
			}
			if f, err := parseRuNumber(val); err == nil {
				return f
			}
		case float64:
			return val
		case int, int64:
			return float64(val.(int64))
		}
	}
	return 0
}

// parseRuNumber - превращает string во float64
func parseRuNumber(s string) (float64, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, ".")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, ",", ".")

	if s == "" {
		return 0, fmt.Errorf("empty")
	}

	return strconv.ParseFloat(s, 64)
}

func roundTo3(f float64) float64 {
	return math.Round(f*1000) / 1000
}
