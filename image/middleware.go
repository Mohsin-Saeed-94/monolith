package image

import (
	"context"
	"image"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type KeyId struct{}
type KeyBody struct{}

func (h *Handler) GetIdFromPath(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if mux.Vars(r)["id"] == "" {
			next.ServeHTTP(w, r)
			return
		}
		id := uuid.MustParse(mux.Vars(r)["id"]).String()
		ctx := context.WithValue(r.Context(), KeyId{}, id)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) GetBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d, err := read(r.Body)
		if err != nil {
			h.logger.Printf("decoding error: %#v", err)
			http.Error(w, "Error reading image", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), KeyBody{}, d)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func read(r io.Reader) (*image.Image, error) {
	d, _, err := image.Decode(r)
	return &d, err
}
