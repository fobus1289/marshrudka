package request

import "errors"

var (
	singleFiles = files{}
	singleFile  = &file{
		Error: errors.New("form data file not found"),
	}
)

type files []IFile

type fileMap map[string]files

type IFiles interface {
	HasMore() bool
	Len() int
	LastIndex() int
	StoreAndCloseAll(dir string) IFiles
	FullPaths() []string
	Files() []IFile
	First() IFile
	Last() IFile
	Get(index int) IFile
}

func (fs files) FullPaths() []string {
	var paths []string

	for _, file := range fs {
		paths = append(paths, file.GetFullPath())
	}

	return paths
}

func (fs files) StoreAndCloseAll(dir string) IFiles {

	for _, file := range fs {
		file.StoreAndClose(dir)
	}

	return fs
}

func (fs files) HasMore() bool {
	return fs.Len() > 1
}

func (fs files) Len() int {
	return len(fs)
}

func (fs files) LastIndex() int {
	l := fs.Len()

	if l-1 >= 0 {
		return l
	}

	return 0
}

func (fs files) Files() []IFile {
	return fs
}

func (fs files) First() IFile {
	return fs.Get(0)
}

func (fs files) Last() IFile {
	return fs.Get(fs.LastIndex())
}

func (fs files) Get(index int) IFile {

	if index >= 0 && fs.Len()-1 >= index {
		return fs[index]
	}

	return singleFile
}

type IFileMap interface {
	Keys() []string
	Get(key string) IFiles
	Len() int
}

func (fm fileMap) Get(key string) IFiles {

	file := fm[key]

	if file == nil {
		return singleFiles
	}

	return file
}

func (fm fileMap) Len() int {
	return len(fm)
}

func (fm fileMap) Keys() []string {

	var keys = []string{}

	if fm.Len() == 0 {
		return keys
	}

	for key := range fm {
		keys = append(keys, key)
	}

	return keys
}
