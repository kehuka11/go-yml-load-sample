package config

import (
	"html/template"
	"io"
	"os"
	"strings"
)

type YAMLTemplate struct{}

func NewYAMLTemplate() *YAMLTemplate {
	return &YAMLTemplate{}
}

func (t *YAMLTemplate) Compile(name string, r io.Reader, w io.Writer) error {
	// デリミタを ${...} に設定
	funcMap := template.FuncMap{
		"getenv": t.Getenv,
	}

	value, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	tmpl, err := template.New(name).
		Delims("${", "}").
		Funcs(funcMap).
		Parse(string(value))
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	return err
}

func (t *YAMLTemplate) Getenv(keyAndDefault string) string {
	parts := strings.SplitN(keyAndDefault, ":", 2)
	key := strings.TrimSpace(parts[0])
	var def string
	if len(parts) == 2 {
		def = strings.TrimSpace(parts[1])
	}

	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return def
}
