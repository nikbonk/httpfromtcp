package headers

import (
	"bytes"
	"errors"
	"strings"
	"unicode"
)

type Headers map[string]string

const (
	clrf                = "\r\n"
	allowedSpecialChars = "!#$%&'*+-.^_`|~"
)

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(clrf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return idx + len(clrf), true, nil
	}
	headerText := data[:idx]
	normalizedHeaderText := strings.TrimLeftFunc(string(headerText), unicode.IsSpace)
	if len(normalizedHeaderText) < len(headerText) {
		return 0, false, errors.New("header text starts with space")
	}
	normalizedHeaderText = strings.TrimFunc(normalizedHeaderText, unicode.IsSpace)
	headerParts := strings.SplitN(normalizedHeaderText, ":", 2)
	if len(headerParts) != 2 {
		return 0, false, errors.New("invalid header format")
	}
	headerName := strings.TrimSpace(headerParts[0])
	if len(headerName) < len(headerParts[0]) {
		return 0, false, errors.New("header name includes spaces")
	}
	headerName = strings.ToLower(headerName)
	if !isAllowedString(headerName) {
		return 0, false, errors.New("header name contains invalid characters")
	}

	headerValue := strings.TrimSpace(headerParts[1])
	h.Set(headerName, headerValue)
	return idx + len(clrf), false, nil
}
func (h Headers) Set(name, value string) {
	h[name] = value
}

func isAllowedString(header string) bool {
	for _, char := range header {
		if ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || ('0' <= char && char <= '9') || strings.ContainsRune(allowedSpecialChars, char) {
			continue
		}
		return false
	}
	return true
}
