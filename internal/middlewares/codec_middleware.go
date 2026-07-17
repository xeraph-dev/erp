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
		raw := r.Header.Get("Content-Type")
		if raw == "" {
			// TODO: add missing content type error
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mediaType, _, err := mime.ParseMediaType(raw)
		if err != nil {
			http.Error(w, "malformed Content-Type", http.StatusBadRequest)
			return
		}

		codec, ok := codecs.Get(mediaType)
		if !ok {
			http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
			return
		}

		w.Header().Set("Content-Type", codec.ContentType())

		ctx := context.WithValue(r.Context(), codecKey, codec)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetCodec(ctx context.Context) codecs.Codec {
	return ctx.Value(codecKey).(codecs.Codec)
}
