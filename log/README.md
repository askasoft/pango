## log
log is a Go log manager. It can use many log writers. The repo is inspired by https://github.com/pandafw/panda/tree/master/panda-core/src/main/java/panda/log .


## How to install?

	go get github.com/pandafw/pango/log


## What writers are supported?

As of now this log support console, file, slack, smtp, connection.


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
log.AddWriter(&log.ConsoleWriter{})
```

Use it like this:

```golang
log.Trace("trace")
log.Debug("debug")
log.Info("info")
log.Warn("warning")
log.Fatal("fatal")
```

## File writer

Configure file writer like this:

```golang
log := log.NewLog()
log.AddWriter("file", `{"filename":"test.log"}`)
```

## Conn writer

Configure like this:

```golang
log := log.NewLog()
log.AddWriter("conn", `{"net":"tcp","addr":":7020"}`)
log.Info("info")
```

## Smtp writer

Configure like this:

```golang
log := log.NewLog()
log.AddWriter("smtp", `{"username":"pangotest@gmail.com","password":"xxxxxxxx","host":"smtp.gmail.com:587","sendTos":["someone@gmail.com"]}`)
log.Fatal("oh my god!")
```
