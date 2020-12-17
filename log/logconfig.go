package log

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config config log by configuration file
func Config(log *Log, file string) error {
	fp, err := os.Open(file)
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

	if err = configLogLevel(log, c); err != nil {
		return err
	}
	if err = configLogFormat(log, c); err != nil {
		return err
	}
	if err = configLogAsync(log, c); err != nil {
		return err
	}
	if err = configLogWriter(log, c); err != nil {
		return err
	}
	return nil
}

func configLogAsync(log *Log, m map[string]interface{}) error {
	if v, ok := m["async"]; ok {
		if i, ok := v.(float64); ok {
			log.Async(int(i))
		} else if i, ok := v.(int); ok {
			log.Async(i)
		} else {
			return fmt.Errorf("Invalid async value: %v", v)
		}
	}
	return nil
}

func configLogFormat(log *Log, m map[string]interface{}) error {
	if v, ok := m["format"]; ok {
		if s, ok := v.(string); ok {
			log.SetFormatter(NewLogFormatter(s))
		} else {
			return fmt.Errorf("Invalid format value: %v", v)
		}
	}
	return nil
}

func configLogLevel(log *Log, m map[string]interface{}) error {
	if lvl, ok := m["level"]; ok {
		switch lvl.(type) {
		case string:
			log.SetLevel(ParseLevel(lvl.(string)))
		case map[string]interface{}:
			for k, v := range lvl.(map[string]interface{}) {
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
		}
	}
	return nil
}

func configLogWriter(log *Log, m map[string]interface{}) error {
	if v, ok := m["writer"]; ok {
		if a, ok := v.([]interface{}); ok {
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
						return fmt.Errorf("Missing writer name: %v", v)
					}
				} else {
					return fmt.Errorf("Invalid writer item: %v", v)
				}
			}
			if len(ws) < 0 {
				return fmt.Errorf("Empty writer settings: %v", v)
			}
			if len(ws) == 1 {
				log.SetWriter(ws[0])
			} else {
				log.SetWriter(&MultiWriter{Writers: ws})
			}
			return nil
		}
		return fmt.Errorf("Invalid writer value: %v", v)
	}
	return fmt.Errorf("Missing writer settings: %v", m)
}
