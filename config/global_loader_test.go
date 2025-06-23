package config

import (
	"os"
	"testing"
	"testing/fstest"

	"sync"

	"github.com/stretchr/testify/assert"
)

func TestGlobalLoader(t *testing.T) {

	t.Run("正常系-環境設定取得", func(t *testing.T) {
		os.Setenv("ENV", "test")
		defer os.Unsetenv("ENV")

		loader := NewGlobalLoader[TestConfig](ConfFS)
		err := loader.Load()
		assert.NoError(t, err)
		assert.Equal(t, true, loader.conf.Test.Normal)
	})

	t.Run("正常系-スレッドセーフであること", func(t *testing.T) {
		loader := NewGlobalLoader[TestConfig](ConfFS)

		var wg sync.WaitGroup
		const goroutines = 10
		errs := make(chan error, goroutines)

		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := loader.Load()
				errs <- err
			}()
		}
		wg.Wait()
		close(errs)

		for err := range errs {
			assert.NoError(t, err)
		}
	})

	t.Run("異常系-環境設定取得失敗", func(t *testing.T) {
		os.Setenv("ENV", "test")
		defer os.Unsetenv("ENV")

		// テンプレート構文エラーのあるYAMLファイルを作成
		fs := fstest.MapFS{
			"config.yaml": &fstest.MapFile{
				Data: []byte(`server:port: ${ // }  # 存在しないテンプレート変数`),
			},
		}

		loader := NewGlobalLoader[TestConfig](fs)
		err := loader.Load()
		assert.Error(t, err)
	})
}

func TestGetConfig(t *testing.T) {
	t.Run("正常系-初回実行", func(t *testing.T) {
		os.Setenv("ENV", "test")
		defer os.Unsetenv("ENV")

		loader := NewGlobalLoader[TestConfig](ConfFS)
		cfg := loader.GetConfig()
		assert.Equal(t, true, cfg.Test.Normal)
	})

	t.Run("正常系-2回目以降の実行", func(t *testing.T) {

		// テスト用のfs.FSを作成
		fs := fstest.MapFS{
			"config.yaml": &fstest.MapFile{
				Data: []byte(`test:¥n	check: true`),
			},
			"config.test.yaml": &fstest.MapFile{
				Data: []byte(`test:¥n	check: false`),
			},
			"config.test2.yaml": &fstest.MapFile{
				Data: []byte(`test:¥n	check: true`),
			},
		}

		os.Setenv("ENV", "test")
		defer os.Unsetenv("ENV")

		loader := NewGlobalLoader[TestConfig](fs)

		// 1回目実行
		cfg := loader.GetConfig()
		assert.Equal(t, false, cfg.Test.Check)

		// 2回目実行
		os.Unsetenv("ENV")
		os.Setenv("ENV", "test2")
		cfg = loader.GetConfig()
		assert.Equal(t, false, cfg.Test.Check)
	})

	t.Run("異常系-環境設定取得失敗", func(t *testing.T) {

		// テスト用のfs.FSを作成
		fs := fstest.MapFS{
			"config.yaml": &fstest.MapFile{
				Data: []byte(`server/invalid: yaml: [  # 不正なYAML形式`),
			},
		}

		os.Setenv("ENV", "test")
		defer os.Unsetenv("ENV")

		loader := NewGlobalLoader[TestConfig](fs)
		assert.Panics(t, func() {
			loader.GetConfig()
		})
	})
}
