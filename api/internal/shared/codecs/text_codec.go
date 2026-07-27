package codecs

import (
	"io"
	"strings"
)

type textDecoder struct{}
type textEncoder struct{}

var _ Decoder = (*textDecoder)(nil)
var _ Encoder = (*textEncoder)(nil)

func (textDecoder) __internal()         {}
func (textEncoder) __internal()         {}
func (textDecoder) ContentType() string { return "text/plain" }
func (textEncoder) ContentType() string { return "text/plain" }
func (textDecoder) Decode(r io.Reader, v any) error {
	b := new(strings.Builder)
	_, err := io.Copy(b, r)
	if err != nil {
		return err
	}
	v = b.String()
	return nil
}
func (textEncoder) Encode(w io.Writer, v any) error {
	_, err := strings.NewReader(v.(string)).WriteTo(w)
	return err
}
