package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/julianstephens/connectify/server/internal/db"
)

type PostMediaStore struct {
	q *db.Queries
}

func NewPostMediaStore(dbConn *sql.DB) *PostMediaStore {
	return &PostMediaStore{
		q: db.New(dbConn),
	}
}

// CreatePostMedia creates a new media attachment for a post
func (s *PostMediaStore) CreatePostMedia(ctx context.Context, media *db.CreatePostMediaParams) (db.PostMedia, error) {
	if media == nil {
		return db.PostMedia{}, errors.New("missing data to create post media")
	}
	return s.q.CreatePostMedia(ctx, *media)
}

// GetPostMedia retrieves a media attachment for a post
func (s *PostMediaStore) GetPostMedia(ctx context.Context, postID uuid.UUID) ([]db.PostMedia, error) {
	return s.q.GetPostMedia(ctx, postID)
}

// DeletePostMedia deletes all media attachments for a post
func (s *PostMediaStore) DeletePostMedia(ctx context.Context, postID uuid.UUID) error {
	_, err := s.q.DeletePostMedia(ctx, postID)
	return err
}
