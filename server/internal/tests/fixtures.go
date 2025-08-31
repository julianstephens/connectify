package tests

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// SetupTestDB creates an in-memory SQLite DB with the posts and post_media tables for testing
func SetupTestDB(t *testing.T) *sql.DB {
	dbConn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	// Create posts table
	postsSchema := `CREATE TABLE posts (
		id TEXT PRIMARY KEY,
		author_id TEXT NOT NULL,
		content TEXT NOT NULL,
		content_html TEXT,
		visibility INTEGER,
		reply_to_post_id TEXT,
		original_post_id TEXT,
		language TEXT,
		meta TEXT,
		likes_count INTEGER DEFAULT 0,
		comments_count INTEGER DEFAULT 0,
		shares_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		search_vector TEXT
	);`
	_, err = dbConn.Exec(postsSchema)
	if err != nil {
		t.Fatalf("failed to create posts table: %v", err)
	}
	// Create post_media table
	mediaSchema := `CREATE TABLE post_media (
		id TEXT PRIMARY KEY,
		post_id TEXT NOT NULL,
		url TEXT NOT NULL,
		media_type TEXT NOT NULL,
		width INTEGER,
		height INTEGER,
		size_bytes INTEGER,
		meta TEXT,
		sort_index INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = dbConn.Exec(mediaSchema)
	if err != nil {
		t.Fatalf("failed to create post_media table: %v", err)
	}
	return dbConn
}
