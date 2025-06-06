package vad

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/askasoft/pango/num"
)

// per validate construct
type validate struct {
	v              *Validate
	top            reflect.Value
	ns             []byte
	actualNs       []byte
	errs           ValidationErrors
	includeExclude map[string]struct{} // reset only if StructPartial or StructExcept are called, no need otherwise
	ffn            FilterFunc
	slflParent     reflect.Value // StructLevel & FieldLevel
	slCurrent      reflect.Value // StructLevel & FieldLevel
	flField        reflect.Value // StructLevel & FieldLevel
	cf             *cField       // StructLevel & FieldLevel
	ct             *cTag         // StructLevel & FieldLevel
	fldIsPointer   bool          // StructLevel & FieldLevel
	isPartial      bool
	hasExcludes    bool
}

// parent and current will be the same the first run of validateStruct
func (v *validate) validateStruct(parent reflect.Value, current reflect.Value, typ reflect.Type, ns []byte, structNs []byte, ct *cTag) {
	cs, ok := v.v.structCache.Get(typ)
	if !ok {
		cs = v.v.extractStructCache(current, typ.Name())
	}

	if len(ns) == 0 && len(cs.name) != 0 {
		ns = append(ns, cs.name...)
		ns = append(ns, '.')

		structNs = append(structNs, cs.name...)
		structNs = append(structNs, '.')
	}

	// ct is nil on top level struct, and structs as fields that have no tag info
	// so if nil or if not nil and the structonly tag isn't present
	if ct == nil || ct.typeof != typeStructOnly {
		var f *cField

		for i := 0; i < len(cs.fields); i++ {
			f = cs.fields[i]

			if v.isPartial {
				if v.ffn != nil {
					// used with StructFiltered
					if v.ffn(append(structNs, f.name...)) {
						continue
					}
				} else {
					// used with StructPartial & StructExcept
					_, ok = v.includeExclude[string(append(structNs, f.name...))]

					if (ok && v.hasExcludes) || (!ok && !v.hasExcludes) {
						continue
					}
				}
			}

			v.traverseField(current, current.Field(f.idx), ns, structNs, f, f.cTags)
		}
	}

	// check if any struct level validations, after all field validations already checked.
	// first iteration will have no info about nostructlevel tag, and is checked prior to
	// calling the next iteration of validateStruct called from traverseField.
	if cs.fn != nil {
		v.slflParent = parent
		v.slCurrent = current
		v.ns = ns
		v.actualNs = structNs

		cs.fn(v)
	}
}

// omitEmpty is the validation function for validating if the current field's value is not the default static value.
// check recursively if the field is a pointer.
func omitEmpty(fl FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return field.IsNil()
	default:
		return !field.IsValid() || field.Interface() == reflect.Zero(field.Type()).Interface()
	}
}

