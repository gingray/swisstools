package api

import (
	"encoding/json"
	"strings"
)

func ParseResponse(output []byte) any {
	var jsonResponse interface{}
	err := json.Unmarshal([]byte(output), &jsonResponse)
	if err == nil {
		return jsonResponse
	}
	linesIter := strings.Lines(string(output))
	var arrStr []string
	for line := range linesIter {
		arrStr = append(arrStr, strings.TrimSpace(line))
	}
	return arrStr
}
