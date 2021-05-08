package user

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (h *Handler) HandlePutUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(KeyId{}).(uuid.UUID)
	dto := r.Context().Value(KeyBody{}).(Dto)
	dto.Id = id
	err := h.UpdateUser(&dto)
	if err != nil {
		http.Error(w, fmt.Sprintf("error updating user: %#v\n", err), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateUser(dto *Dto) error {
	err := h.db.Run("update users set id=%s where id=%s", dto.Id.String(), dto.Id.String())
	return err
}
