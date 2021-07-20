## log
log is a Go log manager. It can use many log writers. This package is inspired by https://github.com/pandafw/panda/tree/master/panda-core/src/main/java/panda/log .


## How to install?

	go get github.com/pandafw/pango/log


## What writers are supported?

As of now this log support stream(console), file, slack, smtp, connection(tcp), webhook.


## How to use it?

First you must import it

```golang
import (
	"github.com/pandafw/pango/log"
)
```

Then init a Log (example with console writer)

```golang
log := log.NewLog()
log.SetWriter(&log.StreamWriter{Color:true})
```

Use it like this:

```golang
log.Trace("trace")
log.Debug("debug")
log.Info("info")
log.Warn("warning")
log.Fatal("fatal")
```

### File writer

Configure file writer like this:

```golang
log := log.NewLog()
log.SetWriter(&log.FileWriter{Path:"test.log"})
```

### Conn writer

Configure like this:

```golang
log := log.NewLog()
log.SetWriter(&log.ConnWriter{Net:"tcp",Addr:":7020"})
log.Info("info")
```

### Slack writer

Configure like this:

```golang
log := log.NewLog()
log.SetWriter(&log.SlackWriter{
	Webhook: "https://hooks.slack.com/services/...",
	Channel: "alert",
	Username: "gotest",
)
log.Error("error")
```

### SMTP writer

Configure like this:

```golang
log := log.NewLog()
log.SetWriter(&log.SMTPWriter{
	Host: "smtp.gmail.com",
	Port: 587,
	Username: "pangotest@gmail.com",
	Password: "xxxxxxxx",
	From: "xxxx@gmail.com",
	Tos: []string{"someone@gmail.com"},
)
log.Fatal("oh my god!")
```

### Webhook writer

Configure like this:

```golang
log := log.NewLog()
log.SetWriter(&log.WebhookWriter{
	Webhook: "http://localhost:9200/pango/logs",
	ContentType: "application/json",
	Timeout: time.Second*5,
)
log.Fatal("fatal error!")
```


## Configure from ini file
```golang
log := log.NewLog()
log.Config("log.ini")
```

log.ini
```ini
# log configuration #

### log async ###
async = 1000

### log format ###
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
subject = %l - %m 
channel = develop
username = gotest
webhook = https://hooks.slack.com/services/...
timeout = 5s
format = %l - %m%n%T
filter = level:error

### smtp writer ###
[writer.smtp]
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

### webhook writer ###
[writer.webhook]
webhook = http://localhost:9200/pango/logs
contentType = application/json
timeout = 5s
format = json:{"when":%t{2006-01-02T15:04:05.000Z07:00}, "level":%l, "file":%S, "line":%L, "func":%F, "msg": %m, "stack": %T}%n
filter = level:error
```
