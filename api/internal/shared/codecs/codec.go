package codecs

import (
	"io"
)

type Encoder interface {
	__internal()
	ContentType() string
	Encode(r io.Writer, v any) error
}

type Decoder interface {
	__internal()
	ContentType() string
	Decode(w io.Reader, v any) error
}

type Codec struct {
	decoder Decoder
	encoder Encoder
}

func NewCodec(contentType string, accept string) (codec Codec, ok bool) {
	ok = true
	switch contentType {
	case "text/plain":
		codec.decoder = textDecoder{}
	case "*/*", "application/json":
		codec.decoder = jsonDecoder{}
	case "application/xml", "text/xml":
		codec.decoder = xmlDecoder{}
	default:
		ok = false
	}
	switch accept {
	case "text/plain":
		codec.encoder = textEncoder{}
	case "*/*", "application/json":
		codec.encoder = jsonEncoder{}
	case "application/xml", "text/xml":
		codec.encoder = xmlEncoder{}
	default:
		ok = false
	}
	return
}

func Default() Codec {
	codec, _ := NewCodec("application/json", "application/json")
	return codec
}

func (codec Codec) ContentType() string             { return codec.encoder.ContentType() }
func (codec Codec) Decode(r io.Reader, v any) error { return codec.decoder.Decode(r, v) }
func (codec Codec) Encode(w io.Writer, v any) error { return codec.encoder.Encode(w, v) }
