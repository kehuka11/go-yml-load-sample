package config

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"

	"gopkg.in/yaml.v3"

	"dario.cat/mergo"
)

const (
	BaseConfigFile = "config.yaml"
	EnvConfigFile  = "config.%s.yaml"
)

type Loader[C any] struct {
	fs fs.FS
}

func NewLoader[C any](fs_ fs.FS) *Loader[C] {
	return &Loader[C]{fs: fs_}
}

// 環境に応じてYAMLをマージ
func (l *Loader[C]) LoadConfig() (*C, error) {
	env := os.Getenv("ENV")

	files := []string{BaseConfigFile}
	if env != "" {
		files = append(files, fmt.Sprintf(EnvConfigFile, env))
	}

	finalConfig := new(C)

	for _, file := range files {
		reader, err := l.fs.Open(file)
		if err != nil {
			// ベースは必須、環境ファイルはoptional
			if file == files[0] || !os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to open config file: %w", err)
			}
			continue
		}
		defer reader.Close()

		buf := bytes.NewBuffer(nil)
		if err := NewYAMLTemplate().Compile("config", reader, buf); err != nil {
			return nil, fmt.Errorf("failed to compile template: %w", err)
		}

		var cfg C
		if err := yaml.Unmarshal(buf.Bytes(), &cfg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
		}

		// マージ: 環境ファイルの値で上書き
		mergeConfig(finalConfig, &cfg)
	}

	return finalConfig, nil
}

func mergeConfig[C any](base, override *C) {
	mergo.Merge(base, override, mergo.WithOverride)
}
