package files

import "strings"

// addExtensionJSON - добавляет расширение .json к названию файла, если его нет
func addExtensionJSON(filename string) string {
	if !strings.HasSuffix(filename, ".json") {
		return filename + ".json"
	}
	return filename
}
