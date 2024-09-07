package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/subosito/gotenv"
	"github.com/suxrobshukurov/gonews/pkg/storage"
)

func init() {
	gotenv.Load()
}

// DB represents a storage
type DB struct {
	pool *pgxpool.Pool
}

// New creates a new postgres storage
func New() (*DB, error) {
	connstr := os.Getenv("connstr")

	if connstr == "" {
		return nil, errors.New("empty connstr")
	}

	p, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	db := DB{
		pool: p,
	}
	return &db, nil
}

// Posts returns a list of posts from the database
// n is the number of posts to return
func (db *DB) Posts(n int) ([]storage.Post, error) {
	if n == 0 {
		n = 10
	}
	rows, err := db.pool.Query(context.Background(), `
		SELECT id, title, content, pub_time, link
		FROM posts
		ORDER BY pub_time DESC
		LIMIT $1
	`, n)
	if err != nil {
		return nil, fmt.Errorf("can't get posts from db: %w", err)
	}

	var posts []storage.Post
	for rows.Next() {
		var post storage.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.PubTime, &post.Link); err != nil {
			return nil, fmt.Errorf("can't scan post: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// AddPosts adds a list of posts to the database
func (db *DB) AddPosts(posts []storage.Post) error {
	for _, p := range posts {
		_, err := db.pool.Exec(context.Background(), `
				INSERT INTO posts ( title, content, pub_time, link)
				VALUES ( $1, $2, $3, $4)
				ON CONFLICT (link) DO UPDATE SET title = EXCLUDED.title, content = EXCLUDED.content, pub_time = EXCLUDED.pub_time
		`,
			p.Title,
			p.Content,
			p.PubTime,
			p.Link,
		)
		if err != nil {
			return fmt.Errorf("can't insert post in db: %w", err)
		}
	}
	return nil
}
