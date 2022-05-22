package wz

import (
	"io/ioutil"
	"os"
)

func Read(path string) ([]FileEntry, error) {
	return read(path)
}

func read(path string) ([]FileEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		fs, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, err
		}

		var fes []FileEntry
		for _, cf := range fs {
			ces, err := read(path + "/" + cf.Name())
			if err != nil {
				return nil, err
			}
			fes = append(fes, ces...)
		}
		return fes, nil
	} else {
		return []FileEntry{NewFileEntry(stat.Name(), f.Name())}, nil
	}
}
