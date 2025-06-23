package config

import (
	"io/fs"
	"sync"
)

type GlobalLoader[C any] struct {
	loader *Loader[C]
	conf   *C
	mu     sync.Mutex
}

func NewGlobalLoader[C any](fs_ fs.FS) *GlobalLoader[C] {
	return &GlobalLoader[C]{
		loader: NewLoader[C](fs_),
	}
}

func (gl *GlobalLoader[C]) Load() error {
	gl.mu.Lock()
	defer gl.mu.Unlock()

	conf, err := gl.loader.LoadConfig()
	if err != nil {
		return err
	}
	gl.conf = conf

	return nil
}

func (gl *GlobalLoader[C]) GetConfig() *C {
	if gl.conf == nil {
		err := gl.Load()
		if err != nil {
			panic(err)
		}
	}
	return gl.conf
}
