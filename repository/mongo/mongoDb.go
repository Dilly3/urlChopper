package mongo

import (
	"context"
	"time"

	"github.com/dilly3/urlshortner/internal"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

const (
	REDIRECTS = "redirects"
)

func NewMongoRepository(mongoUrl string, mongoTimeout int, mongoDbName string) (*MongoRepository, error) {

	repo := &MongoRepository{}

	client, err := newMongoClient(mongoUrl, mongoTimeout)
	if err != nil {
		return nil, err
	}
	repo.client = client
	repo.database = mongoDbName
	repo.timeout = time.Duration(mongoTimeout)
	return repo, nil
}

func newMongoClient(mongoUrl string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))

	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (m *MongoRepository) Find(code string) (*internal.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	redirect := &internal.Redirect{}
	collection := getCollection(m.client, m.database, REDIRECTS)
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(redirect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(internal.ErrRedirectNotFound, "repository.mongoDb.Find")
		}
		return nil, errors.Wrap(err, "repository.mongoDb.Find")
	}
	return redirect, nil

}

func (m *MongoRepository) Store(redirect internal.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	collection := getCollection(m.client, m.database, REDIRECTS)
	_, err := collection.InsertOne(
		ctx, bson.M{
			"code":       redirect.Code,
			"url":        redirect.Url,
			"created_at": redirect.CreatedAt,
		},
	)
	if err != nil {
		return errors.Wrap(err, "repository.mongoDb.Store")
	}
	return nil
}

func getCollection(mClient *mongo.Client, database string, collection string) *mongo.Collection {
	return mClient.Database(database).Collection(collection)
}
