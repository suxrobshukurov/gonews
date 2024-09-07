package memdb

import (
	"sync"

	"github.com/suxrobshukurov/gonews/pkg/storage"
)

type DB struct {
	m     sync.Mutex
	id    int
	store map[int]storage.Post
}

// New creates a new memdb storage
func New() (*DB, error) {
	db := DB{
		id:    1,
		store: make(map[int]storage.Post),
	}
	return &db, nil
}

// Posts returns a list of posts from the database
// n is the number of posts to return
func (db *DB) Posts(n int) ([]storage.Post, error) {
	db.m.Lock()
	defer db.m.Unlock()
	if n == 0 {
		n = 10
	}
	if n > len(db.store) {
		n = len(db.store)
	}
	posts := make([]storage.Post, 0, n)
	for _, post := range db.store {
		posts = append(posts, post)
	}
	return posts, nil
}

// AddPosts adds a list of posts to the database
func (db *DB) AddPosts(posts []storage.Post) error {
	db.m.Lock()
	defer db.m.Unlock()
	for _, p := range posts {
		p.ID = db.id
		db.store[p.ID] = p
		db.id++
	}
	return nil
}
