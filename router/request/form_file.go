package request

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
)

var noSuchFileError = errors.New("no such file")

func newNoSuchFile() *file {
	return &file{
		error: noSuchFileError,
	}
}

type fileContainer struct {
	files  []IFile
	errors []error
}

func (fc *fileContainer) Files() []IFile {
	return fc.files
}

func (fc *fileContainer) GetFirst() IFile {
	if fc.Has() {
		return fc.files[0]
	}
	return newNoSuchFile()
}

func (fc *fileContainer) Count() int {
	return len(fc.files)
}

func (fc *fileContainer) Has() bool {
	return fc.Count() > 0
}

func (fc *fileContainer) HasMultiple() bool {
	return fc.Count() > 1
}

func (fc *fileContainer) Errors() []error {
	return fc.errors
}

func (fc *fileContainer) RollbackAll() IFileContainer {
	for _, iFile := range fc.files {
		if err := iFile.Rollback().Error(); err != nil {
			fc.errors = append(fc.errors, err)
		}
	}
	return fc
}

func (fc *fileContainer) StoreAll(dir string, storagePaths *[]string) IFileContainer {

	for _, iFile := range fc.files {
		var storagePath string

		if err := iFile.RandomFileName().Store(dir, &storagePath).Error(); err != nil {
			fc.errors = append(fc.errors, err)
		} else {
			if storagePaths != nil {
				*storagePaths = append(*storagePaths, storagePath)
			}
		}
	}

	return fc
}

type file struct {
	w      http.ResponseWriter
	r      *http.Request
	osFile *os.File
	*sync.Once
	newFilename   string
	FileHeader    *multipart.FileHeader
	multipartFile multipart.File
	error         error
	perm          os.FileMode
}

func (f *file) Read(writer io.Writer) IFile {
	if !f.open() {
		return f
	}

	if _, err := io.Copy(writer, f.multipartFile); err != nil {
		f.error = err
	}

	return f
}

func (f *file) Rollback() IFile {
	if !f.IsValid() || f.osFile == nil {
		return f
	}

	var name = f.osFile.Name()
	_ = f.osFile.Close()

	if stat, err := os.Stat(name); err != nil || stat.IsDir() {
		return f
	}

	if err := os.Remove(name); err != nil {
		f.error = err
	}

	return f
}

func (f *file) Store(dir string, storagePath *string) IFile {
	if !f.open() {
		return f
	}

	if f.newFilename != "" {
		dir = path.Join(dir, f.newFilename)
	} else {
		dir = path.Join(dir, f.FileHeader.Filename)
	}

	outFile, err := os.OpenFile(dir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.perm)

	if err != nil {
		f.error = err
		return f
	}

	f.osFile = outFile

	if _, err = io.Copy(outFile, f.multipartFile); err != nil {
		f.error = err
		return f
	}

	if storagePath != nil {
		*storagePath = dir
	}
	return f
}

func (f *file) open() bool {

	if !f.IsValid() || f.Once == nil {
		return false
	}

	f.Do(func() {

		if f.multipartFile, f.error = f.FileHeader.Open(); f.error != nil {
			return
		}
		go func() {
			<-f.r.Context().Done()
			if f.multipartFile != nil {
				_ = f.multipartFile.Close()
			}
			if f.osFile != nil {
				_ = f.osFile.Close()
			}
		}()
	})

	if f.multipartFile == nil {
		return false
	}

	if _, f.error = f.multipartFile.Seek(0, 0); f.error != nil {
		return false
	}

	return true
}

func (f *file) RandomFileName() IFile {
	if !f.IsValid() {
		return f
	}
	f.newFilename = randFileName(f.FileHeader.Filename)
	return f
}

func (f *file) Info() IFileInfo {
	return f
}

func (f *file) SetNewName(filename string) IFile {
	f.newFilename = filename
	return f
}

func (f *file) GetNewName() string {
	return f.newFilename
}

func (f *file) SetPrem(perm os.FileMode) IFile {
	f.perm = perm
	return f
}

func (f *file) IsValid() bool {
	return f.error == nil
}

func (f file) Error() error {
	return f.error
}

func (f *file) Size() int64 {
	if !f.IsValid() {
		return 0
	}
	return f.FileHeader.Size
}

func (f *file) Name() string {
	if !f.IsValid() {
		return ""
	}
	return f.FileHeader.Filename
}

func (f *file) ContentType() string {
	if !f.IsValid() {
		return ""
	}

	return f.FileHeader.Filename
}

func (f *file) Extension() string {
	if !f.IsValid() {
		return ""
	}
	return filepath.Ext(f.FileHeader.Filename)
}
