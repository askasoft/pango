package urlx

import (
	"io"
	"net/url"
	"strings"
)

// CleanURL returns the URL without user, query and fragment.
func CleanURL(link string) (string, error) {
	ulink, err := url.Parse(link)
	if err != nil {
		return link, err
	}

	ulink.ForceQuery = false
	ulink.RawQuery = ""
	ulink.Fragment = ""
	ulink.User = nil

	return ulink.String(), nil
}

// ParentURL returns the parent URL relative to the `link`.
func ParentURL(link string) (string, error) {
	ulink, err := url.Parse(link)
	if err != nil {
		return link, err
	}

	ulink.ForceQuery = false
	ulink.RawQuery = ""
	ulink.Fragment = ""
	ulink.User = nil
	if ulink.Path == "" {
		ulink.Path = "/"
		return ulink.String(), nil
	}

	udir, err := url.Parse("./")
	if err != nil {
		return link, err
	}

	ubase := ulink.ResolveReference(udir)
	return ubase.String(), nil
}

func EncodeQuery(kvs ...string) string {
	if len(kvs) == 0 {
		return ""
	}

	sb := &strings.Builder{}
	_ = sb.WriteByte('?')
	_ = AppendValues(sb, kvs...)
	return sb.String()
}

func EncodeValues(kvs ...string) string {
	if len(kvs) == 0 {
		return ""
	}

	sb := &strings.Builder{}
	_ = AppendValues(sb, kvs...)
	return sb.String()
}

func AppendValues(w io.Writer, kvs ...string) error {
	z := len(kvs)
	if z == 0 {
		return nil
	}

	for i := 0; i < z; i++ {
		if i > 0 {
			if _, err := w.Write([]byte{'&'}); err != nil {
				return err
			}
		}

		if _, err := io.WriteString(w, url.QueryEscape(kvs[i])); err != nil {
			return err
		}
		if _, err := w.Write([]byte{'='}); err != nil {
			return err
		}
		i++
		if i < z {
			if _, err := io.WriteString(w, url.QueryEscape(kvs[i])); err != nil {
				return err
			}
		}
	}
	return nil
}
