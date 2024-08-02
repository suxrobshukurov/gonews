package storage

type Post struct {
	ID          int
	AuthorID    int
	Title       string
	Content     string
	CreatedAt   int64
	PublishedAt int64
}

type Interface interface {
	Posts() ([]Post, error) // возвращает все посты
	AddPost(Post) error     // добавляет пост
	UpdatePost(Post) error  // обновляет пост
	DeletePost(Post) error  // удаляет пост
}
