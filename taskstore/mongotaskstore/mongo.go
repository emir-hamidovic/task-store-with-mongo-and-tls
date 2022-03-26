package mongotaskstore

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"rest/taskstore"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	client            *mongo.Client
	ctx               context.Context
	cancel            context.CancelFunc
	database          string
	collection        string
	collectionHandler *mongo.Collection
}

func NewMongoServer(uri, database, collection string) *Mongo {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	if database == "" {
		database = "rest_db"
	}

	if collection == "" {
		collection = "tasks"
	}
	collectionHandler := client.Database(database).Collection(collection)

	return &Mongo{client, ctx, cancel, database, collection, collectionHandler}
}

func (m *Mongo) CloseMongoServer() {
	defer m.cancel()

	defer func() {

		if err := m.client.Disconnect(m.ctx); err != nil {
			panic(err)
		}
	}()
}

func (m *Mongo) CreateTask(text string, tags []string, due string) string {
	document := bson.D{
		primitive.E{Key: "text", Value: text},
		primitive.E{Key: "tags", Value: tags},
		primitive.E{Key: "due", Value: due},
	}

	result, err := m.collectionHandler.InsertOne(m.ctx, document)
	if err != nil {
		panic(err)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return id.Hex()
	}

	return ""
}

func (m *Mongo) GetTaskById(id string) (taskstore.Task, error) {
	idPrimitive, err := primitive.ObjectIDFromHex("5d36277e024f042ff4837ad5")
	if err != nil {
		return taskstore.Task{}, err
	}

	var task taskstore.Task
	if err := m.collectionHandler.FindOne(m.ctx, bson.M{"_id": idPrimitive}).Decode(&task); err != nil {
		return taskstore.Task{}, err
	}

	return task, nil
}

func (m *Mongo) DeleteTask(id string) {
	idPrimitive, err := primitive.ObjectIDFromHex("5d36277e024f042ff4837ad5")
	if err != nil {
		log.Fatal("primitive.ObjectIDFromHex ERROR:", err)
	}

	res, err := m.collectionHandler.DeleteOne(m.ctx, bson.M{"_id": idPrimitive})
	fmt.Println("DeleteOne Result TYPE:", reflect.TypeOf(res))

	if err != nil {
		log.Fatal("DeleteOne() ERROR:", err)
	}

	if res.DeletedCount == 0 {
		fmt.Println("DeleteOne() document not found:", res)
	} else {
		fmt.Println("DeleteOne Result:", res)
	}
}

func (m *Mongo) DeleteAll() {
	res, err := m.collectionHandler.DeleteMany(m.ctx, bson.M{})
	if err != nil {
		log.Fatal("DeleteMany() ERROR:", err)
	}

	fmt.Printf("DeleteMany() TOTAL: %d", res.DeletedCount)
}

func (m *Mongo) GetAllTasks() []taskstore.Task {
	cursor, err := m.collectionHandler.Find(m.ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var tasks []taskstore.Task
	defer cursor.Close(m.ctx)
	for cursor.Next(m.ctx) {
		var task taskstore.Task
		if err = cursor.Decode(&task); err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}

	return tasks
}
