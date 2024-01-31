package xin

import (
	"encoding/xml"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/askasoft/pango/net/httpx"
)

// BindKey indicates a default bind key.
const BindKey = "XIN_BIND_KEY"

// MustBind is a helper function for given interface object and returns a Xin middleware.
func MustBind(val any) HandlerFunc {
	value := reflect.ValueOf(val)
	if value.Kind() == reflect.Ptr {
		panic(`Bind struct can not be a pointer. Use: xin.MustBind(Struct{}) instead of xin.MustBind(&Struct{})`)
	}

	type_ := value.Type()
	return func(c *Context) {
		obj := reflect.New(type_).Interface()
		if c.MustBind(obj) == nil {
			c.Set(BindKey, obj)
		}
	}
}

// WrapF is a helper function for wrapping http.HandlerFunc and returns a Xin middleware.
func WrapF(f http.HandlerFunc) HandlerFunc {
	return func(c *Context) {
		f(c.Writer, c.Request)
	}
}

// WrapH is a helper function for wrapping http.Handler and returns a Xin middleware.
func WrapH(h http.Handler) HandlerFunc {
	return func(c *Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// H is a shortcut for map[string]any
type H map[string]any

// MarshalXML allows type H to be used with xml.Marshal.
func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "map",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range h {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

// SaveUploadedFile uploads the form file to specific dst.
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	dir := path.Dir(dst)
	if err := os.MkdirAll(dir, os.FileMode(0770)); err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// Dir returns a http.FileSystem that can be used by http.FileServer().
// if browsable == true, then it works the same as http.Dir() otherwise it returns
// a filesystem that prevents http.FileServer() to list the directory files.
func Dir(root string, browsable ...bool) http.FileSystem {
	return httpx.Dir(root, browsable...)
}

// FS returns a http.FileSystem that can be used by http.FileServer().
// if browsable == true, then it works the same as http.FS() otherwise it returns
// a filesystem that prevents http.FileServer() to list the directory files.
func FS(fsys fs.FS, browsable ...bool) http.FileSystem {
	return httpx.FS(fsys, browsable...)
}

// FixedModTimeFS returns a FileSystem with fixed ModTime
func FixedModTimeFS(hfs http.FileSystem, mt time.Time) http.FileSystem {
	return httpx.FixedModTimeFS(hfs, mt)
}

// --------------------------------------------------------------------
func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

func chooseData(custom, wildcard any) any {
	if custom != nil {
		return custom
	}
	if wildcard != nil {
		return wildcard
	}
	panic("negotiation config is invalid")
}

func parseAccept(acceptHeader string) []string {
	parts := strings.Split(acceptHeader, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if i := strings.IndexByte(part, ';'); i > 0 {
			part = part[:i]
		}
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func nameOfFunction(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}
