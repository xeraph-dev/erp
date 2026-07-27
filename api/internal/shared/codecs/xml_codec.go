package codecs

import (
	"encoding/xml"
	"io"
)

type xmlDecoder struct{}
type xmlEncoder struct{}

var _ Decoder = (*xmlDecoder)(nil)
var _ Encoder = (*xmlEncoder)(nil)

func (xmlDecoder) __internal()                     {}
func (xmlEncoder) __internal()                     {}
func (xmlDecoder) ContentType() string             { return "application/xml" }
func (xmlEncoder) ContentType() string             { return "application/xml" }
func (xmlDecoder) Decode(r io.Reader, v any) error { return xml.NewDecoder(r).Decode(v) }
func (xmlEncoder) Encode(w io.Writer, v any) error { return xml.NewEncoder(w).Encode(v) }
