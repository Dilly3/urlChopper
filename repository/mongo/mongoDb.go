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
	DB        = "urlshorter"
)

func NewMongoRepository(mongoUrl string, mongoTimeout int, mongoDbName string) (*MongoRepository, error) {

	repo := &MongoRepository{}

	client, err := newMongoClient(mongoUrl, mongoTimeout)
	if err != nil {
		return nil, err
	}

	repo.client = client
	repo.database = DB
	repo.timeout = time.Duration(mongoTimeout)
	return repo, nil
}

func newMongoClient(mongoUrl string, mongoTimeout int) (*mongo.Client, error) {

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(mongoUrl).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	DBMigration(client, DB, REDIRECTS)
	return client, nil
}
func (m *MongoRepository) Find(code string) (*internal.Redirect, error) {
	redirect := &internal.Redirect{}

	filter := bson.M{"code": code}
	err := m.getCollection(REDIRECTS).FindOne(context.Background(), filter).Decode(redirect)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(internal.ErrRedirectNotFound, "repository.mongoDb.Find")
		}
		return nil, errors.Wrap(err, "repository.mongoDb.Find")
	}
	return redirect, nil

}

func (m *MongoRepository) Store(redirect *internal.Redirect) error {

	collection := getCollection(m.client, m.database, REDIRECTS)
	_, err := collection.InsertOne(
		context.Background(), bson.M{
			"name":       redirect.Name,
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
	return mClient.Database(DB).Collection(REDIRECTS)
}
func (m *MongoRepository) getCollection(collection string) *mongo.Collection {
	return m.client.Database(DB).Collection(REDIRECTS)
}

func DBMigration(client *mongo.Client, dbName string, collectionName string) {
	db := client.Database(dbName)
	command := bson.D{{Key: "create", Value: REDIRECTS}}
	var result bson.M
	if err := db.RunCommand(context.TODO(), command).Decode(&result); err != nil {
		if err == errors.New("collection already exists") {
			//do nothing
		}
	}
}
