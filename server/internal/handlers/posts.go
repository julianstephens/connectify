package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/julianstephens/connectify/server/internal/config"
	"github.com/julianstephens/connectify/server/internal/db"
	"github.com/julianstephens/connectify/server/internal/store"
)

type PostHandler struct {
	store  *store.PostStore
	logger config.LoggerInterface
}

func NewPostHandler(store *store.PostStore) *PostHandler {
	return &PostHandler{
		store:  store,
		logger: config.GetLogger(),
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req db.CreatePostParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	post, err := h.store.CreatePost(r.Context(), &req)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	var id uuid.UUID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing post ID", http.StatusBadRequest)
		return
	}

	var err error
	id, err = uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.store.GetPost(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to get post", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(post)
}
