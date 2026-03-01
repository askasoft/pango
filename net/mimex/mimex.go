package mimex

import (
	"mime"
	"path/filepath"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/str"
)

// MediaTypeByFilename returns the MIME type (no charset) associated with the extension of filename.
func MediaTypeByFilename(filename string, defaults ...string) string {
	return str.IfEmpty(str.SubstrBeforeByte(ContentTypeByFilename(filename), ';'), asg.First(defaults))
}

// MediaTypeByExtension returns the MIME type (no charset) associated with the file extension ext.
func MediaTypeByExtension(ext string, defaults ...string) string {
	return str.IfEmpty(str.SubstrBeforeByte(ContentTypeByExtension(ext), ';'), asg.First(defaults))
}

// ContextTypeByFilename returns the MIME type associated with the extension of filename.
func ContentTypeByFilename(filename string, defaults ...string) string {
	return ContentTypeByExtension(filepath.Ext(filename), defaults...)
}

// ContextTypeByExtension returns the MIME type associated with the file extension ext.
// The extension ext should begin with a leading dot, as in ".html".
// When ext has no associated type, TypeByExtension returns "".
//
// Extensions are looked up first case-sensitively, then case-insensitively.
//
// The built-in table is small but on unix it is augmented by the local
// system's MIME-info database or mime.types file(s) if available under one or
// more of these names:
//
//	/usr/local/share/mime/globs2
//	/usr/share/mime/globs2
//	/etc/mime.types
//	/etc/apache2/mime.types
//	/etc/apache/mime.types
//
// On Windows, MIME types are extracted from the registry.
//
// Text types have the charset parameter set to "utf-8" by default.
func ContentTypeByExtension(ext string, defaults ...string) string {
	return str.IfEmpty(mime.TypeByExtension(ext), asg.First(defaults))
}

// ExtensionsByType returns the extensions known to be associated with the MIME
// type typ. The returned extensions will each begin with a leading dot, as in
// ".html". When typ has no associated extensions, ExtensionsByType returns an
// nil slice.
func ExtensionsByType(typ string) ([]string, error) {
	return mime.ExtensionsByType(typ)
}

// AddExtensionType sets the MIME type associated with
// the extension ext to typ. The extension should begin with
// a leading dot, as in ".html".
func AddExtensionType(ext, typ string) error {
	return mime.AddExtensionType(ext, typ)
}

// ParseMediaType parses a media type value and any optional
// parameters, per RFC 1521.  Media types are the values in
// Content-Type and Content-Disposition headers (RFC 2183).
// On success, ParseMediaType returns the media type converted
// to lowercase and trimmed of white space and a non-nil map.
// If there is an error parsing the optional parameter,
// the media type will be returned along with the error
// [ErrInvalidMediaParameter].
// The returned map, params, maps from the lowercase
// attribute to the attribute value with its case preserved.
func ParseMediaType(v string) (mediatype string, params map[string]string, err error) {
	return mime.ParseMediaType(v)
}
