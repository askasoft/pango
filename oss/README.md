# OS system statistics library for Go
This is a library to get system metrics like cpu load and memory usage.

## Example
```go
package main

import (
	"fmt"
	"os"

	"github.com/askasoft/pango/oss/mem"
)

func main() {
	ms, err := mem.GetMemoryStats()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	fmt.Printf("memory total: %d bytes\n", ms.Total)
	fmt.Printf("memory used: %d bytes\n", ms.Used())
	fmt.Printf("memory buffer: %d bytes\n", ms.Buffer)
	fmt.Printf("memory cached: %d bytes\n", ms.Cached)
	fmt.Printf("memory free: %d bytes\n", ms.Free)
}
```

## Supported OS

||loadavg|uptime|cpu|memory|network|disk i/o|
|:--:|:--:|:--:|:--:|:--:|:--:|:--:|
|Linux|yes|yes|yes|yes|yes|yes|
|Darwin|yes|yes|*1|yes|yes|no|
|FreeBSD|yes|yes|no|yes|yes|no|
|NetBSD|yes|yes|no|no|yes|no|
|OpenBSD|yes|yes|no|no|no|no|
|Windows|no|yes|no|yes|no|no|

*1: unavailable without cgo

## Note for counter values
This library returns the counter value for cpu, network and disk I/O statistics by design. 
To get the cpu usage in percent, network traffic in kB/s or disk IOPS, sleep for a while and calculate the difference.

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/askasoft/oss/cpu"
)

func main() {
	cs, err := cpu.GetCPUStatsDelta(time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	fmt.Printf("cpu user: %f %%\n", cpu.UserUsage()*100)
	fmt.Printf("cpu system: %f %%\n", cpu.SystemUsage()*100)
	fmt.Printf("cpu idle: %f %%\n", cpu.IdleUsage()*100)
}
```

