 Pango Log
=====================================================================

A Go Log library like Log4j. It can use many log writers. This package is inspired by https://github.com/pandafw/panda/tree/master/panda-core/src/main/java/panda/log .


### How to install?

	go get github.com/askasoft/pango


### What writers are supported?

As of now this log support stream(console), file, smtp, connection(tcp), webhook(slack, teams).


### How to use it?

First you must import it

```golang
import (
	"github.com/askasoft/pango/log"
)
```

Then init a Log (example with console writer)

```golang
	log := log.NewLog()
	log.SetWriter(log.NewSyncWriter(&log.StreamWriter{Color:true}))
```

Use it like this:

```golang
	log.Trace("trace")
	log.Debug("debug")
	log.Info("info")
	log.Warn("warning")
	log.Fatal("fatal")
```

#### File writer

Configure file writer like this:

```golang
	log := log.NewLog()
	log.SetWriter(log.NewSyncWriter(&log.FileWriter{Path:"test.log"}))
```

#### Conn writer

Configure like this:

```golang
	log := log.NewLog()
	log.SetWriter(log.NewSyncWriter(&log.ConnWriter{Net:"tcp",Addr:":7020"}))
	log.Info("info")
```

#### Slack writer

Configure like this:

```golang
	log := log.NewLog()
	log.SetWriter(log.NewSyncWriter(&log.SlackWriter{
		Webhook: "https://hooks.slack.com/services/...",
		Channel: "alert",
		Username: "gotest",
	}))
	log.Error("error")
```

#### SMTP writer

Configure like this:

```golang
	log := log.NewLog()
	log.SetWriter(log.NewSyncWriter(&log.SMTPWriter{
		Host: "smtp.gmail.com",
		Port: 587,
		Username: "pangotest@gmail.com",
		Password: "xxxxxxxx",
		From: "xxxx@gmail.com",
		Tos: []string{"someone@gmail.com"},
	}))
	log.Fatal("oh my god!")
```

#### Teams writer

Configure like this:

```golang
	log := log.NewLog()
	log.SetWriter(log.NewSyncWriter(&log.TeamsWriter{
		Webhook: "https://xxx.webhook.office.com/webhookb2/...",
	}))
	log.Error("error")
```

#### HTTP writer

Configure like this:

```golang
	log := log.NewLog()
	hw := &log.HTTPWriter{
		URL: "http://localhost:9200/pango_logs/_doc",
		ContentType: "application/json",
		Timeout: time.Second*5,
	}
	hw.SetFormat(`json:{"when": %t{2006-01-02T15:04:05.000Z07:00}, "level": %l, "name": %c, "file": %S, "line": %L, "func": %F, "msg": %m, "trace": %T}%n`)

	log.SetWriter(log.NewSyncWriter(hw))
	log.Fatal("fatal error!")
```

#### HTTP batch writer

Configure like this:

```golang
	log := log.NewLog()
	hw := &log.HTTPWriter{
		URL: "http://localhost:9200/pango_logs/_bulk",
		ContentType: "application/json",
		Timeout: time.Second*5,
		BatchWriter: BatchWriter{
			CacheCount: 6,
			BatchCount: 3,
			FlushLevel: LevelWarn,
			FlushDelta: time.Second,
		},
	}
	hw.SetFormat(`json:{"create": {}}%n{"when": %t{2006-01-02T15:04:05.000Z07:00}, "level": %l, "name": %c, "file": %S, "line": %L, "func": %F, "msg": %m, "trace": %T}%n`)

	log.SetWriter(log.NewSyncWriter(hw))
	log.Fatal("fatal error!")
```

#### Multiple writer

Configure like this:

```golang
	log := log.NewLog()
	fw := &log.FileWriter{Path:"test.log"}
	hw := &log.HTTPWriter{
		URL: "http://localhost:9200/pango_logs/_doc",
		ContentType: "application/json",
		Timeout: time.Second*5,
	}
	log.SetWriter(log.NewSyncWriter(log.NewMultiWriter(fw, hw)))
	log.Fatal("fatal error!")
```


#### Async writer

Configure like this:

```golang
	log := log.NewLog()
	hw := &log.HTTPWriter{
		URL: "http://localhost:9200/pango_logs/_doc",
		ContentType: "application/json",
		Timeout: time.Second*5,
	}
	log.SetWriter(log.NewAsyncWriter(hw, 1000))
	log.Fatal("fatal error!")
```


### Configure from ini file

```golang
	log := log.NewLog()
	log.Config("log.ini")
```

#### log.ini

```ini
# log configuration #

### log async ###
# > 0 : do asynchronize wrap 
# < 0 : do synchronize wrap
# = 0 : do nothing (default)
async = 1000

### global log format ###
#format=json:{"level":%l, "file":%S, "func":%F, "msg": %m}%n
format=text:%l %S %F() - %m%n%T

### log writer ###
writer = stdout, stderr, tcp, dailyfile, slack, smtp, webhook

### log level ###
[level]
* = info
sql = debug
http = trace

### stdout writer ###
[writer.stdout]
format = %l - %m%n%T
filter = name:out level:debug

### tcp writer ###
[writer.tcp]
_async = 1000
addr = localhost:9999
timeout = 5s
format = %l - %m%n%T
filter = level:error

### file writer ###
[writer.dailyfile]
_ = file
path = /tmp/gotest/logs/test.log
dirPerm = 0777
maxDays = 7
format = %l %S:%L %F() - %m%n%T
filter = level:error

### slack writer ###
[writer.slack]
_async = 1000
subject = %l - %m 
channel = develop
username = gotest
webhook = https://hooks.slack.com/services/...
timeout = 5s
format = %l - %m%n%T
filter = level:error

### smtp writer ###
[writer.smtp]
_async = 1000
host = localhost
port = 25
username = -----
password = xxxxxxx
from = pango@google.com
to = to1@test.com, to2@test.com
cc = cc1@test.com, cc2@test.com
timeout = 5s
subject = %l - %m 
format = %l - %m%n%T
filter = level:error

### teams writer ###
[writer.teams]
_async = 1000
webhook = https://xxx.webhook.office.com/webhookb2/...
timeout = 5s
format = %l - %m%n%T
filter = level:error

### opensearch writer ###
[writer.opensearch]
_ = http
_async = 1000
url = http://localhost:9200/pango_logs/_bulk
contentType = application/json
timeout = 5s
batchCount = 10
cacheCount = 20
flushLevel = ERROR
flushDelta = 5s
format = json:{"create": {}}%n{"when": %t{2006-01-02T15:04:05.000Z07:00}, "level": %l, "name": %c, "file": %S, "line": %L, "func": %F, "msg": %m, "trace": %T}%n
filter = level:debug
```
