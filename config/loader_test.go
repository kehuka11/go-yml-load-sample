package config

import (
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {

	t.Run("正常系-環境変数が設定されていない場合", func(t *testing.T) {
		os.Setenv("ENV", "test")
		loader := NewLoader[TestConfig](ConfFS)
		defer teardown()
		cfg, err := loader.LoadConfig()
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}
		// 全ての設定を確認することは省略
		assert.NoError(t, err)
		assert.Equal(t, 8081, cfg.Server.Port)
	})
	t.Run("正常系-環境変数が設定されている場合", func(t *testing.T) {
		os.Setenv("ENV", "test")
		os.Setenv("SERVER_PORT", "8082")
		defer teardown()
		loader := NewLoader[TestConfig](ConfFS)
		cfg, err := loader.LoadConfig()
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}
		// 全ての設定を確認することは省略
		assert.NoError(t, err)
		assert.Equal(t, 8082, cfg.Server.Port)
	})

	t.Run("正常系-環境ごとのyamlの差分が反映されている", func(t *testing.T) {
		os.Setenv("ENV", "test")
		defer teardown()
		loader := NewLoader[TestConfig](ConfFS)
		cfg, err := loader.LoadConfig()
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}
		// 全ての設定を確認することは省略
		assert.NoError(t, err)
		assert.Equal(t, true, cfg.Test.Check)
	})

	t.Run("正常系-環境変数を使用しない設定べた書き", func(t *testing.T) {
		os.Setenv("ENV", "test")
		defer teardown()
		loader := NewLoader[TestConfig](ConfFS)
		cfg, err := loader.LoadConfig()
		assert.NoError(t, err)
		assert.Equal(t, true, cfg.Test.Normal)
	})

	t.Run("異常系-yamlファイルのパースに失敗する", func(t *testing.T) {
		os.Setenv("ENV", "test")
		defer teardown()
		// 不正なYAML形式のファイルを作成
		fs := fstest.MapFS{
			"config.yaml": &fstest.MapFile{
				Data: []byte(`server/invalid: yaml: [  # 不正なYAML形式`),
			},
		}

		loader := NewLoader[TestConfig](fs)
		_, err := loader.LoadConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal YAML")
	})

	t.Run("異常系-yamlCompileに失敗する", func(t *testing.T) {
		os.Setenv("ENV", "test")
		defer teardown()
		// テンプレート構文エラーのあるYAMLファイルを作成
		fs := fstest.MapFS{
			"config.yaml": &fstest.MapFile{
				Data: []byte(`server:port: ${ // }  # 存在しないテンプレート変数`),
			},
		}

		loader := NewLoader[TestConfig](fs)
		_, err := loader.LoadConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to compile template")
	})
}

func teardown() {
	os.Unsetenv("ENV")
	os.Unsetenv("SERVER_PORT")
}
