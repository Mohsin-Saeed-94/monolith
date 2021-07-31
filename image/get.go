package image

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
)

// swagger:route GET /api/image/{id} image getImage
// return a image
// responses:
//	200: imageResponse
//	404: imageErrorResponse

// HandleGetImage handles GET requests
func (h *Handler) HandleGetImage(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(KeyId{}).(string)
	i, err := h.client.GetImage(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting image by id: %s, %#v\n", id, err), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	encodingErr, writeErr := write(i, w)
	if encodingErr != nil {
		h.logger.Printf("encoding error: %#v", err)
	}
	if writeErr != nil {
		h.logger.Printf("write error: %#v", err)
	}
}

func write(i *image.Image, w io.Writer) (error, error) {
	buffer := new(bytes.Buffer)
	encodingErr := png.Encode(buffer, *i)
	d := buffer.Bytes()
	_, writeErr := w.Write(d)
	return encodingErr, writeErr
}
