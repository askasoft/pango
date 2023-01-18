package freshservice

import (
	"github.com/pandafw/pango/sdk/fdk"
)

type RateLimitedError = fdk.RateLimitedError

type Time = fdk.Time

type Attachment = fdk.Attachment

type ListOption = fdk.ListOption

type File = fdk.File

type Files = fdk.Files

type WithFiles = fdk.WithFiles

type Values = fdk.Values

func NewAttachment(file string, data ...[]byte) *Attachment {
	return fdk.NewAttachment(file, data...)
}

func toString(o any) string {
	return fdk.ToString(o)
}
