package request

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"reflect"

	"github.com/fobus1289/marshrudka/util"

	"github.com/gorilla/schema"
)

func (r *request) Json(in any) error {

	if !util.IsEquals(in, reflect.Struct, reflect.Map, reflect.Slice, reflect.Array) {
		return errors.New("type not suported, need Struct | Map | Slice | Array")
	}

	return json.NewDecoder(r.Request.Body).Decode(in)
}

func (r *request) Xml(in any) error {

	if !util.IsEquals(in, reflect.Struct, reflect.Map, reflect.Slice, reflect.Array) {
		return errors.New("type not suported, need Struct | Map | Slice | Array")
	}

	return xml.NewDecoder(r.Request.Body).Decode(in)
}

func (r *request) FormData(in any) error {

	if !util.IsEquals(in, reflect.Struct) {
		return errors.New("type not suported, need Struct")
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return err
	}

	form := r.PostForm

	if len(form) == 0 {
		return errors.New("form data is empty")
	}

	return schema.NewDecoder().Decode(in, form)
}
