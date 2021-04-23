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

	if err := log.configLogAsync(c); err != nil {
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
			if err := log.configLogWriter(a); err != nil {
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

func (log *Log) configINI(filename string) error {
	ini := ini.NewIni()
	if err := ini.LoadFile(filename); err != nil {
		return err
	}

	c := ini.Section("").Kvmap()

	if err := log.configLogAsync(c); err != nil {
		return err
	}
	if err := log.configLogFormat(c); err != nil {
		return err
	}

	sec := ini.Section("level")
	if sec != nil {
		lvls := sec.Kvmap()
		if err := log.configLogLevels(lvls); err != nil {
			return err
		}
	}

	if v, ok := c["writer"]; ok {
		if s, ok := v.(string); ok {
			ss := str.SplitAnyNoEmpty(s, " ,")
			a := make([]interface{}, len(ss))
			for i, w := range ss {
				var es map[string]interface{}

				sec := ini.Section("writer." + w)
				if sec == nil {
					es = make(map[string]interface{}, 1)
				} else {
					es = sec.Kvmap()
				}

				if _, ok := es["_"]; !ok {
					es["_"] = w
				}
				a[i] = es
			}
			if err := log.configLogWriter(a); err != nil {
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

func (log *Log) configLogAsync(m map[string]interface{}) error {
	if v, ok := m["async"]; ok {
		if v != nil {
			n, err := ref.Convert(v, reflect.TypeOf(int(0)))
			if err != nil {
				return fmt.Errorf("Invalid async value %v: %s", v, err.Error())
			}
			log.Async(n.(int))
		}
	}
	return nil
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

func (log *Log) configLogLevels(lvls map[string]interface{}) error {
	for k, v := range lvls {
		if s, ok := v.(string); ok {
			if k == "*" {
				log.SetLevel(ParseLevel(s))
			} else {
				log.SetLoggerLevel(k, ParseLevel(s))
			}
		} else {
			return fmt.Errorf("Invalid level %v", v)
		}
	}
	return nil
}

func (log *Log) configLogWriter(a []interface{}) error {
	var ws []Writer
	for _, i := range a {
		if c, ok := i.(map[string]interface{}); ok {
			if n, ok := c["_"]; ok {
				w := CreateWriter(n.(string))
				if w == nil {
					return fmt.Errorf("Invalid writer name: %v", n)
				}
				if err := ConfigWriter(w, c); err != nil {
					return err
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

	if len(ws) == 1 {
		log.SetWriter(ws[0])
	} else {
		log.SetWriter(&MultiWriter{Writers: ws})
	}
	return nil
}
