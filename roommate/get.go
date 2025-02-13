package roommate

import (
	"fmt"
	"net/http"
)

// swagger:route GET /v1/roommate/{id} roommate getRoommate
// return a roommate
// responses:
//	200: roommateResponse
//	404: roommateErrorResponse

// HandleGetRoommate handles GET requests
func (h *Handler) HandleGetRoommate(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(KeyId{}).(string)
	roommate, err := h.GetRoommate(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting roommate by id: %s, %#v\n", id, err), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = roommate.Write(w)
	if err != nil {
		http.Error(w, fmt.Sprintf("encoding error: %#v", err), http.StatusInternalServerError)
	}
}

// swagger:route GET /v1/roommate roommate getRoommatesByUserId
// return all roommates for a user
// responses:
//	200: roommatesResponse
//	404: roommateErrorResponse

// HandleGetRoommateByUserId handles GET requests
func (h *Handler) HandleGetRoommateByUserId(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(KeyUserId{}).(string)
	roommates, err := h.GetRoommateByUserId(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting roommate by user id: %s, %#v\n", id, err), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = WriteAll(roommates, w)
	if err != nil {
		http.Error(w, fmt.Sprintf("encoding error: %#v", err), http.StatusInternalServerError)
	}
}

func (h *Handler) GetRoommate(id string) (*RoommateDto, error) {
	result, err := NewRoommateDatabase(h.db).SelectRoommate(id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("no results for id: %s", id)
	}
	if len(result) != 1 {
		return nil, fmt.Errorf("duplicate results for id: %s", id)
	}
	return &result[0], err
}

func (h *Handler) GetRoommateByUserId(userId string) (*[]RoommateDto, error) {
	result, err := NewRoommateDatabase(h.db).SelectRoommatesByUserId(userId)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("no results for user id: %s", userId)
	}
	return &result, err
}
