package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/askasoft/pango/cas"
	"github.com/askasoft/pango/ini"
	"github.com/askasoft/pango/str"
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

	c := make(map[string]any)
	jd := json.NewDecoder(fp)
	err = jd.Decode(&c)
	if err != nil {
		return err
	}

	var async int
	if async, err = log.configGetIntValue(c, "async"); err != nil {
		return err
	}

	if lvl, ok := c["level"]; ok {
		switch lvls := lvl.(type) {
		case string:
			log.SetLevel(ParseLevel(lvls))
		case map[string]any:
			if err := log.configLogLevels(lvls); err != nil {
				return err
			}
		}
	}

	if v, ok := c["writer"]; ok {
		if a, ok := v.([]any); ok {
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

	var async int
	if async, err = log.configGetIntValue(c, "async"); err != nil {
		return err
	}

	sec := ini.GetSection("level")
	if sec != nil {
		lvls := sec.Map()
		if err = log.configLogLevels(lvls); err != nil {
			return err
		}
	}

	if v, ok := c["writer"]; ok {
		if s, ok := v.(string); ok {
			ss := str.FieldsAny(s, " ,")
			a := make([]any, len(ss))
			for i, w := range ss {
				var es map[string]any

				sec := ini.GetSection("writer." + w)
				if sec == nil {
					es = make(map[string]any, 1)
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

func (log *Log) configGetIntValue(m map[string]any, k string) (int, error) {
	if v, ok := m[k]; ok {
		if v != nil {
			n, err := cas.ToInt(v)
			if err != nil {
				return 0, fmt.Errorf("Invalid %s value %v: %w", k, v, err)
			}
			return n, nil
		}
	}
	return 0, nil
}

func (log *Log) configLogLevels(lls map[string]any) error {
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

func (log *Log) configLogWriter(a []any, async int) (err error) {
	var ws []Writer
	for _, i := range a {
		if c, ok := i.(map[string]any); ok {
			if n, ok := c["_"]; ok {
				w := CreateWriter(n.(string))
				if w == nil {
					return fmt.Errorf("Invalid writer name: %v", n)
				}
				if err = ConfigWriter(w, c); err != nil {
					return err
				}

				var a int
				if a, err = log.configGetIntValue(c, "_async"); err != nil {
					return err
				}
				if a > 0 {
					w = NewAsyncWriter(w, a)
				} else if a < 0 {
					w = NewSyncWriter(w)
				}
				ws = append(ws, w)
			} else {
				return fmt.Errorf("Missing writer type: %v", c)
			}
		} else {
			return fmt.Errorf("Invalid writer item: %v", i)
		}
	}
	if len(ws) == 0 {
		return fmt.Errorf("Empty writer configuration: %v", a)
	}

	var lw Writer
	if len(ws) == 1 {
		lw = ws[0]
	} else {
		lw = NewMultiWriter(ws...)
	}

	if async > 0 {
		if _, ok := lw.(*AsyncWriter); !ok {
			lw = NewAsyncWriter(lw, async)
		}
	} else if async < 0 {
		if _, ok := lw.(*SyncWriter); !ok {
			lw = NewSyncWriter(lw)
		}
	}

	log.SwitchWriter(lw)
	return nil
}
