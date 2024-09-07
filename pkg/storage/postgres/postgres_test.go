package postgres

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/suxrobshukurov/gonews/pkg/storage"
)

func TestNew(t *testing.T) {
	_, err := New()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDB_Posts(t *testing.T) {
	db, err := New()
	if err != nil {
		t.Fatal(err)
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	posts := []storage.Post{
		{
			Title: "Test Post",
			Link:  strconv.Itoa(r.Intn(1_000_000)),
		},
		{
			Title: "Test Post",
			Link:  strconv.Itoa(r.Intn(1_000_000)),
		},
	}

	err = db.AddPosts(posts)
	if err != nil {
		t.Fatal(err)
	}
	posts, err = db.Posts(2)

	if err != nil {
		t.Fatal(err)
	}
	t.Log(posts)
}
