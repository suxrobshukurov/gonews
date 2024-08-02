package main

import (
	"gonews/pkg/api"
	"gonews/pkg/storage"
	"gonews/pkg/storage/memdb"
	"gonews/pkg/storage/mongo"
	"gonews/pkg/storage/postgres"
	"log"
	"net/http"
	"os"
)

// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.Api
}

func main() {
	// Создаём объект сервера.
	var srv server

	// Создаём объекты баз данных.
	//
	// БД в памяти.
	db1 := memdb.New()

	// Реляционная БД PostgreSQL.
	pwd := os.Getenv("dbpass")
	if pwd == "" {
		os.Exit(1)
	}

	db2, err := postgres.New("postgres://postgres:" + pwd + "@localhost:5432/news?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	
	// Документная БД MongoDB.
	db3, err := mongo.New("mongodb://localhost:27017/")
	if err != nil {
		log.Fatal(err)
	}
	// _, _ = db2, db3
	_, _, _ = db1, db2, db3

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db3

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	http.ListenAndServe(":8080", srv.api.Router())
}
