package dataurl

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/askasoft/pango/str"
)

func Encode(mediaType string, data []byte) string {
	sb := &strings.Builder{}
	_, _ = EncodeWrite(sb, mediaType, data)
	return sb.String()
}

func EncodeWrite(w io.Writer, mediaType string, data []byte) (cnt int, err error) {
	var n int

	n, err = fmt.Fprintf(w, "data:%s;base64,", mediaType)
	cnt += n
	if err != nil {
		return
	}

	enc := base64.NewEncoder(base64.StdEncoding, w)
	n, err = enc.Write(data)
	cnt += n
	if err != nil {
		return
	}

	err = enc.Close()
	return
}

func Decode(dataurl string) (mediaType string, data []byte, err error) {
	prefix, data64, ok := str.CutByte(dataurl, ',')
	if !ok || !str.StartsWith(prefix, "data:") || !str.EndsWith(prefix, ";base64") {
		err = errors.New("invalid dataurl")
		return
	}

	mediaType = prefix[5 : len(prefix)-7]
	data, err = base64.StdEncoding.DecodeString(data64)
	return
}
