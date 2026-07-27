package middlewares

import (
	"context"
	"erp/internal/shared/codecs"
	"mime"
	"net/http"
)

type codecCtxKey int

const codecKey codecCtxKey = iota

func Codec(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentTypeString := r.Header.Get("Content-Type")
		if contentTypeString == "" {
			contentTypeString = "*/*"
		}
		contentType, _, err := mime.ParseMediaType(contentTypeString)
		if err != nil {
			http.Error(w, "malformed Content-Type header", http.StatusBadRequest)
			return
		}

		acceptString := r.Header.Get("Accept")
		if acceptString == "" {
			acceptString = "*/*"
		}
		// TODO: implement propper parsing
		// - Accept header often includes more than one media type
		accept, _, err := mime.ParseMediaType(acceptString)
		if err != nil {
			http.Error(w, "malformed Accept header", http.StatusBadRequest)
			return
		}

		codec, ok := codecs.NewCodec(contentType, accept)
		if !ok {
			http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
			return
		}

		ctx := context.WithValue(r.Context(), codecKey, codec)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetCodec(ctx context.Context) codecs.Codec {
	if codec, ok := ctx.Value(codecKey).(codecs.Codec); ok {
		return codec
	}
	return codecs.Default()
}
