package codecs

import (
	"encoding/xml"
	"io"
)

type xmlCodec struct{}

var _ Codec = (*xmlCodec)(nil)

func (xmlCodec) __internal()                     {}
func (xmlCodec) ContentType() string             { return "text/xml" }
func (xmlCodec) Decode(r io.Reader, v any) error { return xml.NewDecoder(r).Decode(v) }
func (xmlCodec) Encode(w io.Writer, v any) error { return xml.NewEncoder(w).Encode(v) }
