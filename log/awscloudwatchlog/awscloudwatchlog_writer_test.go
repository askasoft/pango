package awscloudwatchlog

import (
	"os"
	"testing"

	"github.com/askasoft/pango/log"
)

func TestAwsCloudWatchLogWriter(t *testing.T) {
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	if awsAccessKey == "" {
		t.Skip("AWS_ACCESS_KEY not set")
		return
	}

	awsSecretKey := os.Getenv("AWS_SECRET_KEY")
	if awsSecretKey == "" {
		t.Skip("AWS_SECRET_KEY not set")
		return
	}

	lg := log.NewLog()
	lg.SetLevel(log.LevelTrace)
	lg.SetProp("HOST", "localhost")
	lg.SetProp("VERSION", "1.0")

	aw := &AWSCloudWatchLogWriter{
		AccessKey:     awsAccessKey,
		SecretKey:     awsSecretKey,
		Region:        "ap-northeast-1",
		LogGroupName:  "testloggroup",
		LogStreamName: "testlogstream",
	}
	aw.SetFormat(`json:{"time": %t{2006-01-02T15:04:05.000Z07:00}, "level": %l, "host":%x{HOST}, "version":%x{VERSON}, "name": %c, "file": %S, "line": %L, "func": %F, "msg": %m, "trace": %T}%n`)

	aw.Filter = log.NewLevelFilter(log.LevelDebug)
	lg.SetWriter(log.NewMultiWriter(
		aw,
		&log.StreamWriter{Color: true},
	))

	lg.Trace("This is a AwsCloudWatchLogWriter trace log")
	lg.Debug("This is a AwsCloudWatchLogWriter debug log")
	lg.Info("This is a AwsCloudWatchLogWriter info log")
	lg.Warn("This is a AwsCloudWatchLogWriter warn log")
	lg.Error("This is a AwsCloudWatchLogWriter error log")
	lg.Fatal("This is a AwsCloudWatchLogWriter fatal log")

	lg.Close()
}
