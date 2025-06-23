package config

import (
	"os"
	"testing"
	"testing/fstest"

	"bytes"
	"fmt"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

// テスト用のエラーを返すReader
type errorReader struct{}

func (e *errorReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("dummy read error")
}

// テスト用のエラーを返すWriter
type errorWriter struct{}

func (e *errorWriter) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("dummy write error")
}

func TestYAMLTemplate_Compile(t *testing.T) {
	t.Run("正常系-ymlを展開できること-環境変数の設定なし", func(t *testing.T) {
		os.Setenv("ENV", "test")
		defer os.Unsetenv("ENV")

		fs := fstest.MapFS{
			"config.yaml": &fstest.MapFile{
				Data: []byte("test:\n  check: ${getenv \"FLAG:true\"}\n"),
			},
		}

		buf := bytes.NewBuffer(nil)
		reader, err := fs.Open("config.yaml")
		if err != nil {
			t.Fatalf("failed to open config.yaml: %v", err)
		}
		defer reader.Close()

		yamlTemplate := NewYAMLTemplate()
		err = yamlTemplate.Compile("test", reader, buf)
		var c TestConfig
		yaml.Unmarshal(buf.Bytes(), &c)
		assert.NoError(t, err)
		assert.Equal(t, true, c.Test.Check)
	})

	t.Run("正常系-ymlを展開できること-環境変数の設定あり", func(t *testing.T) {
		os.Setenv("ENV", "test")
		os.Setenv("FLAG", "false")
		defer os.Unsetenv("ENV")
		defer os.Unsetenv("FLAG")

		fs := fstest.MapFS{
			"config.yaml": &fstest.MapFile{
				Data: []byte("test:\n  check: ${getenv \"FLAG:true\"}\n"),
			},
		}

		buf := bytes.NewBuffer(nil)
		reader, err := fs.Open("config.yaml")
		if err != nil {
			t.Fatalf("failed to open config.yaml: %v", err)
		}
		defer reader.Close()

		yamlTemplate := NewYAMLTemplate()
		err = yamlTemplate.Compile("test", reader, buf)
		var c TestConfig
		yaml.Unmarshal(buf.Bytes(), &c)
		assert.NoError(t, err)
		assert.Equal(t, false, c.Test.Check)
	})

	t.Run("異常系-リーダーエラーが発生する", func(t *testing.T) {

		buf := bytes.NewBuffer(nil)

		yamlTemplate := NewYAMLTemplate()
		err := yamlTemplate.Compile("test", &errorReader{}, buf)
		assert.Error(t, err)
	})

	t.Run("異常系-ymlのParseに失敗する", func(t *testing.T) {

		fs := fstest.MapFS{
			"config.yaml": &fstest.MapFile{
				Data: []byte(`test/check: ${getenv "FLAG:true}"`),
			},
		}

		buf := bytes.NewBuffer(nil)
		reader, err := fs.Open("config.yaml")
		if err != nil {
			t.Fatalf("failed to open config.yaml: %v", err)
		}
		defer reader.Close()

		yamlTemplate := NewYAMLTemplate()
		err = yamlTemplate.Compile("test", reader, buf)
		assert.Error(t, err)
	})

	t.Run("異常系-ymlのExecuteに失敗する", func(t *testing.T) {

		fs := fstest.MapFS{
			"config.yaml": &fstest.MapFile{
				Data: []byte("test:\n  check: ${getenv \"FLAG:true\"}\n"),
			},
		}

		reader, err := fs.Open("config.yaml")
		if err != nil {
			t.Fatalf("failed to open config.yaml: %v", err)
		}
		defer reader.Close()

		yamlTemplate := NewYAMLTemplate()
		err = yamlTemplate.Compile("test", reader, &errorWriter{})
		assert.Error(t, err)
	})
}

func TestYAMLTemplate_Getenv(t *testing.T) {
	t.Run("正常系-環境変数が設定されている場合", func(t *testing.T) {
		os.Setenv("ENV", "test2")
		defer os.Unsetenv("ENV")

		yamlTemplate := NewYAMLTemplate()
		result := yamlTemplate.Getenv("ENV:test")
		assert.Equal(t, "test2", result)
	})

	t.Run("正常系-環境変数が設定されていない場合", func(t *testing.T) {
		yamlTemplate := NewYAMLTemplate()
		result := yamlTemplate.Getenv("ENV:test")
		assert.Equal(t, "test", result)
	})
}
