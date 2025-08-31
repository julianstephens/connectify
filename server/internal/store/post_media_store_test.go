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

func TestCreateAndGetPostMedia(t *testing.T) {
	dbConn := tests.SetupTestDB(t)
	mediaStore := NewPostMediaStore(dbConn)
	ctx := context.Background()

	// Create a post to attach media to
	postParams := db.CreatePostParams{
		ID:             uuid.New(),
		AuthorID:       "author1",
		Content:        "Test post for media",
		ContentHtml:    sql.NullString{String: "<p>Test post for media</p>", Valid: true},
		Visibility:     1,
		ReplyToPostID:  uuid.NullUUID{UUID: uuid.Nil, Valid: false},
		OriginalPostID: uuid.NullUUID{UUID: uuid.Nil, Valid: false},
		Language:       sql.NullString{String: "en", Valid: true},
		Meta:           pqtype.NullRawMessage{Valid: false},
	}
	queries := db.New(dbConn)
	_, err := queries.CreatePost(ctx, postParams)
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	mediaParams := db.CreatePostMediaParams{
		ID:        uuid.New(),
		PostID:    postParams.ID,
		Url:       "https://example.com/image.jpg",
		MediaType: "image/jpeg",
		Width:     sql.NullInt32{Int32: 800, Valid: true},
		Height:    sql.NullInt32{Int32: 600, Valid: true},
		SizeBytes: sql.NullInt64{Int64: 102400, Valid: true},
		Meta:      pqtype.NullRawMessage{Valid: false},
		SortIndex: sql.NullInt32{Int32: 1, Valid: true},
	}

	created, err := mediaStore.CreatePostMedia(ctx, &mediaParams)
	if err != nil {
		t.Fatalf("CreatePostMedia failed: %v", err)
	}
	if created.ID != mediaParams.ID {
		t.Errorf("expected media ID %v, got %v", mediaParams.ID, created.ID)
	}

	mediaList, err := mediaStore.GetPostMedia(ctx, postParams.ID)
	if err != nil {
		t.Fatalf("GetPostMedia failed: %v", err)
	}
	if len(mediaList) != 1 {
		t.Errorf("expected 1 media, got %d", len(mediaList))
	}
	if mediaList[0].Url != mediaParams.Url {
		t.Errorf("expected media URL %v, got %v", mediaParams.Url, mediaList[0].Url)
	}
}
