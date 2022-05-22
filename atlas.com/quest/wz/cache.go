package wz

import (
	"errors"
	"fmt"
	"sync"
)

type fileCache struct {
	files map[string]FileEntry
}

var cache *fileCache
var once sync.Once

func GetFileCache() *fileCache {
	return cache
}

func (e *fileCache) Init(wzPath string) {
	once.Do(func() {
		es, err := Read(wzPath)
		if err != nil {
			panic(err)
		}

		var files = make(map[string]FileEntry)
		for _, e := range es {
			files[e.Name()] = e
		}

		cache = &fileCache{files: files}
	})
}

func (e *fileCache) GetFile(name string) (*FileEntry, error) {
	if val, ok := e.files[name]; ok {
		return &val, nil
	} else {
		return nil, errors.New(fmt.Sprintf("file %s not found", name))
	}
}
