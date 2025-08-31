package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/julianstephens/connectify/server/internal/config"
	"github.com/julianstephens/connectify/server/internal/db"
	"github.com/julianstephens/connectify/server/internal/store"
)

type PostMedia struct {
	store  *store.PostMediaStore
	logger config.LoggerInterface
}

func NewPostMediaHandler(store *store.PostMediaStore) *PostMedia {
	return &PostMedia{
		store:  store,
		logger: config.GetLogger(),
	}
}

func (h *PostMedia) UploadPostMedia(w http.ResponseWriter, r *http.Request) {
	// var req db.CreatePostMediaParams
	var media db.PostMedia

	// TODO: implement media handling

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(media)
}

func (h *PostMedia) GetPostMedia(w http.ResponseWriter, r *http.Request) {
	var id uuid.UUID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing post media ID", http.StatusBadRequest)
		return
	}

	var err error
	id, err = uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid post media ID", http.StatusBadRequest)
		return
	}

	media, err := h.store.GetPostMedia(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to get post media: %v", err)
		http.Error(w, "failed to get post media", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(media)
}

func (h *PostMedia) DeletePostMedia(w http.ResponseWriter, r *http.Request) {
	var id uuid.UUID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing post media ID", http.StatusBadRequest)
		return
	}

	var err error
	id, err = uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid post media ID", http.StatusBadRequest)
		return
	}

	if err := h.store.DeletePostMedia(r.Context(), id); err != nil {
		h.logger.Error("failed to delete post media: %v", err)
		http.Error(w, "failed to delete post media", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
