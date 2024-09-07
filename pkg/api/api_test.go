package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/suxrobshukurov/gonews/pkg/storage"
	"github.com/suxrobshukurov/gonews/pkg/storage/memdb"
)

func TestAPI_postsHandler(t *testing.T) {
	db, _ := memdb.New()
	db.AddPosts([]storage.Post{
		{
			Title: "Test Post",
		},
		{
			Title: "Test Post",
		},
	})
	api := New(db)
	req := httptest.NewRequest("GET", "/news/2", nil)
	rw := httptest.NewRecorder()
	api.r.ServeHTTP(rw, req)

	if !(rw.Code == http.StatusOK) {
		t.Errorf("expected %d, got %d", http.StatusOK, rw.Code)
	}
	b, err := io.ReadAll(rw.Body)
	if err != nil {
		t.Fatal(err)
	}

	var posts []storage.Post
	err = json.Unmarshal(b, &posts)
	if err != nil {
		t.Fatal(err)
	}
	wantLen := 2
	if len(posts) != wantLen {
		t.Errorf("expected %d, got %d", wantLen, len(posts))
	}
}
