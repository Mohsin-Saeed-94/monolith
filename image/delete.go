package image

import (
	"fmt"
	"net/http"
)

// swagger:route DELETE /api/image/{id} image deleteImage
// delete a image
//
// responses:
//	201: noContentResponse
//  404: imageErrorResponse
//  500: imageErrorResponse

// Delete handles DELETE requests and removes items from the database
func (h *Handler) HandleDeleteImage(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(KeyId{}).(string)
	err := h.client.DeleteImage(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("error deleting image by id: %s, %#v\n", id, err), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
