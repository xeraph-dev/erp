package middlewares

import (
	"context"
	"erp/internal/codecs"
	"errors"
	"mime"
	"net/http"
)

var (
	ErrNotAllowedContentType = errors.New("not allowed content type")
)

const codecKey = "codec"

func Codec(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentTypeRaw := r.Header.Get("Content-Type")
		if contentTypeRaw == "" {
			http.Error(w, "missing Content-Type header", http.StatusBadRequest)
			return
		}
		acceptRaw := r.Header.Get("Accept")
		if acceptRaw == "" {
			http.Error(w, "missing Accept header", http.StatusBadRequest)
			return
		}

		contentType, _, err := mime.ParseMediaType(contentTypeRaw)
		if err != nil {
			http.Error(w, "malformed Content-Type header", http.StatusBadRequest)
			return
		}
		accept, _, err := mime.ParseMediaType(acceptRaw)
		if err != nil {
			http.Error(w, "malformed Accept header", http.StatusBadRequest)
			return
		}

		codec, ok := codecs.NewCodec(contentType, accept)
		if !ok {
			http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
			return
		}

		w.Header().Set("Content-Type", codec.EncodeContentType())

		ctx := context.WithValue(r.Context(), codecKey, codec)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetCodec(ctx context.Context) codecs.Codec {
	return ctx.Value(codecKey).(codecs.Codec)
}
