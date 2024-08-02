package mongo

import (
	"context"
	"gonews/pkg/storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	database       = "news"
	collectionName = "posts"
)

type Store struct {
	db *mongo.Client
}

func New(constr string) (*Store, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(constr))
	if err != nil {
		return nil, err
	}
	return &Store{db: client}, nil
}
func (s *Store) Posts() ([]storage.Post, error) {
	collection := s.db.Database(database).Collection(collectionName)
	filter := bson.M{}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var posts []storage.Post

	for cursor.Next(context.Background()) {
		var p storage.Post
		err = cursor.Decode(&p)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (s *Store) AddPost(p storage.Post) error {
	collection := s.db.Database(database).Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), p)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdatePost(p storage.Post) error {
	collection := s.db.Database(database).Collection(collectionName)
	_, err := collection.UpdateOne(context.Background(), bson.M{"id": p.ID}, bson.M{"$set": p})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) DeletePost(p storage.Post) error {
	collection := s.db.Database(database).Collection(collectionName)
	_, err := collection.DeleteOne(context.Background(), bson.M{"id": p.ID})
	if err != nil {
		return err
	}
	return nil
}
