package wz

type FileEntry struct {
	name string
	path string
}

func (e FileEntry) Name() string {
	return e.name
}

func (e FileEntry) Path() string {
	return e.path
}

func NewFileEntry(name string, path string) FileEntry {
	return FileEntry{
		name: name,
		path: path,
	}
}
