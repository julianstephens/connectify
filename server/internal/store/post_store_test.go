package store

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"

	"github.com/julianstephens/connectify/server/internal/db"
	"github.com/julianstephens/connectify/server/internal/tests"
)

func setupTestDB(t *testing.T) *sql.DB {
	return tests.SetupTestDB(t)
}

func TestCreateAndGetPost(t *testing.T) {
	dbConn := setupTestDB(t)
	store := NewPostStore(dbConn)
	ctx := context.Background()

	postParams := db.CreatePostParams{
		ID:             uuid.New(),
		AuthorID:       "author1",
		Content:        "Hello, world!",
		ContentHtml:    sql.NullString{String: "<p>Hello, world!</p>", Valid: true},
		Visibility:     1,
		ReplyToPostID:  uuid.NullUUID{UUID: uuid.Nil, Valid: false},
		OriginalPostID: uuid.NullUUID{UUID: uuid.Nil, Valid: false},
		Language:       sql.NullString{String: "en", Valid: true},
		Meta:           pqtype.NullRawMessage{Valid: false},
	}

	_, err := store.CreatePost(ctx, &postParams)
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	post, err := store.GetPost(ctx, postParams.ID)
	if err != nil {
		t.Fatalf("GetPost failed: %v", err)
	}
	if post.ID != postParams.ID {
		t.Errorf("expected post ID %v, got %v", postParams.ID, post.ID)
	}
}

func TestListPostsOffset(t *testing.T) {
	dbConn := setupTestDB(t)
	store := NewPostStore(dbConn)
	ctx := context.Background()

	// Insert multiple posts
	for i := 0; i < 5; i++ {
		postParams := db.CreatePostParams{
			ID:             uuid.New(),
			AuthorID:       "author1",
			Content:        "Content",
			ContentHtml:    sql.NullString{String: "<p>Content</p>", Valid: true},
			Visibility:     1,
			ReplyToPostID:  uuid.NullUUID{UUID: uuid.Nil, Valid: false},
			OriginalPostID: uuid.NullUUID{UUID: uuid.Nil, Valid: false},
			Language:       sql.NullString{String: "en", Valid: true},
			Meta:           pqtype.NullRawMessage{Valid: false},
		}
		_, _ = store.CreatePost(ctx, &postParams)
	}

	posts, err := store.ListPostsOffset(ctx, "author1", 2, 0)
	if err != nil {
		t.Fatalf("ListPostsOffset failed: %v", err)
	}
	if len(posts) != 2 {
		t.Errorf("expected 2 posts, got %d", len(posts))
	}
}

func TestCountPosts(t *testing.T) {
	dbConn := setupTestDB(t)
	store := NewPostStore(dbConn)
	ctx := context.Background()

	// Insert posts
	for i := 0; i < 3; i++ {
		postParams := db.CreatePostParams{
			ID:             uuid.New(),
			AuthorID:       "author2",
			Content:        "Content",
			ContentHtml:    sql.NullString{String: "<p>Content</p>", Valid: true},
			Visibility:     1,
			ReplyToPostID:  uuid.NullUUID{UUID: uuid.Nil, Valid: false},
			OriginalPostID: uuid.NullUUID{UUID: uuid.Nil, Valid: false},
			Language:       sql.NullString{String: "en", Valid: true},
			Meta:           pqtype.NullRawMessage{Valid: false},
		}
		_, _ = store.CreatePost(ctx, &postParams)
	}

	count, err := store.CountPosts(ctx, "author2")
	if err != nil {
		t.Fatalf("CountPosts failed: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 posts, got %d", count)
	}
}
