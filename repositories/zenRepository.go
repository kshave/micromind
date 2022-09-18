package repositories

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/micromind/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectionTimeout        = 5
	connectionStringTemplate = "mongodb://%s:%s@%s"
)

type ZenRepository interface {
	GetRandomQuote() (string, string, error)
	GetRandomQuestion() (string, error)
}

type zenRepository struct {
}

func NewInstanceOfZenRepository() zenRepository {
	return zenRepository{}
}

// Retrieves a client to the DocumentDB
func getConnection() (*mongo.Client, context.Context, context.CancelFunc, error) {
	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	clusterEndpoint := os.Getenv("MONGO_ENDPOINT")

	connectionURI := fmt.Sprintf(connectionStringTemplate, username, password, clusterEndpoint)

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Printf("Failed to create new MongoDB client: %v", err)
		return nil, nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Printf("Failed to connect to client: %v", err)
		return client, ctx, cancel, err
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping cluster: %v", err)
		return client, ctx, cancel, err
	}

	fmt.Println("Connected to MongoDB!")
	return client, ctx, cancel, nil
}

func (zenRepository) GetRandomQuote() (string, string, error) {
	client, ctx, cancel, err := getConnection()
	if err != nil {
		log.Printf("Failed to getConnection to DB: %v", err)
		return "", "", err
	}

	defer cancel()
	defer client.Disconnect(ctx)

	qm := models.QuoteModel{}

	randomSample := []bson.D{bson.D{{"$sample", bson.D{{"size", 1}}}}}
	resultCursor, err := client.Database("micromind").Collection("quotes").Aggregate(ctx, randomSample)
	if err != nil {
		log.Printf("Failed to get new quote cursor: %v", err)
	}
	for resultCursor.Next(ctx) {
		err := resultCursor.Decode(&qm)
		if err != nil {
			log.Printf("Failed to decode quote: %v", err)
			return "", "", err
		}
	}
	return qm.Quote, qm.Author, nil
}

func (zenRepository) GetRandomQuestion() (string, error) {
	client, ctx, cancel, err := getConnection()
	if err != nil {
		log.Printf("Failed to getConnection to DB: %v", err)
		return "", err
	}

	defer cancel()
	defer client.Disconnect(ctx)

	qm := models.QuestionModel{}

	randomSample := []bson.D{bson.D{{"$sample", bson.D{{"size", 1}}}}}
	resultCursor, err := client.Database("micromind").Collection("questions").Aggregate(ctx, randomSample)
	if err != nil {
		log.Printf("Failed to get new question cursor: %v", err)
	}
	for resultCursor.Next(ctx) {
		if err := resultCursor.Decode(&qm); err != nil {
			log.Printf("Failed to decode question: %v", err)
			return "", err
		}
	}
	return qm.Question, nil
}
