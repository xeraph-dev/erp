package codecs

import (
	"encoding/json"
	"io"
)

type jsonDecoder struct{}
type jsonEncoder struct{}

var _ Decoder = (*jsonDecoder)(nil)
var _ Encoder = (*jsonEncoder)(nil)

func (jsonDecoder) __internal()                     {}
func (jsonEncoder) __internal()                     {}
func (jsonDecoder) ContentType() string             { return "application/json" }
func (jsonEncoder) ContentType() string             { return "application/json" }
func (jsonDecoder) Decode(r io.Reader, v any) error { return json.NewDecoder(r).Decode(v) }
func (jsonEncoder) Encode(r io.Writer, v any) error { return json.NewEncoder(r).Encode(v) }
