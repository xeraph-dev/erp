package codecs

import (
	"io"
	"net/http"
)

type Codec struct {
	decoder Decoder
	encoder Encoder
}

var DefaultCodec, _ = NewCodec("application/json", "application/json")

func NewCodec(contentType string, accept string) (codec Codec, ok bool) {
	ok = true
	switch contentType {
	case "*/*", "application/json":
		codec.decoder = jsonDecoder{}
	case "application/xml", "text/xml":
		codec.decoder = xmlDecoder{}
	default:
		ok = false
	}
	switch accept {
	case "*/*", "application/json":
		codec.encoder = jsonEncoder{}
	case "application/xml", "text/xml":
		codec.encoder = xmlEncoder{}
	default:
		ok = false
	}
	return
}
func (codec Codec) Decode(r *http.Request, v any) error { return codec.decoder.Decode(r.Body, v) }
func (codec Codec) Encode(w http.ResponseWriter, v any) error {
	w.Header().Set("Content-Type", codec.encoder.ContentType())
	return codec.encoder.Encode(w, v)
}

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
