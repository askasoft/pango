package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pandafw/pango/ini"
	"github.com/pandafw/pango/ref"
	"github.com/pandafw/pango/str"
)

// Config config log by configuration file
func (log *Log) Config(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == ".json" || ext == ".js" {
		return log.configJSON(filename)
	}
	return log.configINI(filename)
}

func (log *Log) configJSON(filename string) error {
	fp, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fp.Close()

	c := make(map[string]interface{})
	jd := json.NewDecoder(fp)
	err = jd.Decode(&c)
	if err != nil {
		return err
	}

	async := 0
	if async, err = log.configGetIntValue(c, "async"); err != nil {
		return err
	}
	if err := log.configLogFormat(c); err != nil {
		return err
	}

	if lvl, ok := c["level"]; ok {
		switch lvls := lvl.(type) {
		case string:
			log.SetLevel(ParseLevel(lvls))
		case map[string]interface{}:
			if err := log.configLogLevels(lvls); err != nil {
				return err
			}
		}
	}

	if v, ok := c["writer"]; ok {
		if a, ok := v.([]interface{}); ok {
			if err = log.configLogWriter(a, async); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Invalid writer configuration: %v", v)
		}
	} else {
		return errors.New("Missing writer configuration")
	}
	return nil
}

func (log *Log) configINI(filename string) (err error) {
	ini := ini.NewIni()
	if err = ini.LoadFile(filename); err != nil {
		return err
	}

	c := ini.Section("").Map()

	async := 0
	if async, err = log.configGetIntValue(c, "async"); err != nil {
		return err
	}
	if err = log.configLogFormat(c); err != nil {
		return err
	}

	sec := ini.Section("level")
	if sec != nil {
		lvls := sec.Map()
		if err = log.configLogLevels(lvls); err != nil {
			return err
		}
	}

	if v, ok := c["writer"]; ok {
		if s, ok := v.(string); ok {
			ss := str.FieldsAny(s, " ,")
			a := make([]interface{}, len(ss))
			for i, w := range ss {
				var es map[string]interface{}

				sec := ini.Section("writer." + w)
				if sec == nil {
					es = make(map[string]interface{}, 1)
				} else {
					es = sec.Map()
				}

				if _, ok := es["_"]; !ok {
					es["_"] = w
				}
				a[i] = es
			}
			if err = log.configLogWriter(a, async); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Invalid writer configuration: %v", v)
		}
	} else {
		return fmt.Errorf("Missing writer configuration")
	}
	return nil
}

func (log *Log) configGetIntValue(m map[string]interface{}, k string) (int, error) {
	if v, ok := m[k]; ok {
		if v != nil {
			n, err := ref.Convert(v, reflect.TypeOf(int(0)))
			if err != nil {
				return 0, fmt.Errorf("Invalid %s value %v: %s", k, v, err.Error())
			}
			return n.(int), nil
		}
	}
	return 0, nil
}

func (log *Log) configLogFormat(m map[string]interface{}) error {
	if v, ok := m["format"]; ok {
		if s, ok := v.(string); ok {
			log.SetFormatter(NewLogFormatter(s))
		} else {
			return fmt.Errorf("Invalid format value: %v", v)
		}
	}
	return nil
}

func (log *Log) configLogLevels(lls map[string]interface{}) error {
	lvls := map[string]Level{}

	for k, v := range lls {
		if s, ok := v.(string); ok {
			if k == "*" {
				log.SetLevel(ParseLevel(s))
			} else {
				lvl := ParseLevel(s)
				if lvl != LevelNone {
					lvls[k] = lvl
				}
			}
		} else {
			return fmt.Errorf("Invalid level %v", v)
		}
	}

	log.levels = lvls
	return nil
}

func (log *Log) configLogWriter(a []interface{}, async int) (err error) {
	var ws []Writer
	for _, i := range a {
		if c, ok := i.(map[string]interface{}); ok {
			if n, ok := c["_"]; ok {
				w := CreateWriter(n.(string))
				if w == nil {
					return fmt.Errorf("Invalid writer name: %v", n)
				}
				if err = ConfigWriter(w, c); err != nil {
					return err
				}

				a := 0
				if a, err = log.configGetIntValue(c, "_async"); err != nil {
					return err
				}
				if a > 0 {
					w = NewAsyncWriter(w, a)
				}
				ws = append(ws, w)
			} else {
				return fmt.Errorf("Missing writer type: %v", c)
			}
		} else {
			return fmt.Errorf("Invalid writer item: %v", i)
		}
	}
	if len(ws) < 0 {
		return fmt.Errorf("Empty writer configuration: %v", a)
	}

	var lw Writer
	if len(ws) == 1 {
		lw = ws[0]
	} else {
		lw = NewMultiWriter(ws...)
	}

	if async > 0 {
		if _, ok := lw.(AsyncWriter); !ok {
			lw = NewAsyncWriter(lw, async)
		}
	}

	log.SetWriter(lw)
	return nil
}
