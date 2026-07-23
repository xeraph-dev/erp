package codecs

import "io"

type Codec struct {
	decoder Decoder
	encoder Encoder
}

var DefaultCodec, _ = NewCodec("", "")

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
		codec.encoder = jsonEncoder{}
	default:
		ok = false
	}
	return
}
func (codec Codec) Decode(r io.Reader, v any) error { return codec.decoder.Decode(r, v) }
func (codec Codec) Encode(w io.Writer, v any) error { return codec.encoder.Encode(w, v) }
func (codec Codec) DecodeContentType() string       { return codec.decoder.ContentType() }
func (codec Codec) EncodeContentType() string       { return codec.encoder.ContentType() }

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
