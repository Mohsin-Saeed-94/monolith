package image

import (
	"context"
	"fmt"
	"image"
	"net/http"

	_ "image/png"

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
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing multipart form %#v", err), http.StatusBadRequest)
			return
		}
		file, _, err := r.FormFile("photo")
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing photo from form %#v", err), http.StatusBadRequest)
			return
		}
		i, _, err := image.Decode(file)
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing photo %#v", err), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), KeyBody{}, &i)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
