package handlers

import (
	"errors"
	"strings"
)

const (
	DocxFormat = "docx"
	PdfFormat  = "pdf"
	JpgFormat  = "jpg"
	PngFormat  = "png"
	TxtFormat  = "txt"
)

var allowedFormats = map[string]bool{
	DocxFormat: true,
	PdfFormat:  true,
	JpgFormat:  true,
	PngFormat:  true,
	TxtFormat:  true,
}

func ParseDocumentFormat(fileName string) (string, error) {
	if fileName == "" {
		return "", errors.New("file name is empty")
	}

	idx := strings.LastIndex(fileName, ".")
	if idx == -1 || idx == len(fileName)-1 {
		return "", errors.New("no valid file extension")
	}

	extension := fileName[idx+1:]
	extension = strings.ToLower(extension)

	if !allowedFormats[extension] {
		return "", errors.New("unsupported format: " + extension)
	}

	return extension, nil
}

func ParseDocumentType(fileName string) ProcessResult {
	format, err := ParseDocumentFormat(fileName)

	if err != nil {
		return ProcessResult{
			Error: err.Error(),
		}
	}

	return ProcessResult{
		DocumentType: format,
	}
}
