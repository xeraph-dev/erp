package codecs

import "io"

type Codec interface {
	__internal()
	ContentType() string
	Decode(r io.Reader, v any) error
	Encode(w io.Writer, v any) error
}

var registry = map[string]Codec{
	"application/json": jsonCodec{},
	"text/xml":         xmlCodec{},
}

func Get(contentType string) (codec Codec, ok bool) {
	codec, ok = registry[contentType]
	return
}
