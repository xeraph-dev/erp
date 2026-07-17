package codecs

import (
	"encoding/json"
	"io"
)

type jsonCodec struct{}

var _ Codec = (*jsonCodec)(nil)

func (jsonCodec) __internal()                     {}
func (jsonCodec) ContentType() string             { return "application/json" }
func (jsonCodec) Decode(r io.Reader, v any) error { return json.NewDecoder(r).Decode(v) }
func (jsonCodec) Encode(w io.Writer, v any) error { return json.NewEncoder(w).Encode(v) }
