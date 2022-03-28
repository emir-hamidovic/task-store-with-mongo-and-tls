package mongotaskstore

import (
	"context"
	"errors"
	"fmt"
	"rest/taskstore"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	client            *mongo.Client
	cancel            context.CancelFunc
	database          string
	collection        string
	collectionHandler *mongo.Collection
}

func NewMongoServer(uri, database, collection string) (*Mongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		defer cancel()
		return nil, err
	}

	if database == "" {
		database = "rest_db"
	}

	if collection == "" {
		collection = "tasks"
	}
	collectionHandler := client.Database(database).Collection(collection)

	return &Mongo{client, cancel, database, collection, collectionHandler}, nil
}

func (m *Mongo) CloseMongoServer() (err error) {
	defer m.cancel()

	defer func() error {
		err = m.client.Disconnect(context.TODO())
		return err
	}()

	return nil
}

func (m *Mongo) CreateTask(text string, tags []string, due string) (string, error) {
	document := bson.D{
		primitive.E{Key: "text", Value: text},
		primitive.E{Key: "tags", Value: tags},
		primitive.E{Key: "due", Value: due},
	}

	result, err := m.collectionHandler.InsertOne(context.TODO(), document)
	if err != nil {
		return "", err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return id.Hex(), nil
	}

	return "", err
}

func (m *Mongo) GetTaskById(id string) (taskstore.Task, error) {
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return taskstore.Task{}, err
	}

	var task taskstore.Task
	if err := m.collectionHandler.FindOne(context.TODO(), bson.M{"_id": idPrimitive}).Decode(&task); err != nil {
		return taskstore.Task{}, err
	}

	return task, nil
}

func (m *Mongo) DeleteTask(id string) error {
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := m.collectionHandler.DeleteOne(context.TODO(), bson.M{"_id": idPrimitive})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("deleteOne() document not found")
	} else {
		fmt.Println("DeleteOne Result:", res)
	}

	return nil
}

func (m *Mongo) DeleteAll() error {
	res, err := m.collectionHandler.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		return err
	}

	fmt.Printf("DeleteMany() TOTAL: %d\n", res.DeletedCount)
	return nil
}

func (m *Mongo) GetAllTasks() ([]taskstore.Task, error) {
	cursor, err := m.collectionHandler.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	var tasks []taskstore.Task
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var task taskstore.Task
		if err = cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
