package request

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type file struct {
	FileHeader    *multipart.FileHeader
	MultipartFile multipart.File
	File          *os.File
	Error         error
	FullPath      string
	Content       string
}

type IFile interface {
	Open() IFile
	Store(dir string) IFile
	StoreAndClose(dir string) IFile
	GetFullPath() string
	Remove() IFile
	Close() IFile
	GetError() error
	Size() int64
	Name() string
	ContentType() string
	Extension() string
}

func (f *file) Remove() IFile {

	if f.File == nil {
		return f
	}

	name := f.File.Name()

	f.Close()

	if err := os.Remove(name); err != nil {
		f.Error = err
	}

	return f
}

func (f *file) Store(dir string) IFile {
	f.Open()

	if f.MultipartFile == nil || f.File != nil {
		return f
	}

	outFile, err := os.CreateTemp(dir,
		fmt.Sprintf("*%s",
			filepath.Ext(f.FileHeader.Filename),
		),
	)

	if err != nil {
		f.Error = err
		return f
	}

	f.File = outFile

	f.FullPath = outFile.Name()

	f.MultipartFile.Seek(0, 0)

	if _, err := io.Copy(f.File, f.MultipartFile); err != nil {
		f.Error = err
		return f
	}

	return f
}

func (f *file) GetFullPath() string {
	return f.FullPath
}

func (f *file) Open() IFile {

	if f.FileHeader == nil || f.MultipartFile != nil {
		return f
	}

	f.MultipartFile, f.Error = f.FileHeader.Open()

	return f
}

func (f *file) StoreAndClose(dir string) IFile {
	return f.Store(dir).Close()
}

func (f *file) Close() IFile {

	if f.MultipartFile != nil {
		f.MultipartFile.Close()
		f.MultipartFile = nil
	}

	if f.File != nil {
		f.File.Close()
		f.File = nil
	}

	return f
}

func (f *file) GetError() error {
	return f.Error
}

func (f *file) Size() int64 {

	if f.FileHeader == nil {
		return 0
	}

	return f.FileHeader.Size
}

func (f *file) Name() string {
	if f.FileHeader == nil {
		return ""
	}

	return f.FileHeader.Filename
}

func (f *file) ContentType() string {

	f.Open()

	if f.MultipartFile == nil {
		return ""
	}

	if f.Content != "" {
		return f.Content
	}

	buf := make([]byte, 512)

	f.MultipartFile.Seek(0, 0)

	if _, err := f.MultipartFile.Read(buf); err != nil {
		f.Error = err
		return ""
	}

	f.Content = http.DetectContentType(buf)

	return f.Content
}

func (f *file) Extension() string {
	if f.FileHeader == nil {
		return ""
	}
	return filepath.Ext(f.FileHeader.Filename)
}
