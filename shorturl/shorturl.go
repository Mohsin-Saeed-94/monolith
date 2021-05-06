package shorturl

import (
	"log"
	"net/http"

	"github.com/turnkeyca/monolith/bitly"
)

type Handler struct {
	logger      *log.Logger
	bitlyClient *bitly.Client
}

func NewHandler(logger *log.Logger, bitlyClient *bitly.Client) *Handler {
	return &Handler{
		logger:      logger,
		bitlyClient: bitlyClient,
	}
}

func (h *Handler) HandleGetShortUrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	shortUrl := h.GetShortUrl(r.URL.Query().Get("url"))
	err := shortUrl.Write(w)
	if err != nil {
		h.logger.Printf("encoding error: %#v", err)
	}
}

func (h *Handler) GetShortUrl(longUrl string) *Dto {
	return New(h.bitlyClient.GetShortUrl(longUrl))
}
