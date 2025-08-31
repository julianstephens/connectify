package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julianstephens/connectify/server/internal/db"
	"github.com/julianstephens/connectify/server/internal/store"
)

type PaginatedResponse[T any] struct {
	Items      []T    `json:"items"`
	NextCursor string `json:"next_cursor,omitempty"`
	Limit      int32  `json:"limit"`
	Total      *int64 `json:"total,omitempty"`
}

type PostsHandler struct {
	store *store.PostStore
}

func NewPostsHandler(store *store.PostStore) *PostsHandler {
	return &PostsHandler{
		store: store,
	}
}

func (h *PostsHandler) PostsOffsetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()
	limit := parseIntDefault(q.Get("limit"), 20)
	if limit <= 0 {
		limit = 20
	}
	if limit > 200 {
		limit = 200
	}
	offset := parseIntDefault(q.Get("offset"), 0)
	if offset < 0 {
		offset = 0
	}
	author := q.Get("author")
	if author == "" {
		http.Error(w, "Missing author parameter", http.StatusBadRequest)
		return
	}

	posts, err := h.store.ListPostsOffset(ctx, author, int64(limit), int64(offset))
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	var next *string
	if len(posts) == limit {
		n := strconv.Itoa(offset + limit)
		next = &n
	}

	var total int64
	total, err = h.store.CountPosts(ctx, author)
	if err != nil {
		http.Error(w, "Failed to count posts", http.StatusInternalServerError)
		return
	}

	resp := PaginatedResponse[any]{
		Items:      toAnySlice(posts),
		NextCursor: *next,
		Limit:      int32(limit),
		Total:      &total,
	}
	writeJSON(w, resp)
}

func (h *PostsHandler) PostsCursorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()
	limit := parseIntDefault(q.Get("limit"), 20)
	if limit <= 0 {
		limit = 20
	}
	if limit > 200 {
		limit = 200
	}
	author := q.Get("author")
	if author == "" {
		http.Error(w, "Missing author parameter", http.StatusBadRequest)
		return
	}
	cursor := q.Get("cursor")

	var posts []db.Posts
	var err error

	if cursor == "" {
		posts, err = h.store.ListUserPostsFirstPage(ctx, author, int64(limit))
	} else {
		data, cursorErr := decodeCursor(cursor)
		if cursorErr != nil {
			http.Error(w, "Invalid cursor parameter", http.StatusBadRequest)
			return
		}
		posts, err = h.store.ListUserPostsAfter(ctx, author, data.CreatedAt, int64(limit))
	}
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	var nextCursor *string
	if len(posts) > 0 {
		last := posts[len(posts)-1]
		payload := cursorPayload{
			CreatedAt: last.CreatedAt,
		}
		enc := encodeCursor(payload)
		nextCursor = &enc
	}

	var total *int64
	t, err := h.store.CountPosts(ctx, author)
	if err == nil {
		total = &t
	}

	writeJSON(w, PaginatedResponse[db.Posts]{
		Items:      posts,
		NextCursor: *nextCursor,
		Limit:      int32(limit),
		Total:      total,
	})
}

func parseIntDefault(s string, d int) int {
	if s == "" {
		return d
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return d
	}
	return v
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// toAnySlice converts a slice to []any for generic response typing in UsersOffsetHandler.
// Use concrete types in real code to avoid reflection/encoding surprises.
func toAnySlice[T any](in []T) []any {
	out := make([]any, len(in))
	for i := range in {
		out[i] = in[i]
	}
	return out
}
