package request

import (
	"net/http"
	"sync"
)

func NewFormFile(w http.ResponseWriter, r *http.Request) IFormFile {
	return &formFile{
		w:     w,
		r:     r,
		once:  &sync.Once{},
		files: files{},
	}
}

type formFile struct {
	w     http.ResponseWriter
	r     *http.Request
	once  *sync.Once
	files files
}

func (f *formFile) init() {
	f.once.Do(func() {

		f.r.FormValue("")

		var multipartForm = f.r.MultipartForm

		if multipartForm == nil || len(multipartForm.File) < 1 {
			return
		}

		var multipartFormFiles = multipartForm.File

		for key, headers := range multipartFormFiles {

			if len(headers) < 1 {
				continue
			}

			var _files = make([]IFile, 0, len(headers))

			for _, header := range headers {
				_files = append(_files, &file{
					w:          f.w,
					r:          f.r,
					FileHeader: header,
					perm:       0666,
					Once:       &sync.Once{},
				})
			}

			f.files[key] = &fileContainer{
				files: _files,
			}

		}

	})
}

func (f *formFile) Get(name string) IFileContainer {
	f.init()
	return f.files.Get(name)
}

func (f *formFile) Files() IFiles {
	f.init()
	return f.files
}

func (fs files) Get(name string) IFileContainer {
	if f := fs[name]; f != nil {
		return f
	}
	return &fileContainer{}
}

func (fs files) Len() int {
	return len(fs)
}
