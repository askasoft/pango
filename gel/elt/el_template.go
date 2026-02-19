package elt

import (
	"fmt"
	"io"
	"strings"

	"github.com/askasoft/pango/cas"
	"github.com/askasoft/pango/gel"
	"github.com/askasoft/pango/str"
)

type ELTemplate struct {
	template string
	segments []any
}

func Parse(template string, delims ...string) (*ELTemplate, error) {
	elt := &ELTemplate{
		template: template,
	}

	prefix := "{{"
	suffix := "}}"

	switch len(delims) {
	case 1:
		prefix = str.IfEmpty(delims[0], prefix)
		suffix = prefix
	case 2:
		prefix = str.IfEmpty(delims[0], prefix)
		suffix = str.IfEmpty(delims[1], suffix)
	}

	if err := elt.parse(prefix, suffix); err != nil {
		return nil, err
	}
	return elt, nil
}

func (elt *ELTemplate) addStr(s string) {
	if len(s) == 0 {
		return
	}
	elt.segments = append(elt.segments, s)
}

func (elt *ELTemplate) addEL(el *gel.EL) {
	elt.segments = append(elt.segments, el)
}

func (elt *ELTemplate) parse(prefix, suffix string) (err error) {
	if elt.template == "" {
		return
	}

	s := elt.template

	xp := strings.Index(s, prefix)
	if xp < 0 {
		elt.addStr(s)
		return
	}
	pz := len(prefix)

	xs := strings.Index(s[xp+pz:], suffix)
	if xs < 0 {
		elt.addStr(s)
		return
	}
	sz := len(suffix)

	xn, difffx := 0, prefix != suffix

	if difffx {
		// find next prefix, ignore continuous prefix like `{{ outer {{inner}} }}`
		xn = strings.Index(s[xp+pz:], prefix)
		if xn >= 0 && xn < xs {
			xn += pz
			xs -= xn
			xp += xn
			xn = xs
		}
	}

	cache := map[string]*gel.EL{}

	for {
		elt.addStr(s[:xp])
		s = s[xp+pz:]

		{
			expr := str.Strip(s[:xs])
			if expr != "" {
				el, ok := cache[expr]
				if !ok {
					el, err = gel.Compile(expr)
					if err != nil {
						return
					}
					cache[expr] = el
				}
				elt.addEL(el)
			}
		}
		xs += sz
		s = s[xs:]

		switch {
		case xn < 0:
			xp = xn
		case xn >= xs:
			xp = xn - xs
		default:
			xp = strings.Index(s, prefix)
		}

		if xp < 0 {
			elt.addStr(s)
			break
		}

		xs = strings.Index(s[xp+pz:], suffix)
		if xs < 0 {
			elt.addStr(s)
			break
		}

		if difffx {
			// find next prefix, ignore continuous prefix like `{{ outer {{inner}} }}`
			xn = strings.Index(s[xp+pz:], prefix)
			if xn >= 0 && xn < xs {
				xn += pz
				xs -= xn
				xp += xn
				xn = xs
			}
		}
	}
	return
}

func (elt *ELTemplate) String() string {
	return elt.template
}

func (elt *ELTemplate) wrapError(err error) error {
	return fmt.Errorf("elt: %w", err)
}

// Evaluate applies a parsed template to the specified data object,
// and returns the evaluated string.
func (elt *ELTemplate) Evaluate(data any, strict ...bool) (string, error) {
	sb := &strings.Builder{}
	err := elt.Execute(sb, data, strict...)
	return sb.String(), err
}

// Execute applies a parsed template to the specified data object,
// and writes the output to w.
// If an error occurs executing the template or writing its output,
// execution stops, but partial results may already have been written to
// the output writer.
// A template may be executed safely in parallel, although if parallel
// executions share a Writer the output may be interleaved.
func (elt *ELTemplate) Execute(w io.Writer, data any, strict ...bool) error {
	for _, s := range elt.segments {
		switch a := s.(type) {
		case *gel.EL:
			val, err := a.Calculate(data, strict...)
			if err != nil {
				return elt.wrapError(err)
			}

			str, err := cas.ToString(val)
			if err != nil {
				return elt.wrapError(err)
			}

			if _, err := io.WriteString(w, str); err != nil {
				return elt.wrapError(err)
			}
		case string:
			if _, err := io.WriteString(w, a); err != nil {
				return elt.wrapError(err)
			}
		}
	}
	return nil
}

// Evaluate parse and apply the template to the specified data object,
// and returns the evaluated string.
func Evaluate(template string, data any, strict ...bool) (string, error) {
	elt, err := Parse(template)
	if err != nil {
		return "", err
	}
	return elt.Evaluate(data, strict...)
}

// Execute applies a parsed template to the specified data object,
// and writes the output to w.
// If an error occurs executing the template or writing its output,
// execution stops, but partial results may already have been written to
// the output writer.
// A template may be executed safely in parallel, although if parallel
// executions share a Writer the output may be interleaved.
func Execute(w io.Writer, template string, data any, strict ...bool) error {
	elt, err := Parse(template)
	if err != nil {
		return err
	}
	return elt.Execute(w, data, strict...)
}
