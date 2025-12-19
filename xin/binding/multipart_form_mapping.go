package binding

import (
	"errors"
	"mime/multipart"
	"net/http"
	"path"
	"reflect"
)

type multipartRequest http.Request

var (
	// ErrMultiFileHeader multipart.FileHeader invalid
	ErrMultiFileHeader = errors.New("unsupported field type for multipart.FileHeader")

	// ErrMultiFileHeaderLenInvalid array for []*multipart.FileHeader len invalid
	ErrMultiFileHeaderLenInvalid = errors.New("unsupported len of array for []*multipart.FileHeader")
)

// TrySet tries to set a value by the multipart request with the binding a form file
func (r *multipartRequest) TrySet(field reflect.Value, key string, opt options) (bool, *FieldBindError) {
	if files := r.MultipartForm.File[key]; len(files) != 0 {
		ok, err := setByMultipartFormFile(field, files)
		if err != nil {
			be := &FieldBindError{
				Err:   err,
				Field: key,
			}
			for _, f := range files {
				be.Values = append(be.Values, path.Base(f.Filename))
			}
			return ok, be
		}
		return ok, nil
	}

	return setByForm(field, r.MultipartForm.Value, key, opt)
}

func setByMultipartFormFile(field reflect.Value, files []*multipart.FileHeader) (isSet bool, err error) {
	switch field.Kind() {
	case reflect.Pointer:
		switch field.Interface().(type) {
		case *multipart.FileHeader:
			field.Set(reflect.ValueOf(files[0]))
			return true, nil
		}
	case reflect.Struct:
		switch field.Interface().(type) {
		case multipart.FileHeader:
			field.Set(reflect.ValueOf(*files[0]))
			return true, nil
		}
	case reflect.Slice:
		slice := reflect.MakeSlice(field.Type(), len(files), len(files))
		isSet, err = setArrayOfMultipartFormFiles(slice, files)
		if err != nil || !isSet {
			return isSet, err
		}
		field.Set(slice)
		return true, nil
	case reflect.Array:
		return setArrayOfMultipartFormFiles(field, files)
	}
	return false, ErrMultiFileHeader
}

func setArrayOfMultipartFormFiles(field reflect.Value, files []*multipart.FileHeader) (isSet bool, err error) {
	if field.Len() != len(files) {
		return false, ErrMultiFileHeaderLenInvalid
	}
	for i := range files {
		set, err := setByMultipartFormFile(field.Index(i), files[i:i+1])
		if err != nil || !set {
			return set, err
		}
	}
	return true, nil
}
