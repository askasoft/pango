package teamslog

import (
	"fmt"
	"net/url"
	"time"

	"github.com/askasoft/pango/gog"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tmu"
	"github.com/askasoft/pango/whk/teams"
)

// TeamsWriter implements log Writer Interface and send log message to teams.
type TeamsWriter struct {
	log.RetrySupport
	log.FilterSupport
	log.FormatSupport
	log.SubjectSuport

	Webhook      string
	Timeout      time.Duration
	Style        string
	MaxSubLength int // default: 200
	MaxMsgLength int // default: 2000

	message teams.Message
}

// SetWebhook set the webhook URL
func (tw *TeamsWriter) SetWebhook(webhook string) error {
	_, err := url.ParseRequestURI(webhook)
	if err != nil {
		return fmt.Errorf("teamslog: invalid webhook %q: %w", webhook, err)
	}
	tw.Webhook = webhook
	return nil
}

// SetTimeout set timeout
func (tw *TeamsWriter) SetTimeout(timeout string) error {
	td, err := tmu.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("teamslog: invalid timeout %q: %w", timeout, err)
	}
	tw.Timeout = td
	return nil
}

// Write send log message to teams
func (tw *TeamsWriter) Write(le *log.Event) {
	if tw.Reject(le) {
		tw.Flush()
		return
	}

	tw.RetryWrite(le, tw.write)
}

// Flush retry send failed events.
func (tw *TeamsWriter) Flush() {
	tw.RetryFlush(tw.write)
}

// Close flush and close.
func (tw *TeamsWriter) Close() {
	tw.Flush()
}

func (tw *TeamsWriter) write(le *log.Event) (err error) {
	if tw.Timeout.Milliseconds() == 0 {
		tw.Timeout = time.Second * 5
	}

	title, text := tw.format(le)

	switch tw.Style {
	case "hero":
		if len(tw.message.Attachments) == 0 {
			tw.message.Type = "message"
			tw.message.Attachments = []any{
				map[string]any{
					"contentType": "application/vnd.microsoft.card.hero",
					"content": map[string]any{
						"title": "",
						"text":  "",
					},
				},
			}
		}
		content := tw.message.Attachments[0].(map[string]any)["content"].(map[string]any)
		content["title"] = title
		content["text"] = text
	case "adaptive":
		if len(tw.message.Attachments) == 0 {
			tw.message.Type = "message"
			tw.message.Attachments = []any{
				map[string]any{
					"contentType": "application/vnd.microsoft.card.adaptive",
					"content": map[string]any{
						"$schema": "https://adaptivecards.io/schemas/adaptive-card.json",
						"type":    "AdaptiveCard",
						"version": "1.4",
						"body": []map[string]any{
							{
								"type":  "TextBlock",
								"style": "heading",
								"wrap":  true,
								"text":  "",
							},
							{
								"type": "RichTextBlock",
								"inlines": []map[string]any{
									{
										"type":     "TextRun",
										"fontType": "Monospace",
										"text":     "",
									},
								},
							},
						},
					},
				},
			}
		}
		body := tw.message.Attachments[0].(map[string]any)["content"].(map[string]any)["body"].([]map[string]any)
		body[0]["text"] = title
		body[1]["inlines"].([]map[string]any)[0]["text"] = text
	default:
		tw.message.Title, tw.message.Text = title, text
	}

	if err = teams.Post(tw.Webhook, tw.Timeout, &tw.message); err != nil {
		err = fmt.Errorf("teamslog: Post(%q): %w", tw.Webhook, err)
	}
	return
}

// format format log event to (message)
func (tw *TeamsWriter) format(le *log.Event) (sub, msg string) {
	sbs := tw.SubFormat(le)
	sub = str.UnsafeString(sbs)
	msl := gog.If(tw.MaxSubLength <= 0, 200, tw.MaxSubLength)
	sub = str.Ellipsis(sub, msl)

	mbs := tw.Format(le)
	msg = str.UnsafeString(mbs)
	mml := gog.If(tw.MaxMsgLength <= 0, 2000, tw.MaxMsgLength)
	msg = str.Ellipsis(msg, mml)
	return
}

func init() {
	log.RegisterWriter("teams", func() log.Writer {
		return &TeamsWriter{}
	})
}
