package awscloudwatchlog

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/log/internal"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

// AWSCloudWatchLogWriter implements log.Writer interface and batch send log events to aws cloud watch logs.
type AWSCloudWatchLogWriter struct {
	log.BatchSupport
	log.FilterSupport
	log.FormatSupport

	AccessKey     string // optional
	SecretKey     string // optional
	Region        string // required
	LogGroupName  string // required
	LogStreamName string // optional: if empty, use instance-id for log stream name

	config *aws.Config
	client *cloudwatchlogs.Client
	group  bool
	stream bool
}

// Write cache log message, flush if needed
func (aw *AWSCloudWatchLogWriter) Write(le *log.Event) {
	if aw.Reject(le) {
		le = nil
	}

	aw.BatchWrite(le, aw.flush)
}

// Flush flush cached events
func (aw *AWSCloudWatchLogWriter) Flush() {
	aw.BatchFlush(aw.flush)
}

// Close flush and close the writer
func (aw *AWSCloudWatchLogWriter) Close() {
	aw.Flush()
}

func (aw *AWSCloudWatchLogWriter) flush(eb *log.EventBuffer) error {
	if err := aw.init(); err != nil {
		return err
	}

	plei := &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  &aw.LogGroupName,
		LogStreamName: &aw.LogStreamName,
	}

	for it := eb.Iterator(); it.Next(); {
		le := it.Value()
		tsp := le.Time.UnixMilli()
		msg := string(aw.Format(le, log.JSONFmtDefault))

		ile := types.InputLogEvent{
			Timestamp: &tsp,
			Message:   &msg,
		}
		plei.LogEvents = append(plei.LogEvents, ile)
	}

	res, err := aw.client.PutLogEvents(context.TODO(), plei, aw.optRegion)
	if err != nil {
		return err
	}

	if res.RejectedLogEventsInfo != nil {
		// The expired log events.
		expiredLogEventEndIndex := int32(-1)
		if res.RejectedLogEventsInfo.ExpiredLogEventEndIndex != nil {
			expiredLogEventEndIndex = *res.RejectedLogEventsInfo.ExpiredLogEventEndIndex
		}

		// The index of the first log event that is too new. This field is inclusive.
		tooNewLogEventStartIndex := int32(-1)
		if res.RejectedLogEventsInfo.TooNewLogEventStartIndex != nil {
			tooNewLogEventStartIndex = *res.RejectedLogEventsInfo.TooNewLogEventStartIndex
		}

		// The index of the last log event that is too old. This field is exclusive.
		tooOldLogEventEndIndex := int32(-1)
		if res.RejectedLogEventsInfo.TooOldLogEventEndIndex != nil {
			tooOldLogEventEndIndex = *res.RejectedLogEventsInfo.TooOldLogEventEndIndex
		}

		internal.Perrorf("awscloudwatchlog: Rejected: %d %d %d", expiredLogEventEndIndex, tooNewLogEventStartIndex, tooOldLogEventEndIndex)
	}

	return nil
}

func (aw *AWSCloudWatchLogWriter) optRegion(op *cloudwatchlogs.Options) {
	op.Region = aw.Region
}

func (aw *AWSCloudWatchLogWriter) init() error {
	if err := aw.initAwsClient(); err != nil {
		return err
	}
	if err := aw.createLogGroup(); err != nil {
		return err
	}
	return aw.createLogStream()
}

func (aw *AWSCloudWatchLogWriter) initAwsClient() error {
	if aw.config != nil {
		return nil
	}

	var optFns []func(*config.LoadOptions) error
	if aw.AccessKey != "" && aw.SecretKey != "" {
		optFn := config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(aw.AccessKey, aw.SecretKey, ""))
		optFns = append(optFns, optFn)
	}
	if aw.Region != "" {
		optFn := config.WithDefaultRegion(aw.Region)
		optFns = append(optFns, optFn)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), optFns...)
	if err != nil {
		return err
	}

	aw.config = &cfg

	if aw.LogStreamName == "" {
		aw.LogStreamName = aw.getInstanceID()
	}

	aw.client = cloudwatchlogs.NewFromConfig(*aw.config)
	return nil
}

func (aw *AWSCloudWatchLogWriter) createLogGroup() error {
	if aw.group {
		return nil
	}

	clgi := &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: &aw.LogGroupName,
	}

	_, err := aw.client.CreateLogGroup(context.TODO(), clgi, aw.optRegion)
	if err == nil || isAwsResourceAlreadyExistsError(err) {
		aw.group = true
		return nil
	}

	return fmt.Errorf("awscloudwatchlog: CreateLogGroup(%q): %w", aw.LogGroupName, err)
}

func (aw *AWSCloudWatchLogWriter) createLogStream() error {
	if aw.stream {
		return nil
	}

	clsi := &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  &aw.LogGroupName,
		LogStreamName: &aw.LogStreamName,
	}

	_, err := aw.client.CreateLogStream(context.TODO(), clsi, aw.optRegion)
	if err == nil || isAwsResourceAlreadyExistsError(err) {
		aw.stream = true
		return nil
	}

	return fmt.Errorf("awscloudwatchlog: CreateLogStream(%q, %q): %w", aw.LogGroupName, aw.LogStreamName, err)
}

func isAwsResourceAlreadyExistsError(err error) bool {
	var raee *types.ResourceAlreadyExistsException
	return errors.As(err, &raee)
}

func (aw *AWSCloudWatchLogWriter) getMetadata(client *imds.Client, path string) string {
	res, err := client.GetMetadata(context.TODO(), &imds.GetMetadataInput{
		Path: path,
	})
	if err == nil {
		defer res.Content.Close()

		if bs, err := io.ReadAll(res.Content); err == nil {
			return string(bs)
		}
	}
	return "unknown"
}

func (aw *AWSCloudWatchLogWriter) getInstanceID() string {
	client := imds.NewFromConfig(*aw.config)
	return aw.getMetadata(client, "instance-id")
}

func init() {
	log.RegisterWriter("awscloudwatchlog", func() log.Writer {
		return &AWSCloudWatchLogWriter{}
	})
}
