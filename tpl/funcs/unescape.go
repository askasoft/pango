package funcs

import (
	"fmt"
	"html/template"
)

func HTML(a any) (template.HTML, error) {
	switch v := a.(type) {
	case template.HTML:
		return template.HTML(v), nil //nolint: gosec
	case template.HTMLAttr:
		return template.HTML(v), nil //nolint: gosec
	case template.CSS:
		return template.HTML(v), nil //nolint: gosec
	case template.JS:
		return template.HTML(v), nil //nolint: gosec
	case template.JSStr:
		return template.HTML(v), nil //nolint: gosec
	case template.URL:
		return template.HTML(v), nil //nolint: gosec
	case template.Srcset:
		return template.HTML(v), nil //nolint: gosec
	case string:
		return template.HTML(v), nil //nolint: gosec
	case []byte:
		return template.HTML(v), nil //nolint: gosec
	case rune:
		return template.HTML(v), nil //nolint: gosec
	default:
		return "", fmt.Errorf("HTML: unknown type for '%v' (%T)", v, v)
	}
}

func HTMLAttr(a any) (template.HTMLAttr, error) {
	switch v := a.(type) {
	case template.HTML:
		return template.HTMLAttr(v), nil //nolint: gosec
	case template.HTMLAttr:
		return template.HTMLAttr(v), nil //nolint: gosec
	case template.CSS:
		return template.HTMLAttr(v), nil //nolint: gosec
	case template.JS:
		return template.HTMLAttr(v), nil //nolint: gosec
	case template.JSStr:
		return template.HTMLAttr(v), nil //nolint: gosec
	case template.URL:
		return template.HTMLAttr(v), nil //nolint: gosec
	case template.Srcset:
		return template.HTMLAttr(v), nil //nolint: gosec
	case string:
		return template.HTMLAttr(v), nil //nolint: gosec
	case []byte:
		return template.HTMLAttr(v), nil //nolint: gosec
	case rune:
		return template.HTMLAttr(v), nil //nolint: gosec
	default:
		return "", fmt.Errorf("HTML: unknown type for '%v' (%T)", v, v)
	}
}

func CSS(a any) (template.CSS, error) {
	switch v := a.(type) {
	case template.HTML:
		return template.CSS(v), nil
	case template.HTMLAttr:
		return template.CSS(v), nil
	case template.CSS:
		return template.CSS(v), nil
	case template.JS:
		return template.CSS(v), nil
	case template.JSStr:
		return template.CSS(v), nil
	case template.URL:
		return template.CSS(v), nil
	case template.Srcset:
		return template.CSS(v), nil
	case string:
		return template.CSS(v), nil
	case []byte:
		return template.CSS(v), nil
	case rune:
		return template.CSS(v), nil
	default:
		return "", fmt.Errorf("CSS: unknown type for '%v' (%T)", v, v)
	}
}

func JS(a any) (template.JS, error) {
	switch v := a.(type) {
	case template.HTML:
		return template.JS(v), nil //nolint: gosec
	case template.HTMLAttr:
		return template.JS(v), nil //nolint: gosec
	case template.CSS:
		return template.JS(v), nil //nolint: gosec
	case template.JS:
		return template.JS(v), nil //nolint: gosec
	case template.JSStr:
		return template.JS(v), nil //nolint: gosec
	case template.URL:
		return template.JS(v), nil //nolint: gosec
	case template.Srcset:
		return template.JS(v), nil //nolint: gosec
	case string:
		return template.JS(v), nil //nolint: gosec
	case []byte:
		return template.JS(v), nil //nolint: gosec
	case rune:
		return template.JS(v), nil //nolint: gosec
	default:
		return "", fmt.Errorf("JS: unknown type for '%v' (%T)", v, v)
	}
}

func JSStr(a any) (template.JSStr, error) {
	switch v := a.(type) {
	case template.HTML:
		return template.JSStr(v), nil
	case template.HTMLAttr:
		return template.JSStr(v), nil
	case template.CSS:
		return template.JSStr(v), nil
	case template.JS:
		return template.JSStr(v), nil
	case template.JSStr:
		return template.JSStr(v), nil
	case template.URL:
		return template.JSStr(v), nil
	case template.Srcset:
		return template.JSStr(v), nil
	case string:
		return template.JSStr(v), nil
	case []byte:
		return template.JSStr(v), nil
	case rune:
		return template.JSStr(v), nil
	default:
		return "", fmt.Errorf("JSStr: unknown type for '%v' (%T)", v, v)
	}
}

func URL(a any) (template.URL, error) {
	switch v := a.(type) {
	case template.HTML:
		return template.URL(v), nil //nolint: gosec
	case template.HTMLAttr:
		return template.URL(v), nil //nolint: gosec
	case template.CSS:
		return template.URL(v), nil //nolint: gosec
	case template.JS:
		return template.URL(v), nil //nolint: gosec
	case template.JSStr:
		return template.URL(v), nil //nolint: gosec
	case template.URL:
		return template.URL(v), nil //nolint: gosec
	case template.Srcset:
		return template.URL(v), nil //nolint: gosec
	case string:
		return template.URL(v), nil //nolint: gosec
	case []byte:
		return template.URL(v), nil //nolint: gosec
	case rune:
		return template.URL(v), nil //nolint: gosec
	default:
		return "", fmt.Errorf("URL: unknown type for '%v' (%T)", v, v)
	}
}

func Srcset(a any) (template.Srcset, error) {
	switch v := a.(type) {
	case template.HTML:
		return template.Srcset(v), nil
	case template.HTMLAttr:
		return template.Srcset(v), nil
	case template.CSS:
		return template.Srcset(v), nil
	case template.JS:
		return template.Srcset(v), nil
	case template.JSStr:
		return template.Srcset(v), nil
	case template.URL:
		return template.Srcset(v), nil
	case template.Srcset:
		return template.Srcset(v), nil
	case string:
		return template.Srcset(v), nil
	case []byte:
		return template.Srcset(v), nil
	case rune:
		return template.Srcset(v), nil
	default:
		return "", fmt.Errorf("Srcset: unknown type for '%v' (%T)", v, v)
	}
}
