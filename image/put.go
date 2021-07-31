package image

import (
	"encoding/json"
	"fmt"
	"image"
	"net/http"

	"github.com/google/uuid"
)

type ImageIdDto struct {
	Id string `json:"id"`
}

// swagger:route PUT /api/image image putImage
// update a image
//
// responses:
//	200: imageIdResponse
//  404: imageErrorResponse
//  422: imageErrorValidation

// Update handles PUT requests to update images
func (h *Handler) HandlePutImage(w http.ResponseWriter, r *http.Request) {
	i := r.Context().Value(KeyBody{}).(*image.Image)
	id := uuid.New().String()
	h.logger.Printf("id: %s", id)
	err := h.client.PutImage(id, i)
	if err != nil {
		http.Error(w, fmt.Sprintf("error putting image: %#v\n", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	d := getImageIdDtoFromId(id)
	err = json.NewEncoder(w).Encode(d)
	if err != nil {
		http.Error(w, fmt.Sprintf("error writing image id %s: %#v\n", id, err), http.StatusInternalServerError)
		return
	}
}

func getImageIdDtoFromId(id string) *ImageIdDto {
	return &ImageIdDto{
		Id: id,
	}
}
