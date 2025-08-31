package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/julianstephens/connectify/server/internal/db"
)

type PostStore struct {
	q *db.Queries
}

func NewPostStore(dbConn *sql.DB) *PostStore {
	return &PostStore{
		q: db.New(dbConn),
	}
}

// CreatePost creates a new post in the database
func (s *PostStore) CreatePost(ctx context.Context, post *db.CreatePostParams) (db.Posts, error) {
	return s.q.CreatePost(ctx, *post)
}

// GetPost retrieves a post by its ID
func (s *PostStore) GetPost(ctx context.Context, id uuid.UUID) (db.Posts, error) {
	return s.q.GetPost(ctx, id)
}

// ListPostsOffset retrieves a list of posts with pagination
func (s *PostStore) ListPostsOffset(ctx context.Context, author string, limit, offset int64) ([]db.Posts, error) {
	return s.q.ListUserPostsOffset(ctx, db.ListUserPostsOffsetParams{
		AuthorID: author,
		Limit:    limit,
		Offset:   offset,
	})
}

// CountPosts retrieves the total number of posts for a given author
func (s *PostStore) CountPosts(ctx context.Context, author string) (int64, error) {
	return s.q.CountUserPosts(ctx, author)
}

// ListUserPostsFirstPage retrieves the first page of posts for a given user
func (s *PostStore) ListUserPostsFirstPage(ctx context.Context, author string, limit int64) ([]db.Posts, error) {
	return s.q.ListUserPostsFirstPage(ctx, db.ListUserPostsFirstPageParams{
		AuthorID: author,
		Limit:    limit,
	})
}

// ListUserPostsAfter retrieves posts for a given user after a specific creation date
func (s *PostStore) ListUserPostsAfter(ctx context.Context, author string, createdAt time.Time, limit int64) ([]db.Posts, error) {
	return s.q.ListUserPostsAfter(ctx, db.ListUserPostsAfterParams{
		AuthorID:  author,
		CreatedAt: createdAt,
		Limit:     limit,
	})
}