// traverseField validates any field, be it a struct or single field, ensures it's validity and passes it along to be validated via it's tag options
func (v *validate) traverseField(parent reflect.Value, current reflect.Value, ns []byte, structNs []byte, cf *cField, ct *cTag) {
	var typ reflect.Type
	var kind reflect.Kind

	current, kind, v.fldIsPointer = v.extractTypeInternal(current, false)

	switch kind {
	case reflect.Ptr, reflect.Interface, reflect.Invalid:
		if ct == nil {
			return
		}

		if ct.typeof == typeOmitEmpty || ct.typeof == typeIsEmpty {
			return
		}

		if ct.hasTag {
			if kind == reflect.Invalid {
				vns := string(append(ns, cf.altName...))
				sns := vns
				if v.v.hasTagNameFunc {
					sns = string(append(structNs, cf.name...))
				}
				v.errs = append(v.errs,
					&fieldError{
						v:              v.v,
						tag:            ct.aliasTag,
						actualTag:      ct.tag,
						ns:             vns,
						structNs:       sns,
						fieldLen:       uint8(len(cf.altName)),
						structfieldLen: uint8(len(cf.name)),
						param:          ct.param,
						kind:           kind,
						cause:          errInvalidField,
					},
				)
				return
			}

			vns := string(append(ns, cf.altName...))
			sns := vns
			if v.v.hasTagNameFunc {
				sns = string(append(structNs, cf.name...))
			}
			if !ct.runValidationWhenNil {
				v.errs = append(v.errs,
					&fieldError{
						v:              v.v,
						tag:            ct.aliasTag,
						actualTag:      ct.tag,
						ns:             vns,
						structNs:       sns,
						fieldLen:       uint8(len(cf.altName)),
						structfieldLen: uint8(len(cf.name)),
						value:          getValue(current),
						param:          ct.param,
						kind:           kind,
						typ:            current.Type(),
						cause:          errNilField,
					},
				)
				return
			}
		}

	case reflect.Struct:
		typ = current.Type()
		if typ != timeType {
			if ct != nil {
				if ct.typeof == typeStructOnly {
					goto CONTINUE
				}

				if ct.typeof == typeIsEmpty {
					// set Field Level fields
					v.slflParent = parent
					v.flField = current
					v.cf = cf
					v.ct = ct

					if err := ct.fn(v); err != nil {
						vns := string(append(ns, cf.altName...))
						sns := vns
						if v.v.hasTagNameFunc {
							sns = string(append(structNs, cf.name...))
						}

						v.errs = append(v.errs,
							&fieldError{
								v:              v.v,
								tag:            ct.aliasTag,
								actualTag:      ct.tag,
								ns:             vns,
								structNs:       sns,
								fieldLen:       uint8(len(cf.altName)),
								structfieldLen: uint8(len(cf.name)),
								value:          getValue(current),
								param:          ct.param,
								kind:           kind,
								typ:            typ,
								cause:          err,
							},
						)
						return
					}
				}

				ct = ct.next
			}

			if ct != nil && ct.typeof == typeNoStructLevel {
				return
			}

		CONTINUE:
			// if len == 0 then validating using 'Var' or 'VarWithValue'
			// Var - doesn't make much sense to do it that way, should call 'Struct', but no harm...
			// VarWithField - this allows for validating against each field within the struct against a specific value
			//                pretty handy in certain situations
			if len(cf.name) > 0 {
				ns = append(append(ns, cf.altName...), '.')
				structNs = append(append(structNs, cf.name...), '.')
			}

			v.validateStruct(parent, current, typ, ns, structNs, ct)
			return
		}
	}

	if ct == nil || !ct.hasTag {
		return
	}

	typ = current.Type()

OUTER:
	for {
		if ct == nil {
			return
		}

		switch ct.typeof {
		case typeOmitEmpty:
			// set Field Level fields
			v.slflParent = parent
			v.flField = current
			v.cf = cf
			v.ct = ct

			if omitEmpty(v) {
				return
			}

			ct = ct.next
			continue

		case typeEndKeys:
			return

		case typeDive:
			ct = ct.next

			// traverse slice or map here or panic ;)
			switch kind {
			case reflect.Slice, reflect.Array:
				reusableCF := &cField{}

				for i := 0; i < current.Len(); i++ {
					reusableCF.name = cf.name + "[" + num.Itoa(i) + "]"

					if cf.namesEqual {
						reusableCF.altName = reusableCF.name
					} else {
						reusableCF.altName = cf.altName + "[" + num.Itoa(i) + "]"
					}
					v.traverseField(parent, current.Index(i), ns, structNs, reusableCF, ct)
				}

			case reflect.Map:
				reusableCF := &cField{}

				for _, key := range current.MapKeys() {
					reusableCF.name = fmt.Sprintf("%s[%v]", cf.name, key.Interface())

					if cf.namesEqual {
						reusableCF.altName = reusableCF.name
					} else {
						reusableCF.altName = fmt.Sprintf("%s[%v]", cf.altName, key.Interface())
					}

					if ct != nil && ct.typeof == typeKeys && ct.keys != nil {
						v.traverseField(parent, key, ns, structNs, reusableCF, ct.keys)
						// can be nil when just keys being validated
						if ct.next != nil {
							v.traverseField(parent, current.MapIndex(key), ns, structNs, reusableCF, ct.next)
						}
					} else {
						v.traverseField(parent, current.MapIndex(key), ns, structNs, reusableCF, ct)
					}
				}

			default:
				// throw error, if not a slice or map then should not have gotten here
				// bad dive tag
				panic("dive error! can't dive on a non slice or map")
			}

			return

		case typeOr:
			var misc []byte
			for {
				// set Field Level fields
				v.slflParent = parent
				v.flField = current
				v.cf = cf
				v.ct = ct

				err := ct.fn(v)
				if err == nil {
					// drain rest of the 'or' values, then continue or leave
					for {
						ct = ct.next

						if ct == nil {
							return
						}

						if ct.typeof != typeOr {
							continue OUTER
						}
					}
				}

				misc = append(misc, '|')
				misc = append(misc, ct.tag...)

				if ct.hasParam {
					misc = append(misc, '=')
					misc = append(misc, ct.param...)
				}

				if ct.isBlockEnd || ct.next == nil {
					// if we get here, no valid 'or' value and no more tags
					vns := string(append(ns, cf.altName...))
					sns := vns
					if v.v.hasTagNameFunc {
						sns = string(append(structNs, cf.name...))
					}

					if ct.hasAlias {
						v.errs = append(v.errs, &fieldError{
							v:              v.v,
							tag:            ct.aliasTag,
							actualTag:      ct.actualAliasTag,
							ns:             vns,
							structNs:       sns,
							fieldLen:       uint8(len(cf.altName)),
							structfieldLen: uint8(len(cf.name)),
							value:          getValue(current),
							param:          ct.param,
							kind:           kind,
							typ:            typ,
							cause:          err,
						})
					} else {
						tVal := string(misc[1:])

						v.errs = append(v.errs, &fieldError{
							v:              v.v,
							tag:            tVal,
							actualTag:      tVal,
							ns:             vns,
							structNs:       sns,
							fieldLen:       uint8(len(cf.altName)),
							structfieldLen: uint8(len(cf.name)),
							value:          getValue(current),
							param:          ct.param,
							kind:           kind,
							typ:            typ,
							cause:          err,
						})
					}
					return
				}

				ct = ct.next
			}
		default:
			// set Field Level fields
			v.slflParent = parent
			v.flField = current
			v.cf = cf
			v.ct = ct

			if err := ct.fn(v); err != nil {
				vns := string(append(ns, cf.altName...))
				sns := vns
				if v.v.hasTagNameFunc {
					sns = string(append(structNs, cf.name...))
				}

				v.errs = append(v.errs, &fieldError{
					v:              v.v,
					tag:            ct.aliasTag,
					actualTag:      ct.tag,
					ns:             vns,
					structNs:       sns,
					fieldLen:       uint8(len(cf.altName)),
					structfieldLen: uint8(len(cf.name)),
					value:          getValue(current),
					param:          ct.param,
					kind:           kind,
					typ:            typ,
					cause:          err,
				})
				return
			}
			ct = ct.next
		}
	}
}

func getValue(val reflect.Value) interface{} {
	if val.CanInterface() {
		return val.Interface()
	}

	if val.CanAddr() {
		return reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr())).Elem().Interface()
	}

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return val.Uint()
	case reflect.Complex64, reflect.Complex128:
		return val.Complex()
	case reflect.Float32, reflect.Float64:
		return val.Float()
	default:
		return val.String()
	}
}
