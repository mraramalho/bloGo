package main

import (
	"bytes"
	"os"

	"github.com/yuin/goldmark"
)

// ConvertMarkdownToHTML lê um arquivo .md e converte para HTML
func ConvertMarkdownToHTML(filepath string) (string, error) {
	// Lê o conteúdo do arquivo Markdown
	mdContent, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	// Converte para HTML
	var buf bytes.Buffer
	if err := goldmark.Convert(mdContent, &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
